// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/gateway/gateway.proto

package proto_gateway_service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type CallType int32

const (
	CallType_User      CallType = 0
	CallType_Container CallType = 1
)

var CallType_name = map[int32]string{
	0: "User",
	1: "Container",
}

var CallType_value = map[string]int32{
	"User":      0,
	"Container": 1,
}

func (x CallType) String() string {
	return proto.EnumName(CallType_name, int32(x))
}

func (CallType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{0}
}

type CallerType int32

const (
	CallerType_Register           CallerType = 0
	CallerType_SendCode           CallerType = 1
	CallerType_UserLogin          CallerType = 2
	CallerType_UpdateToken        CallerType = 3
	CallerType_UserLogout         CallerType = 4
	CallerType_CreateContainer    CallerType = 5
	CallerType_GetContainerStatus CallerType = 6
	CallerType_DeleteContainer    CallerType = 7
	CallerType_GetImageList       CallerType = 8
)

var CallerType_name = map[int32]string{
	0: "Register",
	1: "SendCode",
	2: "UserLogin",
	3: "UpdateToken",
	4: "UserLogout",
	5: "CreateContainer",
	6: "GetContainerStatus",
	7: "DeleteContainer",
	8: "GetImageList",
}

var CallerType_value = map[string]int32{
	"Register":           0,
	"SendCode":           1,
	"UserLogin":          2,
	"UpdateToken":        3,
	"UserLogout":         4,
	"CreateContainer":    5,
	"GetContainerStatus": 6,
	"DeleteContainer":    7,
	"GetImageList":       8,
}

func (x CallerType) String() string {
	return proto.EnumName(CallerType_name, int32(x))
}

func (CallerType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{1}
}

type UserOpt struct {
	Username             string   `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Phone                string   `protobuf:"bytes,2,opt,name=Phone,proto3" json:"Phone,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=Password,proto3" json:"Password,omitempty"`
	Verify               string   `protobuf:"bytes,4,opt,name=Verify,proto3" json:"Verify,omitempty"`
	Token                string   `protobuf:"bytes,5,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserOpt) Reset()         { *m = UserOpt{} }
func (m *UserOpt) String() string { return proto.CompactTextString(m) }
func (*UserOpt) ProtoMessage()    {}
func (*UserOpt) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{0}
}

func (m *UserOpt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserOpt.Unmarshal(m, b)
}
func (m *UserOpt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserOpt.Marshal(b, m, deterministic)
}
func (m *UserOpt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserOpt.Merge(m, src)
}
func (m *UserOpt) XXX_Size() int {
	return xxx_messageInfo_UserOpt.Size(m)
}
func (m *UserOpt) XXX_DiscardUnknown() {
	xxx_messageInfo_UserOpt.DiscardUnknown(m)
}

var xxx_messageInfo_UserOpt proto.InternalMessageInfo

func (m *UserOpt) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserOpt) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *UserOpt) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *UserOpt) GetVerify() string {
	if m != nil {
		return m.Verify
	}
	return ""
}

func (m *UserOpt) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Options struct {
	User *UserOpt `protobuf:"bytes,1,opt,name=User,proto3" json:"User,omitempty"`
	//proto.cnode.service.CreateOpt Create = 2;
	Create               *ImageInfo `protobuf:"bytes,2,opt,name=Create,proto3" json:"Create,omitempty"`
	Cid                  string     `protobuf:"bytes,3,opt,name=Cid,proto3" json:"Cid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Options) Reset()         { *m = Options{} }
func (m *Options) String() string { return proto.CompactTextString(m) }
func (*Options) ProtoMessage()    {}
func (*Options) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{1}
}

func (m *Options) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Options.Unmarshal(m, b)
}
func (m *Options) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Options.Marshal(b, m, deterministic)
}
func (m *Options) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Options.Merge(m, src)
}
func (m *Options) XXX_Size() int {
	return xxx_messageInfo_Options.Size(m)
}
func (m *Options) XXX_DiscardUnknown() {
	xxx_messageInfo_Options.DiscardUnknown(m)
}

var xxx_messageInfo_Options proto.InternalMessageInfo

func (m *Options) GetUser() *UserOpt {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *Options) GetCreate() *ImageInfo {
	if m != nil {
		return m.Create
	}
	return nil
}

func (m *Options) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

type ImageInfo struct {
	ID                   uint32   `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Cid                  string   `protobuf:"bytes,2,opt,name=Cid,proto3" json:"Cid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageInfo) Reset()         { *m = ImageInfo{} }
func (m *ImageInfo) String() string { return proto.CompactTextString(m) }
func (*ImageInfo) ProtoMessage()    {}
func (*ImageInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{2}
}

func (m *ImageInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageInfo.Unmarshal(m, b)
}
func (m *ImageInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageInfo.Marshal(b, m, deterministic)
}
func (m *ImageInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageInfo.Merge(m, src)
}
func (m *ImageInfo) XXX_Size() int {
	return xxx_messageInfo_ImageInfo.Size(m)
}
func (m *ImageInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ImageInfo proto.InternalMessageInfo

func (m *ImageInfo) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ImageInfo) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

type ContainerStatus struct {
	Cid                  string   `protobuf:"bytes,1,opt,name=Cid,proto3" json:"Cid,omitempty"`
	NodeId               string   `protobuf:"bytes,2,opt,name=NodeId,proto3" json:"NodeId,omitempty"`
	Status               uint32   `protobuf:"varint,3,opt,name=Status,proto3" json:"Status,omitempty"`
	Image                string   `protobuf:"bytes,4,opt,name=Image,proto3" json:"Image,omitempty"`
	NetWorkRecord        uint64   `protobuf:"varint,5,opt,name=NetWorkRecord,proto3" json:"NetWorkRecord,omitempty"`
	NetWorkLimit         uint64   `protobuf:"varint,6,opt,name=NetWorkLimit,proto3" json:"NetWorkLimit,omitempty"`
	Addr                 string   `protobuf:"bytes,7,opt,name=Addr,proto3" json:"Addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContainerStatus) Reset()         { *m = ContainerStatus{} }
func (m *ContainerStatus) String() string { return proto.CompactTextString(m) }
func (*ContainerStatus) ProtoMessage()    {}
func (*ContainerStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{3}
}

func (m *ContainerStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContainerStatus.Unmarshal(m, b)
}
func (m *ContainerStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContainerStatus.Marshal(b, m, deterministic)
}
func (m *ContainerStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContainerStatus.Merge(m, src)
}
func (m *ContainerStatus) XXX_Size() int {
	return xxx_messageInfo_ContainerStatus.Size(m)
}
func (m *ContainerStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ContainerStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ContainerStatus proto.InternalMessageInfo

func (m *ContainerStatus) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *ContainerStatus) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *ContainerStatus) GetStatus() uint32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ContainerStatus) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

func (m *ContainerStatus) GetNetWorkRecord() uint64 {
	if m != nil {
		return m.NetWorkRecord
	}
	return 0
}

func (m *ContainerStatus) GetNetWorkLimit() uint64 {
	if m != nil {
		return m.NetWorkLimit
	}
	return 0
}

func (m *ContainerStatus) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type ImageList struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	ImageName            string   `protobuf:"bytes,2,opt,name=ImageName,proto3" json:"ImageName,omitempty"`
	Logo                 string   `protobuf:"bytes,3,opt,name=Logo,proto3" json:"Logo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImageList) Reset()         { *m = ImageList{} }
func (m *ImageList) String() string { return proto.CompactTextString(m) }
func (*ImageList) ProtoMessage()    {}
func (*ImageList) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{4}
}

func (m *ImageList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImageList.Unmarshal(m, b)
}
func (m *ImageList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImageList.Marshal(b, m, deterministic)
}
func (m *ImageList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageList.Merge(m, src)
}
func (m *ImageList) XXX_Size() int {
	return xxx_messageInfo_ImageList.Size(m)
}
func (m *ImageList) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageList.DiscardUnknown(m)
}

var xxx_messageInfo_ImageList proto.InternalMessageInfo

func (m *ImageList) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ImageList) GetImageName() string {
	if m != nil {
		return m.ImageName
	}
	return ""
}

func (m *ImageList) GetLogo() string {
	if m != nil {
		return m.Logo
	}
	return ""
}

type CallRsp struct {
	Code                 uint32           `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Status               bool             `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	Msg                  string           `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Data                 *ContainerStatus `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Cid                  string           `protobuf:"bytes,5,opt,name=cid,proto3" json:"cid,omitempty"`
	Token                string           `protobuf:"bytes,6,opt,name=token,proto3" json:"token,omitempty"`
	ImageList            []*ImageList     `protobuf:"bytes,7,rep,name=imageList,proto3" json:"imageList,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CallRsp) Reset()         { *m = CallRsp{} }
func (m *CallRsp) String() string { return proto.CompactTextString(m) }
func (*CallRsp) ProtoMessage()    {}
func (*CallRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{5}
}

func (m *CallRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CallRsp.Unmarshal(m, b)
}
func (m *CallRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CallRsp.Marshal(b, m, deterministic)
}
func (m *CallRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallRsp.Merge(m, src)
}
func (m *CallRsp) XXX_Size() int {
	return xxx_messageInfo_CallRsp.Size(m)
}
func (m *CallRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_CallRsp.DiscardUnknown(m)
}

var xxx_messageInfo_CallRsp proto.InternalMessageInfo

func (m *CallRsp) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *CallRsp) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *CallRsp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *CallRsp) GetData() *ContainerStatus {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *CallRsp) GetCid() string {
	if m != nil {
		return m.Cid
	}
	return ""
}

func (m *CallRsp) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *CallRsp) GetImageList() []*ImageList {
	if m != nil {
		return m.ImageList
	}
	return nil
}

type Call struct {
	Type                 CallType   `protobuf:"varint,1,opt,name=Type,proto3,enum=proto.gateway.service.CallType" json:"Type,omitempty"`
	Caller               CallerType `protobuf:"varint,2,opt,name=Caller,proto3,enum=proto.gateway.service.CallerType" json:"Caller,omitempty"`
	Opt                  *Options   `protobuf:"bytes,3,opt,name=Opt,proto3" json:"Opt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Call) Reset()         { *m = Call{} }
func (m *Call) String() string { return proto.CompactTextString(m) }
func (*Call) ProtoMessage()    {}
func (*Call) Descriptor() ([]byte, []int) {
	return fileDescriptor_087b5f483fc0b5df, []int{6}
}

func (m *Call) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Call.Unmarshal(m, b)
}
func (m *Call) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Call.Marshal(b, m, deterministic)
}
func (m *Call) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Call.Merge(m, src)
}
func (m *Call) XXX_Size() int {
	return xxx_messageInfo_Call.Size(m)
}
func (m *Call) XXX_DiscardUnknown() {
	xxx_messageInfo_Call.DiscardUnknown(m)
}

var xxx_messageInfo_Call proto.InternalMessageInfo

func (m *Call) GetType() CallType {
	if m != nil {
		return m.Type
	}
	return CallType_User
}

func (m *Call) GetCaller() CallerType {
	if m != nil {
		return m.Caller
	}
	return CallerType_Register
}

func (m *Call) GetOpt() *Options {
	if m != nil {
		return m.Opt
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.gateway.service.CallType", CallType_name, CallType_value)
	proto.RegisterEnum("proto.gateway.service.CallerType", CallerType_name, CallerType_value)
	proto.RegisterType((*UserOpt)(nil), "proto.gateway.service.UserOpt")
	proto.RegisterType((*Options)(nil), "proto.gateway.service.Options")
	proto.RegisterType((*ImageInfo)(nil), "proto.gateway.service.ImageInfo")
	proto.RegisterType((*ContainerStatus)(nil), "proto.gateway.service.ContainerStatus")
	proto.RegisterType((*ImageList)(nil), "proto.gateway.service.ImageList")
	proto.RegisterType((*CallRsp)(nil), "proto.gateway.service.CallRsp")
	proto.RegisterType((*Call)(nil), "proto.gateway.service.Call")
}

func init() { proto.RegisterFile("proto/gateway/gateway.proto", fileDescriptor_087b5f483fc0b5df) }

var fileDescriptor_087b5f483fc0b5df = []byte{
	// 651 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x89, 0x6b, 0x3b, 0x93, 0x9f, 0x5a, 0x03, 0x54, 0x56, 0x8b, 0x4a, 0x31, 0x08, 0xa1,
	0x4a, 0x04, 0xe4, 0x5e, 0x80, 0x03, 0x12, 0x4a, 0x45, 0x15, 0x29, 0xb4, 0x65, 0xdb, 0xd2, 0xf3,
	0x12, 0x6f, 0x8d, 0xd5, 0xc4, 0x6b, 0xd9, 0x5b, 0xaa, 0x3c, 0x00, 0x07, 0xde, 0x84, 0x03, 0xef,
	0xc1, 0xc3, 0xf0, 0x12, 0x68, 0xc7, 0xeb, 0x84, 0x56, 0x24, 0x9c, 0x76, 0x7e, 0xbe, 0x6f, 0x76,
	0xe6, 0xdb, 0x1d, 0xd8, 0xca, 0x0b, 0xa9, 0xe4, 0x8b, 0x84, 0x2b, 0x71, 0xcd, 0x67, 0xf5, 0xd9,
	0xa7, 0x28, 0xde, 0xa7, 0xa3, 0x5f, 0x07, 0x4b, 0x51, 0x7c, 0x4d, 0xc7, 0x22, 0xfc, 0x66, 0x81,
	0x7b, 0x56, 0x8a, 0xe2, 0x28, 0x57, 0xb8, 0x09, 0x9e, 0x36, 0x33, 0x3e, 0x15, 0x81, 0xb5, 0x63,
	0x3d, 0x6b, 0xb1, 0xb9, 0x8f, 0xf7, 0x60, 0xed, 0xf8, 0x8b, 0xcc, 0x44, 0xd0, 0xa0, 0x44, 0xe5,
	0x68, 0xc6, 0x31, 0x2f, 0xcb, 0x6b, 0x59, 0xc4, 0x41, 0xb3, 0x62, 0xd4, 0x3e, 0x6e, 0x80, 0xf3,
	0x49, 0x14, 0xe9, 0xc5, 0x2c, 0xb0, 0x29, 0x63, 0x3c, 0x5d, 0xe9, 0x54, 0x5e, 0x8a, 0x2c, 0x58,
	0xab, 0x2a, 0x91, 0x13, 0x7e, 0xb7, 0xc0, 0x3d, 0xca, 0x55, 0x2a, 0xb3, 0x12, 0x23, 0xb0, 0xf5,
	0xbd, 0xd4, 0x43, 0x3b, 0xda, 0xee, 0xff, 0xb3, 0xf3, 0xbe, 0xe9, 0x9a, 0x11, 0x16, 0x5f, 0x81,
	0x33, 0x28, 0x04, 0x57, 0x55, 0x83, 0xed, 0x68, 0x67, 0x09, 0x6b, 0x38, 0xe5, 0x89, 0x18, 0x66,
	0x17, 0x92, 0x19, 0x3c, 0xfa, 0xd0, 0x1c, 0xa4, 0x75, 0xfb, 0xda, 0x0c, 0x9f, 0x43, 0x6b, 0x0e,
	0xc3, 0x1e, 0x34, 0x86, 0xfb, 0xd4, 0x4a, 0x97, 0x35, 0x86, 0xfb, 0x35, 0xbc, 0xb1, 0x80, 0xff,
	0xb2, 0x60, 0x7d, 0x20, 0x33, 0xc5, 0xd3, 0x4c, 0x14, 0x27, 0x8a, 0xab, 0xab, 0xb2, 0x46, 0x59,
	0x73, 0x94, 0x96, 0xe3, 0x50, 0xc6, 0x62, 0x58, 0x53, 0x8d, 0xa7, 0xe3, 0x15, 0x87, 0x3a, 0xe8,
	0x32, 0xe3, 0x69, 0x99, 0xa8, 0x09, 0xa3, 0x5e, 0xe5, 0xe0, 0x13, 0xe8, 0x1e, 0x0a, 0x75, 0x2e,
	0x8b, 0x4b, 0x26, 0xc6, 0x5a, 0x75, 0x2d, 0xa2, 0xcd, 0x6e, 0x06, 0x31, 0x84, 0x8e, 0x09, 0x8c,
	0xd2, 0x69, 0xaa, 0x02, 0x87, 0x40, 0x37, 0x62, 0x88, 0x60, 0xbf, 0x8b, 0xe3, 0x22, 0x70, 0xa9,
	0x3c, 0xd9, 0xe1, 0x07, 0x33, 0xf8, 0x28, 0x2d, 0x15, 0x0d, 0x1e, 0xcf, 0x07, 0x8f, 0xf1, 0x81,
	0x49, 0x1e, 0xea, 0xef, 0x51, 0xcd, 0xb0, 0x08, 0xe8, 0x72, 0x23, 0x99, 0x48, 0x23, 0x23, 0xd9,
	0xe1, 0x6f, 0x0b, 0xdc, 0x01, 0x9f, 0x4c, 0x58, 0x99, 0xeb, 0xfc, 0x58, 0xc6, 0xc2, 0xd4, 0x23,
	0x5b, 0x8f, 0x5e, 0x56, 0xa3, 0xeb, 0x72, 0x1e, 0x33, 0x9e, 0x16, 0x6f, 0x5a, 0x26, 0xf5, 0x8b,
	0x4c, 0xcb, 0x04, 0xdf, 0x80, 0x1d, 0x73, 0xc5, 0x49, 0x8b, 0x76, 0xf4, 0x74, 0xc9, 0xdb, 0xde,
	0x7a, 0x04, 0x46, 0x1c, 0x5d, 0x6d, 0x9c, 0xc6, 0xe6, 0xb7, 0x69, 0x53, 0x4b, 0xab, 0xe8, 0x07,
	0x3a, 0x95, 0xb4, 0xe4, 0xe0, 0x5b, 0x68, 0xa5, 0xf5, 0xf0, 0x81, 0xbb, 0xd3, 0xfc, 0xdf, 0x27,
	0xd2, 0x38, 0xb6, 0xa0, 0x84, 0x3f, 0x2c, 0xb0, 0xf5, 0xb4, 0xb8, 0x07, 0xf6, 0xe9, 0x2c, 0xaf,
	0x46, 0xed, 0x45, 0x0f, 0x97, 0x35, 0xcb, 0x27, 0x13, 0x0d, 0x63, 0x04, 0xc6, 0xd7, 0xe0, 0xe8,
	0x88, 0x28, 0x48, 0x8b, 0x5e, 0xf4, 0x68, 0x05, 0x4d, 0x14, 0x44, 0x34, 0x04, 0x7c, 0x09, 0xcd,
	0xa3, 0x5c, 0x91, 0x5c, 0xcb, 0xb7, 0xc5, 0xec, 0x16, 0xd3, 0xd0, 0xdd, 0xc7, 0xe0, 0xd5, 0xd7,
	0xa3, 0x57, 0x2d, 0x9b, 0x7f, 0x07, 0xbb, 0xd0, 0x9a, 0x2b, 0xe8, 0x5b, 0xbb, 0x3f, 0x2d, 0x80,
	0xc5, 0x6d, 0xd8, 0x01, 0x8f, 0x89, 0x24, 0x2d, 0x15, 0x61, 0x3b, 0xe0, 0x9d, 0x88, 0x2c, 0x1e,
	0xc8, 0x58, 0xf8, 0x96, 0x66, 0xea, 0x1a, 0x23, 0x99, 0xa4, 0x99, 0xdf, 0xc0, 0x75, 0x68, 0x9f,
	0xe5, 0x31, 0x57, 0x82, 0x56, 0xdb, 0x6f, 0x62, 0x0f, 0xc0, 0xe4, 0xe5, 0x95, 0xf2, 0x6d, 0xbc,
	0x0b, 0xeb, 0xd5, 0xf2, 0x2d, 0xee, 0x5b, 0xc3, 0x0d, 0xc0, 0x03, 0xa1, 0x6e, 0xbd, 0xa1, 0xef,
	0x68, 0xf0, 0xbe, 0x98, 0x88, 0xbf, 0xc1, 0x2e, 0xfa, 0xd0, 0x39, 0x10, 0x6a, 0xfe, 0x0e, 0xbe,
	0x17, 0x7d, 0x04, 0xf7, 0x80, 0x2b, 0x71, 0xce, 0x67, 0xf8, 0x1e, 0xdc, 0x93, 0x6a, 0x6c, 0xdc,
	0x5a, 0x21, 0xe3, 0xe6, 0xf6, 0x8a, 0x24, 0x2b, 0xf3, 0xcf, 0x0e, 0xa5, 0xf7, 0xfe, 0x04, 0x00,
	0x00, 0xff, 0xff, 0xc9, 0x55, 0xb7, 0x53, 0x58, 0x05, 0x00, 0x00,
}