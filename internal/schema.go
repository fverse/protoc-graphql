package internal

import "strings"

type Schema struct {
	*strings.Builder
}

// Creates new Schema
func NewSchema() *Schema {
	c := new(Schema)
	c.Builder = new(strings.Builder)
	return c
}

// Puts a new line in the generated content
func (c *Schema) NewLine() {
	c.WriteString("\n")
}

func (c *Schema) Space() {
	c.WriteString(" ")
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
