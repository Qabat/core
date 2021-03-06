// Code generated by protoc-gen-go. DO NOT EDIT.
// source: net.proto

package sonm

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Addr struct {
	Protocol string      `protobuf:"bytes,1,opt,name=protocol" json:"protocol,omitempty"`
	Addr     *SocketAddr `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *Addr) Reset()                    { *m = Addr{} }
func (m *Addr) String() string            { return proto.CompactTextString(m) }
func (*Addr) ProtoMessage()               {}
func (*Addr) Descriptor() ([]byte, []int) { return fileDescriptor11, []int{0} }

func (m *Addr) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *Addr) GetAddr() *SocketAddr {
	if m != nil {
		return m.Addr
	}
	return nil
}

type SocketAddr struct {
	// Addr describes an IP address.
	Addr string `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	// Port describes a port number.
	// Actually an `uint16` here. Protobuf is so clear and handy.
	Port uint32 `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
}

func (m *SocketAddr) Reset()                    { *m = SocketAddr{} }
func (m *SocketAddr) String() string            { return proto.CompactTextString(m) }
func (*SocketAddr) ProtoMessage()               {}
func (*SocketAddr) Descriptor() ([]byte, []int) { return fileDescriptor11, []int{1} }

func (m *SocketAddr) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *SocketAddr) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*Addr)(nil), "sonm.Addr")
	proto.RegisterType((*SocketAddr)(nil), "sonm.SocketAddr")
}

func init() { proto.RegisterFile("net.proto", fileDescriptor11) }

var fileDescriptor11 = []byte{
	// 128 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x4b, 0x2d, 0xd1,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0xce, 0xcf, 0xcb, 0x55, 0xf2, 0xe0, 0x62, 0x71,
	0x4c, 0x49, 0x29, 0x12, 0x92, 0xe2, 0xe2, 0x00, 0x0b, 0x27, 0xe7, 0xe7, 0x48, 0x30, 0x2a, 0x30,
	0x6a, 0x70, 0x06, 0xc1, 0xf9, 0x42, 0x2a, 0x5c, 0x2c, 0x89, 0x29, 0x29, 0x45, 0x12, 0x4c, 0x0a,
	0x8c, 0x1a, 0xdc, 0x46, 0x02, 0x7a, 0x20, 0x8d, 0x7a, 0xc1, 0xf9, 0xc9, 0xd9, 0xa9, 0x25, 0x20,
	0xbd, 0x41, 0x60, 0x59, 0x25, 0x13, 0x2e, 0x2e, 0x84, 0x98, 0x90, 0x10, 0x54, 0x0f, 0xc4, 0x2c,
	0x30, 0x1b, 0x24, 0x56, 0x90, 0x5f, 0x54, 0x02, 0x36, 0x87, 0x37, 0x08, 0xcc, 0x4e, 0x62, 0x03,
	0xdb, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xda, 0x8d, 0x12, 0xa6, 0x99, 0x00, 0x00, 0x00,
}
