package internal

import (
	"github.com/fverse/protoc-graphql/internal/descriptor"
	"github.com/fverse/protoc-graphql/internal/syntax"
)

func (schema *Schema) generateType(object *descriptor.ObjectType) {
	schema.WriteTypeName(syntax.ObjectType, object.Name)

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

// Generate enums
func (schema *Schema) generateEnums() {
	for _, enum := range schema.enums {
		schema.WriteTypeName(syntax.Enum, enum.Name)

		for _, value := range enum.Values {
			schema.Space(3)
			schema.Write(*value)
			schema.NewLine()
		}
		schema.Write(string(syntax.RBrace))
		schema.NewLine(2)
	}
}

func (schema *Schema) generate() {
	// Write the header content to the string builder
	schema.WriteHeader()

	// Generate types
	schema.generateTypes()

	// TODO: Generate enums
	schema.generateEnums()

	// TODO: Generate queries

	// TODO: Generate mutations
}

// Writes the type's name
func (schema *Schema) WriteTypeName(keyWord syntax.Keyword, name *string) {
	schema.Write(string(keyWord))
	schema.Space()
	schema.Write(*name)
	schema.Space()
	schema.Write(string(syntax.LBrace))
	schema.NewLine()
}
