package internal

import (
	"path/filepath"
	"strings"

	"github.com/fverse/protoc-graphql/internal/descriptor"
)

type Schema struct {
	*strings.Builder

	ObjectTypes []*descriptor.ObjectType
	InputTypes  []*descriptor.InputType

	Mutations []*descriptor.Mutation
	Queries   []*descriptor.Query
}

// Creates new Schema
func NewSchema() *Schema {
	c := new(Schema)
	c.Builder = new(strings.Builder)
	c.WriteHeader()
	return c
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
