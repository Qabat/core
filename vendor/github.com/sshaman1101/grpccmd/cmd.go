package grpccmd

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sonm-io/core/accounts"
	"github.com/sonm-io/core/cmd/cli/config"
	"github.com/sonm-io/core/util"
	"github.com/sonm-io/core/util/xgrpc"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	rootCmd = &cobra.Command{}
	remote  = new(string)
	input   = new(string)
)

func init() {
	rootCmd.SetOutput(os.Stdout)

	rootCmd.PersistentFlags().StringVar(remote, "remote", "", "gRPC server endpoint")
	rootCmd.PersistentFlags().StringVar(input, "input", "", "JSON file with request body")
}
func SetCmdInfo(name, short string) {
	rootCmd.Use = fmt.Sprintf("%s [command]", name)
	rootCmd.Short = short
}

func RegisterServiceCmd(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func RunE(method, inT string, newClient func(closer io.Closer) interface{}) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		key, err := openConfigAndLoadEthKey()
		if err != nil {
			return fmt.Errorf("cannot open eth key: %v", err)
		}

		conn, err := dial(key)
		if err != nil {
			return err
		}
		defer conn.Close()

		c := newClient(conn)
		cv := reflect.ValueOf(c)
		method := cv.MethodByName(method)
		if method.IsValid() {
			var err error
			var requestBody = []byte("{}")
			if len(*input) > 0 {
				// we have file name, so read it
				requestBody, err = ioutil.ReadFile(*input)
				if err != nil {
					return fmt.Errorf("cannot read request body from %s: %v", *input, err)
				}
			}

			in := reflect.New(proto.MessageType(inT).Elem()).Interface()
			if len(*input) > 0 {
				if err := json.Unmarshal(requestBody, in); err != nil {
					return err
				}
			}

			result := method.Call([]reflect.Value{
				reflect.ValueOf(context.Background()),
				reflect.ValueOf(in),
			})

			if len(result) != 2 {
				panic("service methods should always return 2 values")
			}

			if !result[1].IsNil() {
				return result[1].Interface().(error)
			}

			out := result[0].Interface()
			responseBody, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				return err
			}

			cmd.Println(string(responseBody))
		}

		return nil
	}
}

// TypeToJson takes structure type generated by protoc, instantiates structure
// and marshal structure to JSON. Used for generating requests templates.
func TypeToJson(inT string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		m := jsonpb.Marshaler{
			EnumsAsInts:  true,
			EmitDefaults: true,
			Indent:       "  ",
		}

		in := reflect.New(proto.MessageType(inT).Elem()).Interface().(proto.Message)
		s, err := m.MarshalToString(in)
		if err != nil {
			cmd.Printf("Cannot marshal %s to string: %s\r\n", inT, err)
			os.Exit(1)
		}

		cmd.Println(string(s))
		return nil
	}
}

// dial build ClientConn wrapper with SONM Wallet auth
func dial(key *ecdsa.PrivateKey) (io.Closer, error) {
	if *remote == "" {
		return nil, fmt.Errorf("remote endpoint address is required")
	}

	ctx := context.Background()

	_, TLSConfig, err := util.NewHitlessCertRotator(ctx, key)
	if err != nil {
		return nil, err
	}

	creds := util.NewTLS(TLSConfig)
	return xgrpc.NewClient(ctx, *remote, creds)
}

func openConfigAndLoadEthKey() (*ecdsa.PrivateKey, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	ko, err := accounts.DefaultKeyOpener(accounts.NewSilentPrinter(), cfg.KeyStore(), cfg.PassPhrase())
	if err != nil {
		return nil, err
	}

	_, err = ko.OpenKeystore()
	if err != nil {
		return nil, err
	}

	key, err := ko.GetKey()
	if err != nil {
		return nil, err
	}

	return key, nil
}
