package util

import (
	//protoV1 "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// GetFullNameByMessage 获取协议名称.
func GetFullNameByMessage(msg proto.Message) string {
	name := proto.MessageName(msg)
	return string(name)
}

/*// GetMessageByFullName 根据协议名称和数据反序列化V1 github.com/golang/protobuf.
func GetMessageByFullName(fullName string, data []byte) (protoV1.Message, error) {
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		return nil, err
	}
	msg := protoV1.MessageV1(msgType.New())
	err = protoV1.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}*/

// GetMessageV2ByFullName 根据协议名称和数据反序列化V2 google.golang.org/protobuf/proto.
func GetMessageV2ByFullName(fullName string, data []byte) (proto.Message, error) {
	msgName := protoreflect.FullName(fullName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if err != nil {
		return nil, err
	}
	msg := msgType.New().Interface()
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
