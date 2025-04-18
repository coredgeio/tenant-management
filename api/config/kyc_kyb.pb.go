// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: kyc_kyb.proto

package config

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StatusInfo int32

const (
	// Kyc pending
	StatusInfo_KycPending StatusInfo = 0
	// Kyc in-process
	StatusInfo_KycInProcess StatusInfo = 1
	// Kyc failed
	StatusInfo_KycFailed StatusInfo = 2
	// Kyc Partial
	StatusInfo_KycPartial StatusInfo = 3
	// Kyc done
	StatusInfo_KycDone StatusInfo = 4
	// rekyc needed
	StatusInfo_ReKycNeeded StatusInfo = 5
)

// Enum value maps for StatusInfo.
var (
	StatusInfo_name = map[int32]string{
		0: "KycPending",
		1: "KycInProcess",
		2: "KycFailed",
		3: "KycPartial",
		4: "KycDone",
		5: "ReKycNeeded",
	}
	StatusInfo_value = map[string]int32{
		"KycPending":   0,
		"KycInProcess": 1,
		"KycFailed":    2,
		"KycPartial":   3,
		"KycDone":      4,
		"ReKycNeeded":  5,
	}
)

func (x StatusInfo) Enum() *StatusInfo {
	p := new(StatusInfo)
	*p = x
	return p
}

func (x StatusInfo) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusInfo) Descriptor() protoreflect.EnumDescriptor {
	return file_kyc_kyb_proto_enumTypes[0].Descriptor()
}

func (StatusInfo) Type() protoreflect.EnumType {
	return &file_kyc_kyb_proto_enumTypes[0]
}

func (x StatusInfo) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusInfo.Descriptor instead.
func (StatusInfo) EnumDescriptor() ([]byte, []int) {
	return file_kyc_kyb_proto_rawDescGZIP(), []int{0}
}

type KycStatusGetReq struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// org name for kyc req
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KycStatusGetReq) Reset() {
	*x = KycStatusGetReq{}
	mi := &file_kyc_kyb_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KycStatusGetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KycStatusGetReq) ProtoMessage() {}

func (x *KycStatusGetReq) ProtoReflect() protoreflect.Message {
	mi := &file_kyc_kyb_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KycStatusGetReq.ProtoReflect.Descriptor instead.
func (*KycStatusGetReq) Descriptor() ([]byte, []int) {
	return file_kyc_kyb_proto_rawDescGZIP(), []int{0}
}

func (x *KycStatusGetReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type KybStatusGetReq struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// org name for kyb req
	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KybStatusGetReq) Reset() {
	*x = KybStatusGetReq{}
	mi := &file_kyc_kyb_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KybStatusGetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KybStatusGetReq) ProtoMessage() {}

func (x *KybStatusGetReq) ProtoReflect() protoreflect.Message {
	mi := &file_kyc_kyb_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KybStatusGetReq.ProtoReflect.Descriptor instead.
func (*KybStatusGetReq) Descriptor() ([]byte, []int) {
	return file_kyc_kyb_proto_rawDescGZIP(), []int{1}
}

func (x *KybStatusGetReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type KycStatusResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	KycStatus     StatusInfo             `protobuf:"varint,1,opt,name=kycStatus,proto3,enum=config.StatusInfo" json:"kycStatus,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KycStatusResp) Reset() {
	*x = KycStatusResp{}
	mi := &file_kyc_kyb_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KycStatusResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KycStatusResp) ProtoMessage() {}

func (x *KycStatusResp) ProtoReflect() protoreflect.Message {
	mi := &file_kyc_kyb_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KycStatusResp.ProtoReflect.Descriptor instead.
func (*KycStatusResp) Descriptor() ([]byte, []int) {
	return file_kyc_kyb_proto_rawDescGZIP(), []int{2}
}

func (x *KycStatusResp) GetKycStatus() StatusInfo {
	if x != nil {
		return x.KycStatus
	}
	return StatusInfo_KycPending
}

type KybStatusResp struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	KybStatus     StatusInfo             `protobuf:"varint,1,opt,name=kybStatus,proto3,enum=config.StatusInfo" json:"kybStatus,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KybStatusResp) Reset() {
	*x = KybStatusResp{}
	mi := &file_kyc_kyb_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KybStatusResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KybStatusResp) ProtoMessage() {}

func (x *KybStatusResp) ProtoReflect() protoreflect.Message {
	mi := &file_kyc_kyb_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KybStatusResp.ProtoReflect.Descriptor instead.
func (*KybStatusResp) Descriptor() ([]byte, []int) {
	return file_kyc_kyb_proto_rawDescGZIP(), []int{3}
}

func (x *KybStatusResp) GetKybStatus() StatusInfo {
	if x != nil {
		return x.KybStatus
	}
	return StatusInfo_KycPending
}

var File_kyc_kyb_proto protoreflect.FileDescriptor

const file_kyc_kyb_proto_rawDesc = "" +
	"\n" +
	"\rkyc_kyb.proto\x12\x06config\x1a\x1cgoogle/api/annotations.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"%\n" +
	"\x0fKycStatusGetReq\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"%\n" +
	"\x0fKybStatusGetReq\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\"A\n" +
	"\rKycStatusResp\x120\n" +
	"\tkycStatus\x18\x01 \x01(\x0e2\x12.config.StatusInfoR\tkycStatus\"A\n" +
	"\rKybStatusResp\x120\n" +
	"\tkybStatus\x18\x01 \x01(\x0e2\x12.config.StatusInfoR\tkybStatus*k\n" +
	"\n" +
	"StatusInfo\x12\x0e\n" +
	"\n" +
	"KycPending\x10\x00\x12\x10\n" +
	"\fKycInProcess\x10\x01\x12\r\n" +
	"\tKycFailed\x10\x02\x12\x0e\n" +
	"\n" +
	"KycPartial\x10\x03\x12\v\n" +
	"\aKycDone\x10\x04\x12\x0f\n" +
	"\vReKycNeeded\x10\x052\xf0\x01\n" +
	"\x10TenantManagement\x12m\n" +
	"\fGetKycStatus\x12\x17.config.KycStatusGetReq\x1a\x15.config.KycStatusResp\"-\x82\xd3\xe4\x93\x02'\x12%/api/tenant-mgmt/v1/tenant/{name}/kyc\x12m\n" +
	"\fGetKybStatus\x12\x17.config.KybStatusGetReq\x1a\x15.config.KybStatusResp\"-\x82\xd3\xe4\x93\x02'\x12%/api/tenant-mgmt/v1/tenant/{name}/kybBu\x92A?\x12\x1c\n" +
	"\x15TenantManagement API 2\x031.0r\x1f\n" +
	"\x1dTenantManagement API, KYC_KYBZ1github.com/coredgeio/tenant-management/api/configb\x06proto3"

var (
	file_kyc_kyb_proto_rawDescOnce sync.Once
	file_kyc_kyb_proto_rawDescData []byte
)

func file_kyc_kyb_proto_rawDescGZIP() []byte {
	file_kyc_kyb_proto_rawDescOnce.Do(func() {
		file_kyc_kyb_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_kyc_kyb_proto_rawDesc), len(file_kyc_kyb_proto_rawDesc)))
	})
	return file_kyc_kyb_proto_rawDescData
}

var file_kyc_kyb_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_kyc_kyb_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_kyc_kyb_proto_goTypes = []any{
	(StatusInfo)(0),         // 0: config.StatusInfo
	(*KycStatusGetReq)(nil), // 1: config.KycStatusGetReq
	(*KybStatusGetReq)(nil), // 2: config.KybStatusGetReq
	(*KycStatusResp)(nil),   // 3: config.KycStatusResp
	(*KybStatusResp)(nil),   // 4: config.KybStatusResp
}
var file_kyc_kyb_proto_depIdxs = []int32{
	0, // 0: config.KycStatusResp.kycStatus:type_name -> config.StatusInfo
	0, // 1: config.KybStatusResp.kybStatus:type_name -> config.StatusInfo
	1, // 2: config.TenantManagement.GetKycStatus:input_type -> config.KycStatusGetReq
	2, // 3: config.TenantManagement.GetKybStatus:input_type -> config.KybStatusGetReq
	3, // 4: config.TenantManagement.GetKycStatus:output_type -> config.KycStatusResp
	4, // 5: config.TenantManagement.GetKybStatus:output_type -> config.KybStatusResp
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kyc_kyb_proto_init() }
func file_kyc_kyb_proto_init() {
	if File_kyc_kyb_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_kyc_kyb_proto_rawDesc), len(file_kyc_kyb_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kyc_kyb_proto_goTypes,
		DependencyIndexes: file_kyc_kyb_proto_depIdxs,
		EnumInfos:         file_kyc_kyb_proto_enumTypes,
		MessageInfos:      file_kyc_kyb_proto_msgTypes,
	}.Build()
	File_kyc_kyb_proto = out.File
	file_kyc_kyb_proto_goTypes = nil
	file_kyc_kyb_proto_depIdxs = nil
}
