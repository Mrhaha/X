package msgprocessor

import (
	"fmt"
	"sync"

	"XServer/framework/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Msg interface {
	Marshal() (ext proto.Message, msgid uint32, data []byte, err error)
}

type BytesMsg struct {
	ExtMsg proto.Message
	MsgID  uint32
	Data   []byte
}

func (m *BytesMsg) Marshal() (proto.Message, uint32, []byte, error) {
	return m.ExtMsg, m.MsgID, m.Data, nil
}

type ProtoMsg struct {
	ExtMsg proto.Message
	Msg    proto.Message
}

func (m *ProtoMsg) Marshal() (proto.Message, uint32, []byte, error) {
	msgID, data, err := OnMarshal(m.Msg)
	return m.ExtMsg, msgID, data, err
}

//------------消息处理-------------------------
//
// 关联id与msg结构
//
var (
	msgMutex   sync.Mutex
	msgID2Typ  map[uint32]protoreflect.MessageType
	msgType2ID map[protoreflect.MessageType]uint32
)

func init() {
	msgID2Typ = make(map[uint32]protoreflect.MessageType)
	msgType2ID = make(map[protoreflect.MessageType]uint32)
}

// RegisterMessage ...
func RegisterMessage(msg proto.Message) (uint32, error) {
	id := util.MessageHash(msg)
	msgType := msg.ProtoReflect().Type()

	msgMutex.Lock()
	defer msgMutex.Unlock()

	if usedMsgType, bHave := msgID2Typ[id]; bHave {
		if usedMsgType != msgType {
			return id, fmt.Errorf("register message, message id:%v-type:%v, typein:%v is already registered", id, usedMsgType, msgType)
		}

		return id, nil
	}

	msgID2Typ[id] = msgType
	msgType2ID[msgType] = id

	logrus.WithFields(logrus.Fields{
		"msgid":   id,
		"msgtype": msgType.Descriptor().FullName(),
	}).Debug("RegisterMessage")

	return id, nil
}

// MessageType ...
func MessageType(id uint32) (protoreflect.MessageType, bool) {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	if msgType, bHave := msgID2Typ[id]; bHave {
		return msgType, true
	}

	return nil, false
}

// MessageID ...
func MessageID(msgType protoreflect.MessageType) (uint32, bool) {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	if msgID, bHave := msgType2ID[msgType]; bHave {
		return msgID, true
	}

	return 0, false
}

// OnUnmarshal on read gorutine.
func OnUnmarshal(id uint32, data []byte) (proto.Message, error) {
	if msgType, bHave := MessageType(id); bHave {
		msg := msgType.New().Interface()
		err := proto.Unmarshal(data, msg)
		return msg.(proto.Message), err
	}
	return nil, fmt.Errorf("OnUnmarshal, message %d not registered", id)
}

// OnMarshal on write gorutine.
func OnMarshal(msg proto.Message) (uint32, []byte, error) {
	msgType := msg.ProtoReflect().Type()
	if msgID, bHave := MessageID(msgType); bHave {
		data, err := proto.Marshal(msg.(proto.Message))
		return msgID, data, err
	}
	if msgID, e := RegisterMessage(msg); e == nil {
		data, err := proto.Marshal(msg.(proto.Message))
		return msgID, data, err
	}
	return 0, nil, fmt.Errorf("OnMarshal, message %s auto register failed", msgType.Descriptor().FullName())
}
