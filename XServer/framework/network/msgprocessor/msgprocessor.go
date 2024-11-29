package msgprocessor

import (
	"XServer/framework/ioservice"
	"XServer/framework/network/connection"
	gxbytes "github.com/dubbogo/gost/bytes"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ConnectHandler ...
type ConnectHandler func(connection.Connection)

// CloseHandler ...
type CloseHandler func(connection.Connection)

// BytesHandler ...
type BytesHandler func(conn connection.Connection, ext proto.Message, msgid uint32, bytes []byte)

// MsgHandler ...
type MsgHandler func(conn connection.Connection, ext proto.Message, msg proto.Message)

// MsgHandlerGetter 获取消息处理函数
type MsgHandlerGetter interface {
	GetMsgHandler(typ protoreflect.MessageType) (MsgHandler, bool)
}

// IMsgProcessor 消息处理者
type IMsgProcessor interface {
	// OnConnect on read gorutine 连接.
	OnConnect(conn connection.Connection)
	// OnClose 断开.
	OnClose(conn connection.Connection)
	// OnMessage 消息处理函数.
	OnMessage(conn connection.Connection, ext proto.Message, msgid uint32, msgData []byte) error
}

type msgProcessor struct {
	connectHandler   ConnectHandler
	closeHandler     CloseHandler
	bytesHandler     BytesHandler
	msgHandlerGetter MsgHandlerGetter
	callbackIO       ioservice.IOService
}

// NewMsgProcessor ...
func NewMsgProcessor(io ioservice.IOService, connectHandler ConnectHandler, closeHandler CloseHandler, bytesHandler BytesHandler, msgHandlers MsgHandlerGetter) IMsgProcessor {
	if io == nil {
		panic("init NewMsgProcessor ioservice is nil")
	}
	processor := &msgProcessor{
		connectHandler:   connectHandler,
		closeHandler:     closeHandler,
		bytesHandler:     bytesHandler,
		msgHandlerGetter: msgHandlers,
		callbackIO:       io,
	}
	return processor
}

// OnConnect ...
func (p *msgProcessor) OnConnect(conn connection.Connection) {
	if p.connectHandler != nil {
		p.callbackIO.Post(func() {
			p.connectHandler(conn)
		})
	}
}

// OnClose ...
func (p *msgProcessor) OnClose(conn connection.Connection) {
	if p.closeHandler != nil {
		p.callbackIO.Post(func() {
			p.closeHandler(conn)
		})
	}
}

// OnMessage 消息回调
func (p *msgProcessor) OnMessage(conn connection.Connection, ext proto.Message, msgid uint32, msgdata []byte) error {
	handler, ok := p.findMsgHandler(msgid)
	if !ok {
		if p.bytesHandler != nil {
			if len(msgdata) > 0 {
				pData := gxbytes.AcquireBytes(len(msgdata))
				data := *pData
				copy(data, msgdata)
				p.callbackIO.Post(func() {
					defer gxbytes.ReleaseBytes(pData)
					p.bytesHandler(conn, ext, msgid, data)
				})
			} else {
				p.callbackIO.Post(func() {
					p.bytesHandler(conn, ext, msgid, nil)
				})
			}

		}
		return nil
	}

	msg, err := OnUnmarshal(msgid, msgdata)
	if err != nil {
		return err
	}

	p.callbackIO.Post(func() {
		handler(conn, ext, msg)
	})

	return nil
}

func (p *msgProcessor) findMsgHandler(msgid uint32) (MsgHandler, bool) {
	if p.msgHandlerGetter == nil {
		return nil, false
	}

	typ, ok := MessageType(msgid)
	if !ok {
		return nil, false
	}

	return p.msgHandlerGetter.GetMsgHandler(typ)
}
