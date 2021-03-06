// Code generated by protoc-gen-go. DO NOT EDIT.
// source: servicectl.proto

package pb

import (
    context "context"
    fmt "fmt"
    proto "github.com/golang/protobuf/proto"
    grpc "google.golang.org/grpc"
    codes "google.golang.org/grpc/codes"
    status "google.golang.org/grpc/status"
    math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// 服务监听进程应答消息
type ServiceCntrolReplyOnly struct {
    StatusCode           int32    `protobuf:"varint,1,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
    Reason               string   `protobuf:"bytes,2,opt,name=Reason,proto3" json:"Reason,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceCntrolReplyOnly) Reset()         { *m = ServiceCntrolReplyOnly{} }
func (m *ServiceCntrolReplyOnly) String() string { return proto.CompactTextString(m) }
func (*ServiceCntrolReplyOnly) ProtoMessage()    {}
func (*ServiceCntrolReplyOnly) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{0}
}

func (m *ServiceCntrolReplyOnly) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ServiceCntrolReplyOnly.Unmarshal(m, b)
}
func (m *ServiceCntrolReplyOnly) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ServiceCntrolReplyOnly.Marshal(b, m, deterministic)
}
func (m *ServiceCntrolReplyOnly) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ServiceCntrolReplyOnly.Merge(m, src)
}
func (m *ServiceCntrolReplyOnly) XXX_Size() int {
    return xxx_messageInfo_ServiceCntrolReplyOnly.Size(m)
}
func (m *ServiceCntrolReplyOnly) XXX_DiscardUnknown() {
    xxx_messageInfo_ServiceCntrolReplyOnly.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceCntrolReplyOnly proto.InternalMessageInfo

func (m *ServiceCntrolReplyOnly) GetStatusCode() int32 {
    if m != nil {
        return m.StatusCode
    }
    return 0
}

func (m *ServiceCntrolReplyOnly) GetReason() string {
    if m != nil {
        return m.Reason
    }
    return ""
}

type ProxyMapAddRequest struct {
    Creator              string   `protobuf:"bytes,1,opt,name=Creator,proto3" json:"Creator,omitempty"`
    ListenPort           uint32   `protobuf:"varint,2,opt,name=ListenPort,proto3" json:"ListenPort,omitempty"`
    IP                   string   `protobuf:"bytes,3,opt,name=IP,proto3" json:"IP,omitempty"`
    Port                 uint32   `protobuf:"varint,4,opt,name=Port,proto3" json:"Port,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyMapAddRequest) Reset()         { *m = ProxyMapAddRequest{} }
func (m *ProxyMapAddRequest) String() string { return proto.CompactTextString(m) }
func (*ProxyMapAddRequest) ProtoMessage()    {}
func (*ProxyMapAddRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{1}
}

func (m *ProxyMapAddRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ProxyMapAddRequest.Unmarshal(m, b)
}
func (m *ProxyMapAddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ProxyMapAddRequest.Marshal(b, m, deterministic)
}
func (m *ProxyMapAddRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ProxyMapAddRequest.Merge(m, src)
}
func (m *ProxyMapAddRequest) XXX_Size() int {
    return xxx_messageInfo_ProxyMapAddRequest.Size(m)
}
func (m *ProxyMapAddRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ProxyMapAddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyMapAddRequest proto.InternalMessageInfo

func (m *ProxyMapAddRequest) GetCreator() string {
    if m != nil {
        return m.Creator
    }
    return ""
}

func (m *ProxyMapAddRequest) GetListenPort() uint32 {
    if m != nil {
        return m.ListenPort
    }
    return 0
}

func (m *ProxyMapAddRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

func (m *ProxyMapAddRequest) GetPort() uint32 {
    if m != nil {
        return m.Port
    }
    return 0
}

type ProxyMapDelRequest struct {
    ListenPort           uint32   `protobuf:"varint,1,opt,name=ListenPort,proto3" json:"ListenPort,omitempty"`
    IP                   string   `protobuf:"bytes,2,opt,name=IP,proto3" json:"IP,omitempty"`
    Port                 uint32   `protobuf:"varint,3,opt,name=Port,proto3" json:"Port,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyMapDelRequest) Reset()         { *m = ProxyMapDelRequest{} }
func (m *ProxyMapDelRequest) String() string { return proto.CompactTextString(m) }
func (*ProxyMapDelRequest) ProtoMessage()    {}
func (*ProxyMapDelRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{2}
}

func (m *ProxyMapDelRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ProxyMapDelRequest.Unmarshal(m, b)
}
func (m *ProxyMapDelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ProxyMapDelRequest.Marshal(b, m, deterministic)
}
func (m *ProxyMapDelRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ProxyMapDelRequest.Merge(m, src)
}
func (m *ProxyMapDelRequest) XXX_Size() int {
    return xxx_messageInfo_ProxyMapDelRequest.Size(m)
}
func (m *ProxyMapDelRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ProxyMapDelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyMapDelRequest proto.InternalMessageInfo

func (m *ProxyMapDelRequest) GetListenPort() uint32 {
    if m != nil {
        return m.ListenPort
    }
    return 0
}

func (m *ProxyMapDelRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

func (m *ProxyMapDelRequest) GetPort() uint32 {
    if m != nil {
        return m.Port
    }
    return 0
}

type ProxyMapModifyRequest struct {
    ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
    ListenPort           uint32   `protobuf:"varint,2,opt,name=ListenPort,proto3" json:"ListenPort,omitempty"`
    IP                   string   `protobuf:"bytes,3,opt,name=IP,proto3" json:"IP,omitempty"`
    Port                 uint32   `protobuf:"varint,4,opt,name=Port,proto3" json:"Port,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyMapModifyRequest) Reset()         { *m = ProxyMapModifyRequest{} }
func (m *ProxyMapModifyRequest) String() string { return proto.CompactTextString(m) }
func (*ProxyMapModifyRequest) ProtoMessage()    {}
func (*ProxyMapModifyRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{3}
}

func (m *ProxyMapModifyRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ProxyMapModifyRequest.Unmarshal(m, b)
}
func (m *ProxyMapModifyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ProxyMapModifyRequest.Marshal(b, m, deterministic)
}
func (m *ProxyMapModifyRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ProxyMapModifyRequest.Merge(m, src)
}
func (m *ProxyMapModifyRequest) XXX_Size() int {
    return xxx_messageInfo_ProxyMapModifyRequest.Size(m)
}
func (m *ProxyMapModifyRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ProxyMapModifyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyMapModifyRequest proto.InternalMessageInfo

func (m *ProxyMapModifyRequest) GetID() int64 {
    if m != nil {
        return m.ID
    }
    return 0
}

func (m *ProxyMapModifyRequest) GetListenPort() uint32 {
    if m != nil {
        return m.ListenPort
    }
    return 0
}

func (m *ProxyMapModifyRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

func (m *ProxyMapModifyRequest) GetPort() uint32 {
    if m != nil {
        return m.Port
    }
    return 0
}

// 添加监听IP请求
type ListenAddrAddRequest struct {
    IP                   string   `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
    Createor             string   `protobuf:"bytes,2,opt,name=Createor,proto3" json:"Createor,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ListenAddrAddRequest) Reset()         { *m = ListenAddrAddRequest{} }
func (m *ListenAddrAddRequest) String() string { return proto.CompactTextString(m) }
func (*ListenAddrAddRequest) ProtoMessage()    {}
func (*ListenAddrAddRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{4}
}

func (m *ListenAddrAddRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ListenAddrAddRequest.Unmarshal(m, b)
}
func (m *ListenAddrAddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ListenAddrAddRequest.Marshal(b, m, deterministic)
}
func (m *ListenAddrAddRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ListenAddrAddRequest.Merge(m, src)
}
func (m *ListenAddrAddRequest) XXX_Size() int {
    return xxx_messageInfo_ListenAddrAddRequest.Size(m)
}
func (m *ListenAddrAddRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ListenAddrAddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListenAddrAddRequest proto.InternalMessageInfo

func (m *ListenAddrAddRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

func (m *ListenAddrAddRequest) GetCreateor() string {
    if m != nil {
        return m.Createor
    }
    return ""
}

// 删除监听IP请求
type ListenAddrDelRequest struct {
    IP                   string   `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ListenAddrDelRequest) Reset()         { *m = ListenAddrDelRequest{} }
func (m *ListenAddrDelRequest) String() string { return proto.CompactTextString(m) }
func (*ListenAddrDelRequest) ProtoMessage()    {}
func (*ListenAddrDelRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{5}
}

func (m *ListenAddrDelRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ListenAddrDelRequest.Unmarshal(m, b)
}
func (m *ListenAddrDelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ListenAddrDelRequest.Marshal(b, m, deterministic)
}
func (m *ListenAddrDelRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ListenAddrDelRequest.Merge(m, src)
}
func (m *ListenAddrDelRequest) XXX_Size() int {
    return xxx_messageInfo_ListenAddrDelRequest.Size(m)
}
func (m *ListenAddrDelRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ListenAddrDelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListenAddrDelRequest proto.InternalMessageInfo

func (m *ListenAddrDelRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

// 修改监听IP请求
type ListenAddrModifyRequest struct {
    ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
    IP                   string   `protobuf:"bytes,2,opt,name=IP,proto3" json:"IP,omitempty"`
    XXX_NoUnkeyedLiteral struct{} `json:"-"`
    XXX_unrecognized     []byte   `json:"-"`
    XXX_sizecache        int32    `json:"-"`
}

func (m *ListenAddrModifyRequest) Reset()         { *m = ListenAddrModifyRequest{} }
func (m *ListenAddrModifyRequest) String() string { return proto.CompactTextString(m) }
func (*ListenAddrModifyRequest) ProtoMessage()    {}
func (*ListenAddrModifyRequest) Descriptor() ([]byte, []int) {
    return fileDescriptor_5561ec5be41fd2ff, []int{6}
}

func (m *ListenAddrModifyRequest) XXX_Unmarshal(b []byte) error {
    return xxx_messageInfo_ListenAddrModifyRequest.Unmarshal(m, b)
}
func (m *ListenAddrModifyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
    return xxx_messageInfo_ListenAddrModifyRequest.Marshal(b, m, deterministic)
}
func (m *ListenAddrModifyRequest) XXX_Merge(src proto.Message) {
    xxx_messageInfo_ListenAddrModifyRequest.Merge(m, src)
}
func (m *ListenAddrModifyRequest) XXX_Size() int {
    return xxx_messageInfo_ListenAddrModifyRequest.Size(m)
}
func (m *ListenAddrModifyRequest) XXX_DiscardUnknown() {
    xxx_messageInfo_ListenAddrModifyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListenAddrModifyRequest proto.InternalMessageInfo

func (m *ListenAddrModifyRequest) GetID() int64 {
    if m != nil {
        return m.ID
    }
    return 0
}

func (m *ListenAddrModifyRequest) GetIP() string {
    if m != nil {
        return m.IP
    }
    return ""
}

func init() {
    proto.RegisterType((*ServiceCntrolReplyOnly)(nil), "pb.ServiceCntrolReplyOnly")
    proto.RegisterType((*ProxyMapAddRequest)(nil), "pb.ProxyMapAddRequest")
    proto.RegisterType((*ProxyMapDelRequest)(nil), "pb.ProxyMapDelRequest")
    proto.RegisterType((*ProxyMapModifyRequest)(nil), "pb.ProxyMapModifyRequest")
    proto.RegisterType((*ListenAddrAddRequest)(nil), "pb.ListenAddrAddRequest")
    proto.RegisterType((*ListenAddrDelRequest)(nil), "pb.ListenAddrDelRequest")
    proto.RegisterType((*ListenAddrModifyRequest)(nil), "pb.ListenAddrModifyRequest")
}

func init() {
    proto.RegisterFile("servicectl.proto", fileDescriptor_5561ec5be41fd2ff)
}

var fileDescriptor_5561ec5be41fd2ff = []byte{
    // 366 bytes of a gzipped FileDescriptorProto
    0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x41, 0x4f, 0xc2, 0x40,
    0x10, 0x85, 0xd3, 0x16, 0x51, 0xc6, 0x40, 0xc8, 0x44, 0xb1, 0x62, 0x62, 0x48, 0x0f, 0x86, 0x13,
    0x07, 0x3d, 0x79, 0x44, 0x6a, 0x0c, 0x51, 0x62, 0xb3, 0x5c, 0xbc, 0x16, 0xba, 0x26, 0xc4, 0x4d,
    0xb7, 0x6e, 0x17, 0x63, 0xff, 0x9a, 0xbf, 0xce, 0x74, 0xa1, 0xb0, 0x5b, 0x84, 0xf4, 0xe0, 0x8d,
    0x19, 0xe6, 0x7d, 0xb3, 0x6f, 0xf2, 0x0a, 0xed, 0x94, 0x8a, 0xaf, 0xc5, 0x9c, 0xce, 0x25, 0x1b,
    0x24, 0x82, 0x4b, 0x8e, 0x76, 0x32, 0xf3, 0x02, 0xe8, 0x4c, 0x57, 0xfd, 0x51, 0x2c, 0x05, 0x67,
    0x84, 0x26, 0x2c, 0x7b, 0x8d, 0x59, 0x86, 0xd7, 0x00, 0x53, 0x19, 0xca, 0x65, 0x3a, 0xe2, 0x11,
    0x75, 0xad, 0x9e, 0xd5, 0x3f, 0x22, 0x5a, 0x07, 0x3b, 0x50, 0x27, 0x34, 0x4c, 0x79, 0xec, 0xda,
    0x3d, 0xab, 0xdf, 0x20, 0xeb, 0xca, 0x13, 0x80, 0x81, 0xe0, 0xdf, 0xd9, 0x24, 0x4c, 0x86, 0x51,
    0x44, 0xe8, 0xe7, 0x92, 0xa6, 0x12, 0x5d, 0x38, 0x1e, 0x09, 0x1a, 0x4a, 0x2e, 0x14, 0xaa, 0x41,
    0x8a, 0x32, 0xdf, 0xf3, 0xb2, 0x48, 0x25, 0x8d, 0x03, 0x2e, 0xa4, 0x62, 0x35, 0x89, 0xd6, 0xc1,
    0x16, 0xd8, 0xe3, 0xc0, 0x75, 0x94, 0xc8, 0x1e, 0x07, 0x88, 0x50, 0x53, 0x93, 0x35, 0x35, 0xa9,
    0x7e, 0x7b, 0x6f, 0xdb, 0x9d, 0x3e, 0x65, 0xc5, 0x4e, 0x93, 0x6c, 0xed, 0x21, 0xdb, 0x3b, 0x64,
    0x47, 0x23, 0x7f, 0xc0, 0x79, 0x41, 0x9e, 0xf0, 0x68, 0xf1, 0x9e, 0x15, 0xf0, 0x5c, 0xec, 0x2b,
    0xa8, 0x43, 0xec, 0xb1, 0xff, 0x2f, 0x36, 0x1e, 0xe0, 0x6c, 0xa5, 0x18, 0x46, 0x91, 0xd0, 0x8e,
    0xb7, 0xd2, 0x5a, 0x1b, 0x6d, 0x17, 0x4e, 0xd4, 0xf5, 0x28, 0x17, 0xeb, 0xe7, 0x6f, 0x6a, 0xef,
    0x46, 0x67, 0x68, 0xc7, 0x28, 0x31, 0xbc, 0x7b, 0xb8, 0xd8, 0xce, 0x1d, 0xb6, 0x56, 0xba, 0xd3,
    0xed, 0x8f, 0x03, 0xad, 0x22, 0x34, 0x5c, 0xa5, 0x06, 0x1f, 0xa1, 0x69, 0xbc, 0x1c, 0xdd, 0x41,
    0x32, 0x1b, 0xfc, 0x65, 0xa6, 0xdb, 0xcd, 0xff, 0xd9, 0x93, 0x39, 0x03, 0xe3, 0x53, 0x56, 0xc6,
    0x6c, 0xfd, 0x1c, 0xc4, 0x3c, 0x43, 0xbb, 0xec, 0x0d, 0xaf, 0x4c, 0x92, 0xe1, 0xf8, 0x20, 0x6c,
    0x08, 0xa7, 0x5a, 0x9e, 0xb1, 0x93, 0x8f, 0xee, 0x06, 0xbc, 0x2a, 0x22, 0x37, 0x65, 0x20, 0x2a,
    0x5a, 0x7a, 0x82, 0x96, 0x99, 0x43, 0xbc, 0xd4, 0x29, 0x95, 0xed, 0xcc, 0xea, 0xea, 0xdb, 0xbf,
    0xfb, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xe2, 0x2c, 0xeb, 0x82, 0x0f, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServiceControlClient is the client API for ServiceControl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceControlClient interface {
    ListenAddrAdd(ctx context.Context, in *ListenAddrAddRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
    ListenAddrDel(ctx context.Context, in *ListenAddrDelRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
    ListenAddrModify(ctx context.Context, in *ListenAddrModifyRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
    ProxyMapAdd(ctx context.Context, in *ProxyMapAddRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
    ProxyMapDel(ctx context.Context, in *ProxyMapDelRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
    ProxyMapModify(ctx context.Context, in *ProxyMapModifyRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error)
}

type serviceControlClient struct {
    cc grpc.ClientConnInterface
}

func NewServiceControlClient(cc grpc.ClientConnInterface) ServiceControlClient {
    return &serviceControlClient{cc}
}

func (c *serviceControlClient) ListenAddrAdd(ctx context.Context, in *ListenAddrAddRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ListenAddrAdd", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *serviceControlClient) ListenAddrDel(ctx context.Context, in *ListenAddrDelRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ListenAddrDel", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *serviceControlClient) ListenAddrModify(ctx context.Context, in *ListenAddrModifyRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ListenAddrModify", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *serviceControlClient) ProxyMapAdd(ctx context.Context, in *ProxyMapAddRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ProxyMapAdd", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *serviceControlClient) ProxyMapDel(ctx context.Context, in *ProxyMapDelRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ProxyMapDel", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

func (c *serviceControlClient) ProxyMapModify(ctx context.Context, in *ProxyMapModifyRequest, opts ...grpc.CallOption) (*ServiceCntrolReplyOnly, error) {
    out := new(ServiceCntrolReplyOnly)
    err := c.cc.Invoke(ctx, "/pb.ServiceControl/ProxyMapModify", in, out, opts...)
    if err != nil {
        return nil, err
    }
    return out, nil
}

// ServiceControlServer is the server API for ServiceControl service.
type ServiceControlServer interface {
    ListenAddrAdd(context.Context, *ListenAddrAddRequest) (*ServiceCntrolReplyOnly, error)
    ListenAddrDel(context.Context, *ListenAddrDelRequest) (*ServiceCntrolReplyOnly, error)
    ListenAddrModify(context.Context, *ListenAddrModifyRequest) (*ServiceCntrolReplyOnly, error)
    ProxyMapAdd(context.Context, *ProxyMapAddRequest) (*ServiceCntrolReplyOnly, error)
    ProxyMapDel(context.Context, *ProxyMapDelRequest) (*ServiceCntrolReplyOnly, error)
    ProxyMapModify(context.Context, *ProxyMapModifyRequest) (*ServiceCntrolReplyOnly, error)
}

// UnimplementedServiceControlServer can be embedded to have forward compatible implementations.
type UnimplementedServiceControlServer struct {
}

func (*UnimplementedServiceControlServer) ListenAddrAdd(ctx context.Context, req *ListenAddrAddRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ListenAddrAdd not implemented")
}
func (*UnimplementedServiceControlServer) ListenAddrDel(ctx context.Context, req *ListenAddrDelRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ListenAddrDel not implemented")
}
func (*UnimplementedServiceControlServer) ListenAddrModify(ctx context.Context, req *ListenAddrModifyRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ListenAddrModify not implemented")
}
func (*UnimplementedServiceControlServer) ProxyMapAdd(ctx context.Context, req *ProxyMapAddRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ProxyMapAdd not implemented")
}
func (*UnimplementedServiceControlServer) ProxyMapDel(ctx context.Context, req *ProxyMapDelRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ProxyMapDel not implemented")
}
func (*UnimplementedServiceControlServer) ProxyMapModify(ctx context.Context, req *ProxyMapModifyRequest) (*ServiceCntrolReplyOnly, error) {
    return nil, status.Errorf(codes.Unimplemented, "method ProxyMapModify not implemented")
}

func RegisterServiceControlServer(s *grpc.Server, srv ServiceControlServer) {
    s.RegisterService(&_ServiceControl_serviceDesc, srv)
}

func _ServiceControl_ListenAddrAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ListenAddrAddRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ListenAddrAdd(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ListenAddrAdd",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ListenAddrAdd(ctx, req.(*ListenAddrAddRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _ServiceControl_ListenAddrDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ListenAddrDelRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ListenAddrDel(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ListenAddrDel",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ListenAddrDel(ctx, req.(*ListenAddrDelRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _ServiceControl_ListenAddrModify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ListenAddrModifyRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ListenAddrModify(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ListenAddrModify",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ListenAddrModify(ctx, req.(*ListenAddrModifyRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _ServiceControl_ProxyMapAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ProxyMapAddRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ProxyMapAdd(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ProxyMapAdd",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ProxyMapAdd(ctx, req.(*ProxyMapAddRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _ServiceControl_ProxyMapDel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ProxyMapDelRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ProxyMapDel(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ProxyMapDel",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ProxyMapDel(ctx, req.(*ProxyMapDelRequest))
    }
    return interceptor(ctx, in, info, handler)
}

func _ServiceControl_ProxyMapModify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
    in := new(ProxyMapModifyRequest)
    if err := dec(in); err != nil {
        return nil, err
    }
    if interceptor == nil {
        return srv.(ServiceControlServer).ProxyMapModify(ctx, in)
    }
    info := &grpc.UnaryServerInfo{
        Server:     srv,
        FullMethod: "/pb.ServiceControl/ProxyMapModify",
    }
    handler := func(ctx context.Context, req interface{}) (interface{}, error) {
        return srv.(ServiceControlServer).ProxyMapModify(ctx, req.(*ProxyMapModifyRequest))
    }
    return interceptor(ctx, in, info, handler)
}

var _ServiceControl_serviceDesc = grpc.ServiceDesc{
    ServiceName: "pb.ServiceControl",
    HandlerType: (*ServiceControlServer)(nil),
    Methods: []grpc.MethodDesc{
        {
            MethodName: "ListenAddrAdd",
            Handler:    _ServiceControl_ListenAddrAdd_Handler,
        },
        {
            MethodName: "ListenAddrDel",
            Handler:    _ServiceControl_ListenAddrDel_Handler,
        },
        {
            MethodName: "ListenAddrModify",
            Handler:    _ServiceControl_ListenAddrModify_Handler,
        },
        {
            MethodName: "ProxyMapAdd",
            Handler:    _ServiceControl_ProxyMapAdd_Handler,
        },
        {
            MethodName: "ProxyMapDel",
            Handler:    _ServiceControl_ProxyMapDel_Handler,
        },
        {
            MethodName: "ProxyMapModify",
            Handler:    _ServiceControl_ProxyMapModify_Handler,
        },
    },
    Streams:  []grpc.StreamDesc{},
    Metadata: "servicectl.proto",
}
