package internal

import (
	"fmt"

	"github.com/fverse/protoc-graphql/internal/descriptor"
	"github.com/fverse/protoc-graphql/internal/syntax"
	"github.com/fverse/protoc-graphql/pkg/utils"
)

const dotPackage string = ".hello."

var empty = dotPackage + "Empty"

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

// Generate queries
func (schema *Schema) generateQueries() {
	schema.Write("type Query {\n")
	schema.NewLine()

	for _, query := range schema.queries {
		if query.Input.Empty {
			schema.Write(fmt.Sprintf("  %s: %s\n", utils.LowercaseFirst(*query.Name), query.Payload.Type))
		} else {
			if query.Options.GqlInput.Optional {
				schema.Write(fmt.Sprintf("  %s(%s: %s): %s!\n", utils.LowercaseFirst(*query.Name),
					query.Input.Param, query.Input.Type, query.Payload.Type))
			} else {
				schema.Write(fmt.Sprintf("  %s(%s: %s!): %s!\n", utils.LowercaseFirst(*query.Name),
					query.Input.Param, query.Input.Type, query.Payload.Type))
			}
		}
		schema.NewLine()
		// q(input: InputType): ObjectType
	}
	schema.NewLine()
	schema.Write("}")
	schema.NewLine()
}

func (schema *Schema) generate() {
	// Write the header content to the string builder
	schema.WriteHeader()

	// Generate types
	schema.generateTypes()

	// Generate enums
	schema.generateEnums()

	// Generate queries
	schema.generateQueries()

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

func (schema *Schema) WriteMethod(kind, name *string) {
	schema.Write("")
	schema.Space()
	schema.Write(*name)
	schema.Space()
	schema.Write(string(syntax.LBrace))
	schema.NewLine()
}
