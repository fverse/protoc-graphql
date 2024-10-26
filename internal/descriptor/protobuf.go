package descriptor

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

var wellKnownTypes = map[string]string{
	".google.protobuf.Timestamp": "String",
	".google.protobuf.Any":       "String",
}

func isWellKnownType(field *descriptorpb.FieldDescriptorProto) bool {
	return field.GetTypeName() == ".google.protobuf.Timestamp" || field.GetTypeName() == ".google.protobuf.Any"
}

// Extracts the type's name
func getTypeName(field *descriptorpb.FieldDescriptorProto) *string {
	t := strings.Split(*field.TypeName, ".")
	return &t[len(t)-1]
}

// Check if the field is optional
func isOptional(field *descriptorpb.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
}

// Checks if the field is required
func isRequired(field *descriptorpb.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptorpb.FieldDescriptorProto_LABEL_REQUIRED
}

// Check if the field is repeated
func isRepeated(field *descriptorpb.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}
