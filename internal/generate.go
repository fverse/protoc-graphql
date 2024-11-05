package internal

import (
	"github.com/fverse/protoc-graphql/internal/descriptor"
	"github.com/fverse/protoc-graphql/internal/syntax"
)

func (schema *Schema) generateType(object *descriptor.ObjectType) {
	schema.WriteTypeName(object.Name)

	for _, field := range object.Fields {
		schema.Space(3)
		schema.Write(*field.Name + string(syntax.Colon))
		schema.Space()
		schema.Write(field.Type.String())
		schema.NewLine()
	}
	schema.Write(string(syntax.RBrace))
	schema.NewLine(2)
}

func (schema *Schema) generateTypes() {
	for _, object := range schema.objectTypes {
		schema.generateType(object)
	}
}

func (schema *Schema) generate() {
	// Write the header content to the string builder
	schema.WriteHeader()

	// Generate types
	schema.generateTypes()

	// TODO: Generate enums

	// TODO: Generate queries

	// TODO: Generate mutations
}

// Writes the type's name
func (schema *Schema) WriteTypeName(name *string) {
	schema.Write(string(syntax.Type))
	schema.Space()
	schema.Write(*name)
	schema.Space()
	schema.Write(string(syntax.LBrace))
	schema.NewLine()
}
