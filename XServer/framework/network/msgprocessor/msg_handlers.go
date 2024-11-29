package msgprocessor

import (
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MsgHandlers ...
type MsgHandlers struct {
	handlers map[protoreflect.MessageType]MsgHandler
	rw       sync.RWMutex
}

// NewMsgHandlers ...
func NewMsgHandlers() *MsgHandlers {
	m := new(MsgHandlers)
	m.handlers = make(map[protoreflect.MessageType]MsgHandler)
	return m
}

// GetMsgHandler ...
func (m *MsgHandlers) GetMsgHandler(typ protoreflect.MessageType) (MsgHandler, bool) {
	m.rw.RLock()
	handler, ok := m.handlers[typ]
	m.rw.RUnlock()
	return handler, ok
}

// AddHandler ..
func (m *MsgHandlers) AddHandler(msg proto.Message, handler MsgHandler) {
	msgType := msg.ProtoReflect().Type()
	_, err := RegisterMessage(msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"msgtype": msgType.Descriptor().FullName(),
		}).WithError(err).Error("MsgHandlers RegisterMessage error")
		return
	}

	m.rw.Lock()
	defer m.rw.Unlock()

	if _, exist := m.handlers[msgType]; exist {
		logrus.WithFields(logrus.Fields{
			"msgtype": msgType.Descriptor().FullName(),
		}).Warning("MsgHandlers already have handler")
	}

	m.handlers[msgType] = handler
}
