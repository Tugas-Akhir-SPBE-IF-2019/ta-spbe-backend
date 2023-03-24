// OpenAPI base definition of the service.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: spbe/service/agentmarketing/openapi.proto

package agentmarketing

import (
	_ "github.com/google/gnostic/openapiv3"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_spbe_service_agentmarketing_openapi_proto protoreflect.FileDescriptor

var file_spbe_service_agentmarketing_openapi_proto_rawDesc = []byte{
	0x0a, 0x29, 0x73, 0x70, 0x62, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x6f, 0x70,
	0x65, 0x6e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x62, 0x66, 0x69,
	0x2e, 0x62, 0x72, 0x61, 0x76, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x1b, 0x6f, 0x70,
	0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x33, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0xe8, 0x01, 0x0a, 0x23, 0x63, 0x6f,
	0x6d, 0x2e, 0x62, 0x66, 0x69, 0x2e, 0x62, 0x72, 0x61, 0x76, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x42, 0x07, 0x4f, 0x70, 0x65, 0x6e, 0x41, 0x50, 0x49, 0x48, 0x01, 0x50, 0x01, 0x5a, 0x59,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x66, 0x69, 0x2d, 0x66,
	0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x2f, 0x62, 0x66, 0x69, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x62, 0x66, 0x69, 0x2f, 0x62,
	0x72, 0x61, 0x76, 0x6f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x3b, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0xba, 0x47, 0x58, 0x12, 0x2c, 0x0a, 0x25,
	0x41, 0x67, 0x65, 0x6e, 0x74, 0x20, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x20,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x20, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x28, 0x3a, 0x26, 0x0a, 0x24,
	0x0a, 0x0b, 0x42, 0x65, 0x61, 0x72, 0x65, 0x72, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x15, 0x0a,
	0x13, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x2a, 0x06, 0x62, 0x65, 0x61, 0x72, 0x65, 0x72, 0x32,
	0x03, 0x4a, 0x57, 0x54, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_spbe_service_agentmarketing_openapi_proto_goTypes = []interface{}{}
var file_spbe_service_agentmarketing_openapi_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_spbe_service_agentmarketing_openapi_proto_init() }
func file_spbe_service_agentmarketing_openapi_proto_init() {
	if File_spbe_service_agentmarketing_openapi_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_spbe_service_agentmarketing_openapi_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_spbe_service_agentmarketing_openapi_proto_goTypes,
		DependencyIndexes: file_spbe_service_agentmarketing_openapi_proto_depIdxs,
	}.Build()
	File_spbe_service_agentmarketing_openapi_proto = out.File
	file_spbe_service_agentmarketing_openapi_proto_rawDesc = nil
	file_spbe_service_agentmarketing_openapi_proto_goTypes = nil
	file_spbe_service_agentmarketing_openapi_proto_depIdxs = nil
}
