package descriptor

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
	Fields   []*Field
	Name     *string
	Type     *GraphQLType
	Optional *bool
}

type GraphQLType string
