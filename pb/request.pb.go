// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: request.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type RequestJoinCluster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClusterId     uint64 `protobuf:"varint,1,opt,name=ClusterId,proto3" json:"ClusterId,omitempty"`
	NodeId        uint64 `protobuf:"varint,2,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	RaftAddresses string `protobuf:"bytes,3,opt,name=RaftAddresses,proto3" json:"RaftAddresses,omitempty"`
	Key           string `protobuf:"bytes,4,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *RequestJoinCluster) Reset() {
	*x = RequestJoinCluster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestJoinCluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestJoinCluster) ProtoMessage() {}

func (x *RequestJoinCluster) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestJoinCluster.ProtoReflect.Descriptor instead.
func (*RequestJoinCluster) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{0}
}

func (x *RequestJoinCluster) GetClusterId() uint64 {
	if x != nil {
		return x.ClusterId
	}
	return 0
}

func (x *RequestJoinCluster) GetNodeId() uint64 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

func (x *RequestJoinCluster) GetRaftAddresses() string {
	if x != nil {
		return x.RaftAddresses
	}
	return ""
}

func (x *RequestJoinCluster) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type ResponseJoinResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status int64  `protobuf:"varint,1,opt,name=Status,proto3" json:"Status,omitempty"`
	Msg    string `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"` //uint64 NodeId=3;
}

func (x *ResponseJoinResult) Reset() {
	*x = ResponseJoinResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseJoinResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseJoinResult) ProtoMessage() {}

func (x *ResponseJoinResult) ProtoReflect() protoreflect.Message {
	mi := &file_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseJoinResult.ProtoReflect.Descriptor instead.
func (*ResponseJoinResult) Descriptor() ([]byte, []int) {
	return file_request_proto_rawDescGZIP(), []int{1}
}

func (x *ResponseJoinResult) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *ResponseJoinResult) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

var File_request_proto protoreflect.FileDescriptor

var file_request_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x82, 0x01, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a,
	0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x12, 0x24, 0x0a, 0x0d, 0x52, 0x61, 0x66, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x52, 0x61, 0x66, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x3e, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x32, 0x48, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a, 0x6f,
	0x69, 0x6e, 0x12, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4a,
	0x6f, 0x69, 0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_request_proto_rawDescOnce sync.Once
	file_request_proto_rawDescData = file_request_proto_rawDesc
)

func file_request_proto_rawDescGZIP() []byte {
	file_request_proto_rawDescOnce.Do(func() {
		file_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_request_proto_rawDescData)
	})
	return file_request_proto_rawDescData
}

var file_request_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_request_proto_goTypes = []interface{}{
	(*RequestJoinCluster)(nil), // 0: pb.RequestJoinCluster
	(*ResponseJoinResult)(nil), // 1: pb.ResponseJoinResult
}
var file_request_proto_depIdxs = []int32{
	0, // 0: pb.Request.RequestJoin:input_type -> pb.RequestJoinCluster
	1, // 1: pb.Request.RequestJoin:output_type -> pb.ResponseJoinResult
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_request_proto_init() }
func file_request_proto_init() {
	if File_request_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestJoinCluster); i {
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
		file_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseJoinResult); i {
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
			RawDescriptor: file_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_request_proto_goTypes,
		DependencyIndexes: file_request_proto_depIdxs,
		MessageInfos:      file_request_proto_msgTypes,
	}.Build()
	File_request_proto = out.File
	file_request_proto_rawDesc = nil
	file_request_proto_goTypes = nil
	file_request_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RequestClient is the client API for Request service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RequestClient interface {
	RequestJoin(ctx context.Context, in *RequestJoinCluster, opts ...grpc.CallOption) (*ResponseJoinResult, error)
}

type requestClient struct {
	cc grpc.ClientConnInterface
}

func NewRequestClient(cc grpc.ClientConnInterface) RequestClient {
	return &requestClient{cc}
}

func (c *requestClient) RequestJoin(ctx context.Context, in *RequestJoinCluster, opts ...grpc.CallOption) (*ResponseJoinResult, error) {
	out := new(ResponseJoinResult)
	err := c.cc.Invoke(ctx, "/pb.Request/RequestJoin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RequestServer is the server API for Request service.
type RequestServer interface {
	RequestJoin(context.Context, *RequestJoinCluster) (*ResponseJoinResult, error)
}

// UnimplementedRequestServer can be embedded to have forward compatible implementations.
type UnimplementedRequestServer struct {
}

func (*UnimplementedRequestServer) RequestJoin(context.Context, *RequestJoinCluster) (*ResponseJoinResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestJoin not implemented")
}

func RegisterRequestServer(s *grpc.Server, srv RequestServer) {
	s.RegisterService(&_Request_serviceDesc, srv)
}

func _Request_RequestJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestJoinCluster)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RequestServer).RequestJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Request/RequestJoin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RequestServer).RequestJoin(ctx, req.(*RequestJoinCluster))
	}
	return interceptor(ctx, in, info, handler)
}

var _Request_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Request",
	HandlerType: (*RequestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestJoin",
			Handler:    _Request_RequestJoin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "request.proto",
}
