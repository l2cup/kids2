// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.15.8
// source: node/node.proto

package node

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type NetworkInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	IpAddr string `protobuf:"bytes,2,opt,name=ip_addr,json=ipAddr,proto3" json:"ip_addr,omitempty"`
	Port   int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *NetworkInfo) Reset() {
	*x = NetworkInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkInfo) ProtoMessage() {}

func (x *NetworkInfo) ProtoReflect() protoreflect.Message {
	mi := &file_node_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkInfo.ProtoReflect.Descriptor instead.
func (*NetworkInfo) Descriptor() ([]byte, []int) {
	return file_node_node_proto_rawDescGZIP(), []int{0}
}

func (x *NetworkInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *NetworkInfo) GetIpAddr() string {
	if x != nil {
		return x.IpAddr
	}
	return ""
}

func (x *NetworkInfo) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type Configuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NetworkInfo    *NetworkInfo   `protobuf:"bytes,1,opt,name=network_info,json=networkInfo,proto3" json:"network_info,omitempty"`
	TotalNodeCount uint64         `protobuf:"varint,2,opt,name=TotalNodeCount,proto3" json:"TotalNodeCount,omitempty"`
	ConnectedNodes []*NetworkInfo `protobuf:"bytes,3,rep,name=connected_nodes,json=connectedNodes,proto3" json:"connected_nodes,omitempty"`
}

func (x *Configuration) Reset() {
	*x = Configuration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration) ProtoMessage() {}

func (x *Configuration) ProtoReflect() protoreflect.Message {
	mi := &file_node_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration.ProtoReflect.Descriptor instead.
func (*Configuration) Descriptor() ([]byte, []int) {
	return file_node_node_proto_rawDescGZIP(), []int{1}
}

func (x *Configuration) GetNetworkInfo() *NetworkInfo {
	if x != nil {
		return x.NetworkInfo
	}
	return nil
}

func (x *Configuration) GetTotalNodeCount() uint64 {
	if x != nil {
		return x.TotalNodeCount
	}
	return 0
}

func (x *Configuration) GetConnectedNodes() []*NetworkInfo {
	if x != nil {
		return x.ConnectedNodes
	}
	return nil
}

var File_node_node_proto protoreflect.FileDescriptor

var file_node_node_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6e, 0x6f, 0x64, 0x65, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x0b, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x70, 0x41, 0x64, 0x64, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74,
	0x22, 0xa9, 0x01, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x34, 0x0a, 0x0c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x26, 0x0a, 0x0e, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0e, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x3a, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x6e, 0x6f,
	0x64, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6e, 0x6f, 0x64, 0x65,
	0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0e, 0x63, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x4e, 0x6f, 0x64, 0x65, 0x73, 0x32, 0x39, 0x0a, 0x04,
	0x4e, 0x6f, 0x64, 0x65, 0x12, 0x31, 0x0a, 0x09, 0x49, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x65, 0x12, 0x11, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x11, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x32, 0x44, 0x0a, 0x09, 0x42, 0x6f, 0x6f, 0x74, 0x73,
	0x74, 0x72, 0x61, 0x70, 0x12, 0x37, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x6e, 0x6f, 0x64, 0x65, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x2c, 0x5a,
	0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x32, 0x63, 0x75,
	0x70, 0x2f, 0x6b, 0x69, 0x64, 0x73, 0x32, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_node_node_proto_rawDescOnce sync.Once
	file_node_node_proto_rawDescData = file_node_node_proto_rawDesc
)

func file_node_node_proto_rawDescGZIP() []byte {
	file_node_node_proto_rawDescOnce.Do(func() {
		file_node_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_node_proto_rawDescData)
	})
	return file_node_node_proto_rawDescData
}

var file_node_node_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_node_node_proto_goTypes = []interface{}{
	(*NetworkInfo)(nil),   // 0: node.NetworkInfo
	(*Configuration)(nil), // 1: node.Configuration
	(*emptypb.Empty)(nil), // 2: google.protobuf.Empty
}
var file_node_node_proto_depIdxs = []int32{
	0, // 0: node.Configuration.network_info:type_name -> node.NetworkInfo
	0, // 1: node.Configuration.connected_nodes:type_name -> node.NetworkInfo
	0, // 2: node.Node.Introduce:input_type -> node.NetworkInfo
	2, // 3: node.Bootstrap.Register:input_type -> google.protobuf.Empty
	0, // 4: node.Node.Introduce:output_type -> node.NetworkInfo
	1, // 5: node.Bootstrap.Register:output_type -> node.Configuration
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_node_node_proto_init() }
func file_node_node_proto_init() {
	if File_node_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_node_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkInfo); i {
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
		file_node_node_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration); i {
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
			RawDescriptor: file_node_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_node_node_proto_goTypes,
		DependencyIndexes: file_node_node_proto_depIdxs,
		MessageInfos:      file_node_node_proto_msgTypes,
	}.Build()
	File_node_node_proto = out.File
	file_node_node_proto_rawDesc = nil
	file_node_node_proto_goTypes = nil
	file_node_node_proto_depIdxs = nil
}
