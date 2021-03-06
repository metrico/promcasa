// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: loki.proto

package logproto

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

type Timestamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative
	// second values with fractions must still have non-negative nanos values
	// that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanos int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

func (x *Timestamp) Reset() {
	*x = Timestamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loki_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timestamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timestamp) ProtoMessage() {}

func (x *Timestamp) ProtoReflect() protoreflect.Message {
	mi := &file_loki_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timestamp.ProtoReflect.Descriptor instead.
func (*Timestamp) Descriptor() ([]byte, []int) {
	return file_loki_proto_rawDescGZIP(), []int{0}
}

func (x *Timestamp) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *Timestamp) GetNanos() int32 {
	if x != nil {
		return x.Nanos
	}
	return 0
}

type PushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Streams []*StreamAdapter `protobuf:"bytes,1,rep,name=streams,proto3" json:"streams,omitempty"`
}

func (x *PushRequest) Reset() {
	*x = PushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loki_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRequest) ProtoMessage() {}

func (x *PushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_loki_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRequest.ProtoReflect.Descriptor instead.
func (*PushRequest) Descriptor() ([]byte, []int) {
	return file_loki_proto_rawDescGZIP(), []int{1}
}

func (x *PushRequest) GetStreams() []*StreamAdapter {
	if x != nil {
		return x.Streams
	}
	return nil
}

type PushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushResponse) Reset() {
	*x = PushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loki_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushResponse) ProtoMessage() {}

func (x *PushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_loki_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushResponse.ProtoReflect.Descriptor instead.
func (*PushResponse) Descriptor() ([]byte, []int) {
	return file_loki_proto_rawDescGZIP(), []int{2}
}

type StreamAdapter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Labels  string          `protobuf:"bytes,1,opt,name=labels,proto3" json:"labels,omitempty"`
	Entries []*EntryAdapter `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *StreamAdapter) Reset() {
	*x = StreamAdapter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loki_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamAdapter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamAdapter) ProtoMessage() {}

func (x *StreamAdapter) ProtoReflect() protoreflect.Message {
	mi := &file_loki_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamAdapter.ProtoReflect.Descriptor instead.
func (*StreamAdapter) Descriptor() ([]byte, []int) {
	return file_loki_proto_rawDescGZIP(), []int{3}
}

func (x *StreamAdapter) GetLabels() string {
	if x != nil {
		return x.Labels
	}
	return ""
}

func (x *StreamAdapter) GetEntries() []*EntryAdapter {
	if x != nil {
		return x.Entries
	}
	return nil
}

type EntryAdapter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp *Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Line      string     `protobuf:"bytes,2,opt,name=line,proto3" json:"line,omitempty"`
}

func (x *EntryAdapter) Reset() {
	*x = EntryAdapter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_loki_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntryAdapter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntryAdapter) ProtoMessage() {}

func (x *EntryAdapter) ProtoReflect() protoreflect.Message {
	mi := &file_loki_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntryAdapter.ProtoReflect.Descriptor instead.
func (*EntryAdapter) Descriptor() ([]byte, []int) {
	return file_loki_proto_rawDescGZIP(), []int{4}
}

func (x *EntryAdapter) GetTimestamp() *Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *EntryAdapter) GetLine() string {
	if x != nil {
		return x.Line
	}
	return ""
}

var File_loki_proto protoreflect.FileDescriptor

var file_loki_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6c, 0x6f, 0x6b, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6c, 0x6f,
	0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6e, 0x61,
	0x6e, 0x6f, 0x73, 0x22, 0x40, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x31, 0x0a, 0x07, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x6f, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x41, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x52, 0x07, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x73, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x59, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x41,
	0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x30,
	0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x6c, 0x6f, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x41, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x22, 0x55, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x41, 0x64, 0x61, 0x70, 0x74, 0x65, 0x72,
	0x12, 0x31, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6c, 0x6f, 0x67, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x6c, 0x6f, 0x67, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_loki_proto_rawDescOnce sync.Once
	file_loki_proto_rawDescData = file_loki_proto_rawDesc
)

func file_loki_proto_rawDescGZIP() []byte {
	file_loki_proto_rawDescOnce.Do(func() {
		file_loki_proto_rawDescData = protoimpl.X.CompressGZIP(file_loki_proto_rawDescData)
	})
	return file_loki_proto_rawDescData
}

var file_loki_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_loki_proto_goTypes = []interface{}{
	(*Timestamp)(nil),     // 0: logproto.Timestamp
	(*PushRequest)(nil),   // 1: logproto.PushRequest
	(*PushResponse)(nil),  // 2: logproto.PushResponse
	(*StreamAdapter)(nil), // 3: logproto.StreamAdapter
	(*EntryAdapter)(nil),  // 4: logproto.EntryAdapter
}
var file_loki_proto_depIdxs = []int32{
	3, // 0: logproto.PushRequest.streams:type_name -> logproto.StreamAdapter
	4, // 1: logproto.StreamAdapter.entries:type_name -> logproto.EntryAdapter
	0, // 2: logproto.EntryAdapter.timestamp:type_name -> logproto.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_loki_proto_init() }
func file_loki_proto_init() {
	if File_loki_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_loki_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timestamp); i {
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
		file_loki_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushRequest); i {
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
		file_loki_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushResponse); i {
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
		file_loki_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamAdapter); i {
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
		file_loki_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntryAdapter); i {
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
			RawDescriptor: file_loki_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_loki_proto_goTypes,
		DependencyIndexes: file_loki_proto_depIdxs,
		MessageInfos:      file_loki_proto_msgTypes,
	}.Build()
	File_loki_proto = out.File
	file_loki_proto_rawDesc = nil
	file_loki_proto_goTypes = nil
	file_loki_proto_depIdxs = nil
}
