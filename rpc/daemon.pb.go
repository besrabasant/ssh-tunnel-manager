// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v6.30.2
// source: rpc/daemon.proto

package rpc

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

type ResponseStatus int32

const (
	ResponseStatus_Success ResponseStatus = 0
	ResponseStatus_Error   ResponseStatus = 1
)

// Enum value maps for ResponseStatus.
var (
	ResponseStatus_name = map[int32]string{
		0: "Success",
		1: "Error",
	}
	ResponseStatus_value = map[string]int32{
		"Success": 0,
		"Error":   1,
	}
)

func (x ResponseStatus) Enum() *ResponseStatus {
	p := new(ResponseStatus)
	*p = x
	return p
}

func (x ResponseStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_rpc_daemon_proto_enumTypes[0].Descriptor()
}

func (ResponseStatus) Type() protoreflect.EnumType {
	return &file_rpc_daemon_proto_enumTypes[0]
}

func (x ResponseStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseStatus.Descriptor instead.
func (ResponseStatus) EnumDescriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{0}
}

type TunnelConfig struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Server        string                 `protobuf:"bytes,3,opt,name=server,proto3" json:"server,omitempty"`
	User          string                 `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	KeyFile       string                 `protobuf:"bytes,5,opt,name=key_file,json=keyFile,proto3" json:"key_file,omitempty"`
	RemoteHost    string                 `protobuf:"bytes,6,opt,name=remote_host,json=remoteHost,proto3" json:"remote_host,omitempty"`
	RemotePort    int32                  `protobuf:"varint,7,opt,name=remote_port,json=remotePort,proto3" json:"remote_port,omitempty"`
	LocalPort     int32                  `protobuf:"varint,8,opt,name=local_port,json=localPort,proto3" json:"local_port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TunnelConfig) Reset() {
	*x = TunnelConfig{}
	mi := &file_rpc_daemon_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TunnelConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TunnelConfig) ProtoMessage() {}

func (x *TunnelConfig) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TunnelConfig.ProtoReflect.Descriptor instead.
func (*TunnelConfig) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{0}
}

func (x *TunnelConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TunnelConfig) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *TunnelConfig) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *TunnelConfig) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *TunnelConfig) GetKeyFile() string {
	if x != nil {
		return x.KeyFile
	}
	return ""
}

func (x *TunnelConfig) GetRemoteHost() string {
	if x != nil {
		return x.RemoteHost
	}
	return ""
}

func (x *TunnelConfig) GetRemotePort() int32 {
	if x != nil {
		return x.RemotePort
	}
	return 0
}

func (x *TunnelConfig) GetLocalPort() int32 {
	if x != nil {
		return x.LocalPort
	}
	return 0
}

type ListConfigurationsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SearchPattern string                 `protobuf:"bytes,1,opt,name=search_pattern,json=searchPattern,proto3" json:"search_pattern,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListConfigurationsRequest) Reset() {
	*x = ListConfigurationsRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListConfigurationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListConfigurationsRequest) ProtoMessage() {}

func (x *ListConfigurationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListConfigurationsRequest.ProtoReflect.Descriptor instead.
func (*ListConfigurationsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{1}
}

func (x *ListConfigurationsRequest) GetSearchPattern() string {
	if x != nil {
		return x.SearchPattern
	}
	return ""
}

type ListConfigurationsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListConfigurationsResponse) Reset() {
	*x = ListConfigurationsResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListConfigurationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListConfigurationsResponse) ProtoMessage() {}

func (x *ListConfigurationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListConfigurationsResponse.ProtoReflect.Descriptor instead.
func (*ListConfigurationsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{2}
}

func (x *ListConfigurationsResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type AddOrUpdateConfigurationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data          *TunnelConfig          `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddOrUpdateConfigurationRequest) Reset() {
	*x = AddOrUpdateConfigurationRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddOrUpdateConfigurationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOrUpdateConfigurationRequest) ProtoMessage() {}

func (x *AddOrUpdateConfigurationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOrUpdateConfigurationRequest.ProtoReflect.Descriptor instead.
func (*AddOrUpdateConfigurationRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{3}
}

func (x *AddOrUpdateConfigurationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AddOrUpdateConfigurationRequest) GetData() *TunnelConfig {
	if x != nil {
		return x.Data
	}
	return nil
}

type AddOrUpdateConfigurationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddOrUpdateConfigurationResponse) Reset() {
	*x = AddOrUpdateConfigurationResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddOrUpdateConfigurationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddOrUpdateConfigurationResponse) ProtoMessage() {}

func (x *AddOrUpdateConfigurationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddOrUpdateConfigurationResponse.ProtoReflect.Descriptor instead.
func (*AddOrUpdateConfigurationResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{4}
}

func (x *AddOrUpdateConfigurationResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type FetchConfigurationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FetchConfigurationRequest) Reset() {
	*x = FetchConfigurationRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FetchConfigurationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchConfigurationRequest) ProtoMessage() {}

func (x *FetchConfigurationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchConfigurationRequest.ProtoReflect.Descriptor instead.
func (*FetchConfigurationRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{5}
}

func (x *FetchConfigurationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type FetchConfigurationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        ResponseStatus         `protobuf:"varint,1,opt,name=status,proto3,enum=daemon.ResponseStatus" json:"status,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data          *TunnelConfig          `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *FetchConfigurationResponse) Reset() {
	*x = FetchConfigurationResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FetchConfigurationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchConfigurationResponse) ProtoMessage() {}

func (x *FetchConfigurationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchConfigurationResponse.ProtoReflect.Descriptor instead.
func (*FetchConfigurationResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{6}
}

func (x *FetchConfigurationResponse) GetStatus() ResponseStatus {
	if x != nil {
		return x.Status
	}
	return ResponseStatus_Success
}

func (x *FetchConfigurationResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *FetchConfigurationResponse) GetData() *TunnelConfig {
	if x != nil {
		return x.Data
	}
	return nil
}

type DeleteConfigurationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteConfigurationRequest) Reset() {
	*x = DeleteConfigurationRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteConfigurationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteConfigurationRequest) ProtoMessage() {}

func (x *DeleteConfigurationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteConfigurationRequest.ProtoReflect.Descriptor instead.
func (*DeleteConfigurationRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteConfigurationRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DeleteConfigurationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteConfigurationResponse) Reset() {
	*x = DeleteConfigurationResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteConfigurationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteConfigurationResponse) ProtoMessage() {}

func (x *DeleteConfigurationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteConfigurationResponse.ProtoReflect.Descriptor instead.
func (*DeleteConfigurationResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{8}
}

func (x *DeleteConfigurationResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type StartTunnelRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigName    string                 `protobuf:"bytes,1,opt,name=config_name,json=configName,proto3" json:"config_name,omitempty"`
	LocalPort     int32                  `protobuf:"varint,2,opt,name=local_port,json=localPort,proto3" json:"local_port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartTunnelRequest) Reset() {
	*x = StartTunnelRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartTunnelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartTunnelRequest) ProtoMessage() {}

func (x *StartTunnelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartTunnelRequest.ProtoReflect.Descriptor instead.
func (*StartTunnelRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{9}
}

func (x *StartTunnelRequest) GetConfigName() string {
	if x != nil {
		return x.ConfigName
	}
	return ""
}

func (x *StartTunnelRequest) GetLocalPort() int32 {
	if x != nil {
		return x.LocalPort
	}
	return 0
}

type StartTunnelResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StartTunnelResponse) Reset() {
	*x = StartTunnelResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StartTunnelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartTunnelResponse) ProtoMessage() {}

func (x *StartTunnelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartTunnelResponse.ProtoReflect.Descriptor instead.
func (*StartTunnelResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{10}
}

func (x *StartTunnelResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type ListActiveTunnelsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListActiveTunnelsRequest) Reset() {
	*x = ListActiveTunnelsRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListActiveTunnelsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActiveTunnelsRequest) ProtoMessage() {}

func (x *ListActiveTunnelsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActiveTunnelsRequest.ProtoReflect.Descriptor instead.
func (*ListActiveTunnelsRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{11}
}

type ListActiveTunnelsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListActiveTunnelsResponse) Reset() {
	*x = ListActiveTunnelsResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListActiveTunnelsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListActiveTunnelsResponse) ProtoMessage() {}

func (x *ListActiveTunnelsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListActiveTunnelsResponse.ProtoReflect.Descriptor instead.
func (*ListActiveTunnelsResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{12}
}

func (x *ListActiveTunnelsResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

type KillTunnelRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ConfigName    string                 `protobuf:"bytes,1,opt,name=config_name,json=configName,proto3" json:"config_name,omitempty"`
	LocalPort     int32                  `protobuf:"varint,2,opt,name=local_port,json=localPort,proto3" json:"local_port,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KillTunnelRequest) Reset() {
	*x = KillTunnelRequest{}
	mi := &file_rpc_daemon_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KillTunnelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillTunnelRequest) ProtoMessage() {}

func (x *KillTunnelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillTunnelRequest.ProtoReflect.Descriptor instead.
func (*KillTunnelRequest) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{13}
}

func (x *KillTunnelRequest) GetConfigName() string {
	if x != nil {
		return x.ConfigName
	}
	return ""
}

func (x *KillTunnelRequest) GetLocalPort() int32 {
	if x != nil {
		return x.LocalPort
	}
	return 0
}

type KillTunnelResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        string                 `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *KillTunnelResponse) Reset() {
	*x = KillTunnelResponse{}
	mi := &file_rpc_daemon_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *KillTunnelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillTunnelResponse) ProtoMessage() {}

func (x *KillTunnelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_daemon_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillTunnelResponse.ProtoReflect.Descriptor instead.
func (*KillTunnelResponse) Descriptor() ([]byte, []int) {
	return file_rpc_daemon_proto_rawDescGZIP(), []int{14}
}

func (x *KillTunnelResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_rpc_daemon_proto protoreflect.FileDescriptor

var file_rpc_daemon_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x70, 0x63, 0x2f, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x06, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x22, 0xec, 0x01, 0x0a, 0x0c, 0x54,
	0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x19, 0x0a,
	0x08, 0x6b, 0x65, 0x79, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72,
	0x65, 0x6d, 0x6f, 0x74, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f,
	0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x42, 0x0a, 0x19, 0x4c, 0x69, 0x73,
	0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x5f, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x22, 0x34, 0x0a,
	0x1a, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x5f, 0x0a, 0x1f, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f,
	0x6e, 0x2e, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x0a, 0x20, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x2f, 0x0a, 0x19, 0x46, 0x65, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x90, 0x01, 0x0a, 0x1a, 0x46, 0x65, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f,
	0x6e, 0x2e, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x30, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x35, 0x0a, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x54, 0x0a,
	0x12, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x50,
	0x6f, 0x72, 0x74, 0x22, 0x2d, 0x0a, 0x13, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x75, 0x6e, 0x6e,
	0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x22, 0x1a, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x33,
	0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x54, 0x75, 0x6e, 0x6e,
	0x65, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x22, 0x53, 0x0a, 0x11, 0x4b, 0x69, 0x6c, 0x6c, 0x54, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x2c, 0x0a, 0x12, 0x4b, 0x69, 0x6c, 0x6c,
	0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2a, 0x28, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x01,
	0x32, 0xf1, 0x05, 0x0a, 0x0d, 0x44, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x5d, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f,
	0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x64, 0x61,
	0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x67, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x41,
	0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x6a, 0x0a, 0x13, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x27, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x72,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x64, 0x61, 0x65,
	0x6d, 0x6f, 0x6e, 0x2e, 0x41, 0x64, 0x64, 0x4f, 0x72, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5d, 0x0a, 0x12, 0x46, 0x65, 0x74, 0x63, 0x68, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x64,
	0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x60, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x2e, 0x64,
	0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x23, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1a, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x45, 0x0a, 0x0a, 0x4b, 0x69, 0x6c, 0x6c, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x12,
	0x19, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x6c, 0x6c, 0x54, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x64, 0x61, 0x65,
	0x6d, 0x6f, 0x6e, 0x2e, 0x4b, 0x69, 0x6c, 0x6c, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x12, 0x20, 0x2e,
	0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x65, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x64, 0x61, 0x65, 0x6d, 0x6f, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x65, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x65, 0x73, 0x72, 0x61, 0x62, 0x61, 0x73, 0x61, 0x6e, 0x74, 0x2f, 0x73,
	0x73, 0x68, 0x2d, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x72, 0x2f, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_daemon_proto_rawDescOnce sync.Once
	file_rpc_daemon_proto_rawDescData = file_rpc_daemon_proto_rawDesc
)

func file_rpc_daemon_proto_rawDescGZIP() []byte {
	file_rpc_daemon_proto_rawDescOnce.Do(func() {
		file_rpc_daemon_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_daemon_proto_rawDescData)
	})
	return file_rpc_daemon_proto_rawDescData
}

var file_rpc_daemon_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_rpc_daemon_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_rpc_daemon_proto_goTypes = []any{
	(ResponseStatus)(0),                      // 0: daemon.ResponseStatus
	(*TunnelConfig)(nil),                     // 1: daemon.TunnelConfig
	(*ListConfigurationsRequest)(nil),        // 2: daemon.ListConfigurationsRequest
	(*ListConfigurationsResponse)(nil),       // 3: daemon.ListConfigurationsResponse
	(*AddOrUpdateConfigurationRequest)(nil),  // 4: daemon.AddOrUpdateConfigurationRequest
	(*AddOrUpdateConfigurationResponse)(nil), // 5: daemon.AddOrUpdateConfigurationResponse
	(*FetchConfigurationRequest)(nil),        // 6: daemon.FetchConfigurationRequest
	(*FetchConfigurationResponse)(nil),       // 7: daemon.FetchConfigurationResponse
	(*DeleteConfigurationRequest)(nil),       // 8: daemon.DeleteConfigurationRequest
	(*DeleteConfigurationResponse)(nil),      // 9: daemon.DeleteConfigurationResponse
	(*StartTunnelRequest)(nil),               // 10: daemon.StartTunnelRequest
	(*StartTunnelResponse)(nil),              // 11: daemon.StartTunnelResponse
	(*ListActiveTunnelsRequest)(nil),         // 12: daemon.ListActiveTunnelsRequest
	(*ListActiveTunnelsResponse)(nil),        // 13: daemon.ListActiveTunnelsResponse
	(*KillTunnelRequest)(nil),                // 14: daemon.KillTunnelRequest
	(*KillTunnelResponse)(nil),               // 15: daemon.KillTunnelResponse
}
var file_rpc_daemon_proto_depIdxs = []int32{
	1,  // 0: daemon.AddOrUpdateConfigurationRequest.data:type_name -> daemon.TunnelConfig
	0,  // 1: daemon.FetchConfigurationResponse.status:type_name -> daemon.ResponseStatus
	1,  // 2: daemon.FetchConfigurationResponse.data:type_name -> daemon.TunnelConfig
	2,  // 3: daemon.DaemonService.ListConfigurations:input_type -> daemon.ListConfigurationsRequest
	4,  // 4: daemon.DaemonService.AddConfiguration:input_type -> daemon.AddOrUpdateConfigurationRequest
	4,  // 5: daemon.DaemonService.UpdateConfiguration:input_type -> daemon.AddOrUpdateConfigurationRequest
	6,  // 6: daemon.DaemonService.FetchConfiguration:input_type -> daemon.FetchConfigurationRequest
	8,  // 7: daemon.DaemonService.DeleteConfiguration:input_type -> daemon.DeleteConfigurationRequest
	10, // 8: daemon.DaemonService.StartTunnel:input_type -> daemon.StartTunnelRequest
	14, // 9: daemon.DaemonService.KillTunnel:input_type -> daemon.KillTunnelRequest
	12, // 10: daemon.DaemonService.ListActiveTunnels:input_type -> daemon.ListActiveTunnelsRequest
	3,  // 11: daemon.DaemonService.ListConfigurations:output_type -> daemon.ListConfigurationsResponse
	5,  // 12: daemon.DaemonService.AddConfiguration:output_type -> daemon.AddOrUpdateConfigurationResponse
	5,  // 13: daemon.DaemonService.UpdateConfiguration:output_type -> daemon.AddOrUpdateConfigurationResponse
	7,  // 14: daemon.DaemonService.FetchConfiguration:output_type -> daemon.FetchConfigurationResponse
	9,  // 15: daemon.DaemonService.DeleteConfiguration:output_type -> daemon.DeleteConfigurationResponse
	11, // 16: daemon.DaemonService.StartTunnel:output_type -> daemon.StartTunnelResponse
	15, // 17: daemon.DaemonService.KillTunnel:output_type -> daemon.KillTunnelResponse
	13, // 18: daemon.DaemonService.ListActiveTunnels:output_type -> daemon.ListActiveTunnelsResponse
	11, // [11:19] is the sub-list for method output_type
	3,  // [3:11] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_rpc_daemon_proto_init() }
func file_rpc_daemon_proto_init() {
	if File_rpc_daemon_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_daemon_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_daemon_proto_goTypes,
		DependencyIndexes: file_rpc_daemon_proto_depIdxs,
		EnumInfos:         file_rpc_daemon_proto_enumTypes,
		MessageInfos:      file_rpc_daemon_proto_msgTypes,
	}.Build()
	File_rpc_daemon_proto = out.File
	file_rpc_daemon_proto_rawDesc = nil
	file_rpc_daemon_proto_goTypes = nil
	file_rpc_daemon_proto_depIdxs = nil
}
