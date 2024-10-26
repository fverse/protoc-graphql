package internal

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fverse/protoc-graphql/internal/descriptor"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Schema struct {
	*strings.Builder

	protoFile *descriptorpb.FileDescriptorProto

	objectTypes []*descriptor.ObjectType
	inputTypes  []*descriptor.InputType
	Mutations   []*descriptor.Mutation
	Queries     []*descriptor.Query
}

func generateFields(fields []*descriptorpb.FieldDescriptorProto) []*descriptor.Field {
	result := make([]*descriptor.Field, 0, len(fields))
	for _, field := range fields {
		f := &descriptor.Field{
			Name: field.Name,
		}
		f.GetType(field)
		f.IsOptional(field)
		f.IsRepeated(field)
		f.Print("field type: ", *f.Name, f.Type.String())
		result = append(result, f)
	}
	return result
}

// Constructs the Object types from message types and fills the
func (schema *Schema) ObjectTypes() {
	for _, message := range schema.protoFile.MessageType {
		if len(message.Field) > 0 {
			t := new(descriptor.ObjectType)
			t.Name = message.Name
			t.Fields = generateFields(message.Field)
			schema.objectTypes = append(schema.objectTypes, t)
		}
	}
}

// Creates new Schema
func CreateSchema(protoFile *descriptorpb.FileDescriptorProto) *Schema {
	schema := new(Schema)
	schema.Builder = new(strings.Builder)
	schema.protoFile = protoFile

	// Write the header content to the string builder
	schema.WriteHeader()

	// Construct Object types
	schema.ObjectTypes()

	// TODO: Generate Input types
	return schema
}

// Puts a new line in the generated content
func (c *Schema) NewLine() {
	c.Write("\n")
}

// Adds a space to the generated content
func (c *Schema) Space() {
	c.Write(" ")
}

// Puts a graphql comment in the generated content
func (c *Schema) Comment(s string) {
	c.Write("#")
	c.Space()
	c.Write(s)
}

// Write writes a string to the string builder
func (p *Schema) Write(s string) {
	if len(s) == 0 {
		return
	}
	p.WriteString(s)
}

// Creates a file name based on the given proto file name
func (p *Schema) FileName(filename *string) string {
	ext := filepath.Ext(*filename)
	return strings.TrimSuffix(*filename, ext) + ".graphql"
}

// Prints a message
func (p *Schema) Print(msg ...string) {
	s := strings.Join(msg, " ")
	log.Print(s)
}
