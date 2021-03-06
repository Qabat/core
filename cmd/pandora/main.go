// How to use it?
//
// Direct usage: `pandora_OS_ARCH etc/pandora.yaml`

package main

import (
	"crypto/ecdsa"
	"sync"

	"github.com/sonm-io/core/accounts"
	"github.com/spf13/afero"
	"github.com/yandex/pandora/cli"
	importer "github.com/yandex/pandora/core/import"
	"github.com/yandex/pandora/core/register"
)

var (
	globalPrivateKey *ecdsa.PrivateKey
	once             sync.Once
)

func PrivateKey(ethConfig accounts.EthConfig) *ecdsa.PrivateKey {
	once.Do(func() {
		privateKey, err := ethConfig.LoadKey(accounts.Silent())
		if err != nil {
			panic(err)
		}

		globalPrivateKey = privateKey
	})

	return globalPrivateKey
}

func main() {
	fs := afero.NewOsFs()
	importer.Import(fs)

	register.Provider("marketplace/order/info", NewOrderInfoProvider())
	register.Provider("marketplace/order/place", NewOrderPlaceProvider())
	register.Gun("sonm", NewGun, NewDefaultGunConfig)

	cli.Run()
}
