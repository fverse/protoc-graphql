package descriptor

import (
	"log"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

type Mutation struct {
	Name    *string
	Input   *InputType
	Payload *ObjectType
}

type Query struct {
	Name    *string
	Input   *InputType
	Payload *ObjectType
}

type ObjectType struct {
	Fields []*Field
	Name   *string
}

type InputType struct {
	Fields []*Field
	Name   *string
}

type Field struct {
	Name     *string
	Type     *Scalar
	TypeName *string
	Optional *bool
	Repeated *bool
}

type Scalar string

const (
	Int     Scalar = "Int"
	Float   Scalar = "Float"
	Boolean Scalar = "Bool"
	String  Scalar = "String"
	Object  Scalar = "type"
	Input   Scalar = "input"
	Enum    Scalar = "enum"
	Unknown Scalar = "Unknown"
)

func (s *Scalar) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Copies the given scalar value and returns a pointer to it.
func scalar(v Scalar) *Scalar {
	return &v
}

// Stores value of v and returns a pointer to it
func Bool(v bool) *bool {
	return &v
}

// Check if the field is optional
func (f *Field) IsOptional(field *descriptorpb.FieldDescriptorProto) {
	f.Optional = Bool(isOptional(field))
}

// Checks if the field is required
func (f *Field) IsRequired(field *descriptorpb.FieldDescriptorProto) {
	f.Optional = Bool(!isRequired(field))
}

// Check if the field is repeated
func (f *Field) IsRepeated(field *descriptorpb.FieldDescriptorProto) {
	f.Repeated = Bool(isRepeated(field))
}

// Prints a message
func (p *Field) Print(msg ...string) {
	s := strings.Join(msg, " ")
	log.Print(s)
}

func (f *Field) GetType(field *descriptorpb.FieldDescriptorProto) {
	switch *field.Type {
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		f.Type = scalar(Int)
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		f.Type = scalar(Float)
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		f.Type = scalar(Boolean)
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		f.Type = scalar(String)
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if isWellKnownType(field) {
			// TODO: This needs to mapped to a custom Gql scalar type instead of string
			f.Type = scalar(String)
		} else {
			f.Type = scalar(Object)
			f.TypeName = getTypeName(field)
		}
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		f.Type = scalar(Enum)
		f.TypeName = getTypeName(field)
	default:
		f.Type = scalar(Unknown) // TODO: handle this
	}
}
