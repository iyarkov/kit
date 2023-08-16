// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: iyarkov/kit/hlc.proto

package protobuf

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

type Stamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time  int32 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Count int32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	Node  int32 `protobuf:"varint,3,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *Stamp) Reset() {
	*x = Stamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iyarkov_kit_hlc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stamp) ProtoMessage() {}

func (x *Stamp) ProtoReflect() protoreflect.Message {
	mi := &file_iyarkov_kit_hlc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stamp.ProtoReflect.Descriptor instead.
func (*Stamp) Descriptor() ([]byte, []int) {
	return file_iyarkov_kit_hlc_proto_rawDescGZIP(), []int{0}
}

func (x *Stamp) GetTime() int32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *Stamp) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Stamp) GetNode() int32 {
	if x != nil {
		return x.Node
	}
	return 0
}

var File_iyarkov_kit_hlc_proto protoreflect.FileDescriptor

var file_iyarkov_kit_hlc_proto_rawDesc = []byte{
	0x0a, 0x15, 0x69, 0x79, 0x61, 0x72, 0x6b, 0x6f, 0x76, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x68, 0x6c,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x79, 0x61, 0x72, 0x6b, 0x6f, 0x76,
	0x2e, 0x6b, 0x69, 0x74, 0x22, 0x45, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x42, 0x25, 0x5a, 0x23, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x79, 0x61, 0x72, 0x6b, 0x6f,
	0x76, 0x2f, 0x6b, 0x69, 0x74, 0x2f, 0x68, 0x6c, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iyarkov_kit_hlc_proto_rawDescOnce sync.Once
	file_iyarkov_kit_hlc_proto_rawDescData = file_iyarkov_kit_hlc_proto_rawDesc
)

func file_iyarkov_kit_hlc_proto_rawDescGZIP() []byte {
	file_iyarkov_kit_hlc_proto_rawDescOnce.Do(func() {
		file_iyarkov_kit_hlc_proto_rawDescData = protoimpl.X.CompressGZIP(file_iyarkov_kit_hlc_proto_rawDescData)
	})
	return file_iyarkov_kit_hlc_proto_rawDescData
}

var file_iyarkov_kit_hlc_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_iyarkov_kit_hlc_proto_goTypes = []interface{}{
	(*Stamp)(nil), // 0: iyarkov.kit.Stamp
}
var file_iyarkov_kit_hlc_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_iyarkov_kit_hlc_proto_init() }
func file_iyarkov_kit_hlc_proto_init() {
	if File_iyarkov_kit_hlc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_iyarkov_kit_hlc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stamp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_iyarkov_kit_hlc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iyarkov_kit_hlc_proto_goTypes,
		DependencyIndexes: file_iyarkov_kit_hlc_proto_depIdxs,
		MessageInfos:      file_iyarkov_kit_hlc_proto_msgTypes,
	}.Build()
	File_iyarkov_kit_hlc_proto = out.File
	file_iyarkov_kit_hlc_proto_rawDesc = nil
	file_iyarkov_kit_hlc_proto_goTypes = nil
	file_iyarkov_kit_hlc_proto_depIdxs = nil
}
