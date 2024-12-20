// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.19.3
// source: frame/frame.proto

package frame

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GameState int32

const (
	GameState_GameState_None  GameState = 0
	GameState_GameState_Start GameState = 1
)

// Enum value maps for GameState.
var (
	GameState_name = map[int32]string{
		0: "GameState_None",
		1: "GameState_Start",
	}
	GameState_value = map[string]int32{
		"GameState_None":  0,
		"GameState_Start": 1,
	}
)

func (x GameState) Enum() *GameState {
	p := new(GameState)
	*p = x
	return p
}

func (x GameState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameState) Descriptor() protoreflect.EnumDescriptor {
	return file_frame_frame_proto_enumTypes[0].Descriptor()
}

func (GameState) Type() protoreflect.EnumType {
	return &file_frame_frame_proto_enumTypes[0]
}

func (x GameState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameState.Descriptor instead.
func (GameState) EnumDescriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{0}
}

type ReqSyncFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Frame *SyncFrame `protobuf:"bytes,1,opt,name=frame,proto3" json:"frame,omitempty"`
}

func (x *ReqSyncFrame) Reset() {
	*x = ReqSyncFrame{}
	mi := &file_frame_frame_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReqSyncFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqSyncFrame) ProtoMessage() {}

func (x *ReqSyncFrame) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqSyncFrame.ProtoReflect.Descriptor instead.
func (*ReqSyncFrame) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{0}
}

func (x *ReqSyncFrame) GetFrame() *SyncFrame {
	if x != nil {
		return x.Frame
	}
	return nil
}

type SyncFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerID int64   `protobuf:"varint,1,opt,name=PlayerID,proto3" json:"PlayerID,omitempty"`
	Frame    int64   `protobuf:"varint,2,opt,name=Frame,proto3" json:"Frame,omitempty"`
	X        float32 `protobuf:"fixed32,3,opt,name=X,proto3" json:"X,omitempty"`
	Y        float32 `protobuf:"fixed32,4,opt,name=Y,proto3" json:"Y,omitempty"`
}

func (x *SyncFrame) Reset() {
	*x = SyncFrame{}
	mi := &file_frame_frame_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncFrame) ProtoMessage() {}

func (x *SyncFrame) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncFrame.ProtoReflect.Descriptor instead.
func (*SyncFrame) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{1}
}

func (x *SyncFrame) GetPlayerID() int64 {
	if x != nil {
		return x.PlayerID
	}
	return 0
}

func (x *SyncFrame) GetFrame() int64 {
	if x != nil {
		return x.Frame
	}
	return 0
}

func (x *SyncFrame) GetX() float32 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *SyncFrame) GetY() float32 {
	if x != nil {
		return x.Y
	}
	return 0
}

type RspSyncFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerFrame []*SyncFrame `protobuf:"bytes,1,rep,name=serverFrame,proto3" json:"serverFrame,omitempty"`
}

func (x *RspSyncFrame) Reset() {
	*x = RspSyncFrame{}
	mi := &file_frame_frame_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspSyncFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspSyncFrame) ProtoMessage() {}

func (x *RspSyncFrame) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspSyncFrame.ProtoReflect.Descriptor instead.
func (*RspSyncFrame) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{2}
}

func (x *RspSyncFrame) GetServerFrame() []*SyncFrame {
	if x != nil {
		return x.ServerFrame
	}
	return nil
}

type ReqReadyBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerID int64 `protobuf:"varint,1,opt,name=PlayerID,proto3" json:"PlayerID,omitempty"`
}

func (x *ReqReadyBattle) Reset() {
	*x = ReqReadyBattle{}
	mi := &file_frame_frame_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReqReadyBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqReadyBattle) ProtoMessage() {}

func (x *ReqReadyBattle) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqReadyBattle.ProtoReflect.Descriptor instead.
func (*ReqReadyBattle) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{3}
}

func (x *ReqReadyBattle) GetPlayerID() int64 {
	if x != nil {
		return x.PlayerID
	}
	return 0
}

type RspReadyBattle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerID int64 `protobuf:"varint,1,opt,name=PlayerID,proto3" json:"PlayerID,omitempty"`
}

func (x *RspReadyBattle) Reset() {
	*x = RspReadyBattle{}
	mi := &file_frame_frame_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspReadyBattle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspReadyBattle) ProtoMessage() {}

func (x *RspReadyBattle) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspReadyBattle.ProtoReflect.Descriptor instead.
func (*RspReadyBattle) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{4}
}

func (x *RspReadyBattle) GetPlayerID() int64 {
	if x != nil {
		return x.PlayerID
	}
	return 0
}

type RspNotifyGameStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RspNotifyGameStart) Reset() {
	*x = RspNotifyGameStart{}
	mi := &file_frame_frame_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RspNotifyGameStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RspNotifyGameStart) ProtoMessage() {}

func (x *RspNotifyGameStart) ProtoReflect() protoreflect.Message {
	mi := &file_frame_frame_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RspNotifyGameStart.ProtoReflect.Descriptor instead.
func (*RspNotifyGameStart) Descriptor() ([]byte, []int) {
	return file_frame_frame_proto_rawDescGZIP(), []int{5}
}

var File_frame_frame_proto protoreflect.FileDescriptor

var file_frame_frame_proto_rawDesc = []byte{
	0x0a, 0x11, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x58, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x22,
	0x3b, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x53, 0x79, 0x6e, 0x63, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12,
	0x2b, 0x0a, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x58, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x53, 0x79, 0x6e, 0x63,
	0x46, 0x72, 0x61, 0x6d, 0x65, 0x52, 0x05, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x22, 0x59, 0x0a, 0x09,
	0x53, 0x79, 0x6e, 0x63, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x58,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x58, 0x12, 0x0c, 0x0a, 0x01, 0x59, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x01, 0x59, 0x22, 0x47, 0x0a, 0x0c, 0x52, 0x73, 0x70, 0x53, 0x79,
	0x6e, 0x63, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x58,
	0x46, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x46, 0x72,
	0x61, 0x6d, 0x65, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x46, 0x72, 0x61, 0x6d, 0x65,
	0x22, 0x2c, 0x0a, 0x0e, 0x52, 0x65, 0x71, 0x52, 0x65, 0x61, 0x64, 0x79, 0x42, 0x61, 0x74, 0x74,
	0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x22, 0x2c,
	0x0a, 0x0e, 0x52, 0x73, 0x70, 0x52, 0x65, 0x61, 0x64, 0x79, 0x42, 0x61, 0x74, 0x74, 0x6c, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x22, 0x14, 0x0a, 0x12,
	0x52, 0x73, 0x70, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x2a, 0x34, 0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x12, 0x0a, 0x0e, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x4e, 0x6f, 0x6e,
	0x65, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x5f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x10, 0x01, 0x42, 0x13, 0x5a, 0x11, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_frame_frame_proto_rawDescOnce sync.Once
	file_frame_frame_proto_rawDescData = file_frame_frame_proto_rawDesc
)

func file_frame_frame_proto_rawDescGZIP() []byte {
	file_frame_frame_proto_rawDescOnce.Do(func() {
		file_frame_frame_proto_rawDescData = protoimpl.X.CompressGZIP(file_frame_frame_proto_rawDescData)
	})
	return file_frame_frame_proto_rawDescData
}

var file_frame_frame_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_frame_frame_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_frame_frame_proto_goTypes = []any{
	(GameState)(0),             // 0: XFramework.GameState
	(*ReqSyncFrame)(nil),       // 1: XFramework.ReqSyncFrame
	(*SyncFrame)(nil),          // 2: XFramework.SyncFrame
	(*RspSyncFrame)(nil),       // 3: XFramework.RspSyncFrame
	(*ReqReadyBattle)(nil),     // 4: XFramework.ReqReadyBattle
	(*RspReadyBattle)(nil),     // 5: XFramework.RspReadyBattle
	(*RspNotifyGameStart)(nil), // 6: XFramework.RspNotifyGameStart
}
var file_frame_frame_proto_depIdxs = []int32{
	2, // 0: XFramework.ReqSyncFrame.frame:type_name -> XFramework.SyncFrame
	2, // 1: XFramework.RspSyncFrame.serverFrame:type_name -> XFramework.SyncFrame
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_frame_frame_proto_init() }
func file_frame_frame_proto_init() {
	if File_frame_frame_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_frame_frame_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_frame_frame_proto_goTypes,
		DependencyIndexes: file_frame_frame_proto_depIdxs,
		EnumInfos:         file_frame_frame_proto_enumTypes,
		MessageInfos:      file_frame_frame_proto_msgTypes,
	}.Build()
	File_frame_frame_proto = out.File
	file_frame_frame_proto_rawDesc = nil
	file_frame_frame_proto_goTypes = nil
	file_frame_frame_proto_depIdxs = nil
}
