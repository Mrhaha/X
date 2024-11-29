package connection

import (
	"google.golang.org/protobuf/proto"
	"net"
)

// Connection 连接
type Connection interface {
	Name() string

	WriteMsg(ext proto.Message, msg proto.Message) error
	WriteBytes(ext proto.Message, msgid uint32, bytes []byte) error

	Close()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
