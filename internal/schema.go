package internal

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fverse/protoc-graphql/internal/descriptor"
	"github.com/fverse/protoc-graphql/options"
	"github.com/fverse/protoc-graphql/pkg/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Schema struct {
	*strings.Builder

	// Plugin's parsed command line arguments
	args *Args

	protoFile *descriptorpb.FileDescriptorProto
	fileName  *string

	objectTypes []*descriptor.ObjectType
	enums       []*descriptor.Enumeration
	inputTypes  []*descriptor.InputType
	mutations   []*descriptor.Mutation
	queries     []*descriptor.Query
}

// Checks the keepCase option for the fields
func keepCase(fieldOptions *descriptorpb.FieldOptions) bool {
	if proto.HasExtension(fieldOptions, options.E_KeepCase) {
		ext := proto.GetExtension(fieldOptions, options.E_KeepCase)
		return ext.(bool)
	}
	return false
}

// Constructs the Object types from message types and fills the schema.objectTypes
func (schema *Schema) makeObjectTypes(messages []*descriptorpb.DescriptorProto) {
	for _, message := range messages {
		if len(message.Field) > 0 {
			objectType := new(descriptor.ObjectType)
			objectType.Name = message.Name

			// Generate type fields
			objectType.Fields = generateFields(message.Field)

			// Construct embedded object types
			for _, nested := range message.NestedType {
				schema.makeObjectTypes([]*descriptorpb.DescriptorProto{nested})
			}

			// Construct embedded enums
			for _, enumType := range message.EnumType {
				enum := new(descriptor.Enumeration)
				enum.Name = enumType.Name
				for _, value := range enumType.Value {
					enum.Values = append(enum.Values, enumValues(value))
				}
				schema.enums = append(schema.enums, enum)
			}
			schema.objectTypes = append(schema.objectTypes, objectType)
		}
	}
}

// Return the string value of the provided enum value
func enumValues(value *descriptorpb.EnumValueDescriptorProto) *string {
	return value.Name
}

// Constructs the fields of an object type
func generateFields(fields []*descriptorpb.FieldDescriptorProto) []*descriptor.Field {
	result := make([]*descriptor.Field, 0, len(fields))

	for _, field := range fields {
		f := &descriptor.Field{
			Name: field.Name,
		}
		// Obtain the type of field
		f.GetType(field)

		// Sets wether the field is optional or not
		f.IsOptional(field)

		// Sets wether the field is required or not
		f.IsRepeated(field)

		if !keepCase(field.GetOptions()) {
			f.Name = utils.String(utils.CamelCase(*field.Name))
		}
		result = append(result, f)
	}
	return result
}

// Constructs mutations from service methods and fills the schema.mutations
func (schema *Schema) Mutations() {
}

func getMethodOptions(method *descriptorpb.MethodDescriptorProto) *descriptor.MethodOptions {
	opts := method.GetOptions()
	if proto.HasExtension(opts, options.E_Method) {
		ext := proto.GetExtension(opts, options.E_Method)
		return ext.(*descriptor.MethodOptions)
	}
	return &descriptor.MethodOptions{}
}

func getGqlOutputType(outputType string, mo *string) *string {
	if outputType != "" {
		outputType = utils.UppercaseFirst(outputType)
		return &outputType
	}
	outputType = strings.TrimPrefix(*mo, "."+PACKAGE+".")
	return &outputType
}

func isBoolean(t *string) bool {
	return strings.Contains(*t, "Bool")
}

func isEmpty(t *string) bool {
	// query.Input.Type == empty || query.Input.Type == "Empty" || query.Input.Type == "empty"
	return *t == "Empty"
}

func isArray(t *descriptor.GqlInput, length int) bool {
	f := t.Type[:1]
	l := t.Type[length-1:]
	return f == "[" && l == "]"
}

func parseType(t *descriptor.GqlInput) *descriptor.GqlInput {
	length := len(t.Type)
	if length == 0 {
		return nil
	}

	if isArray(t, length) {
		t.Array = true
		t.Type = utils.UppercaseFirst(t.Type[1 : length-1])
	} else {
		t.Type = utils.UppercaseFirst(t.Type)
	}

	if isPrimitive(&t.Type) {
		t.Primitive = true
		if isBoolean(&t.Type) {
			t.Type = "Boolean"
		}
	} else if isEmpty(&t.Type) {
		t.Empty = true
	}
	return t
}

func isPrimitive(t *string) bool {
	switch *t {
	case "String", "Boolean", "Bool", "Int", "Float":
		return true
	default:
		return false
	}
}

const PACKAGE = ".hello."

func getGqlInputType(input *descriptor.GqlInput, mi *string) *descriptor.GqlInput {
	if input.Type != "" {
		input = parseType(input)
		if !input.Primitive && !input.Empty {
			input.Type = "I" + input.Type
		}
		// if input.Array {
		// 	t.Type = "[" + t.Type + "]"
		// }
		return input
	}
	input.Type = "I" + strings.TrimPrefix(*mi, "."+PACKAGE+".")
	return input
}

// Check the compiler target and the method's target
func checkCompilerTarget(compilerTarget *string, options *descriptor.MethodOptions) bool {
	return *compilerTarget == utils.CastUit32ToString(options.Target) || utils.CompareStringInt(*compilerTarget, 3)
}

// Check if the compiler target is not the same as the method's target
func skipMethod(compilerTarget *string, options *descriptor.MethodOptions) bool {
	return options.Skip || !checkCompilerTarget(compilerTarget, options) && options.Target != 3
}

// Constructs the Object types from message types and fills the schema.objectTypes
func (schema *Schema) AddQueriesAndMutations() {
	for _, service := range schema.protoFile.Service {
		for _, method := range service.Method {

			methodOptions := getMethodOptions(method)
			if skipMethod(&schema.args.Target, methodOptions) {
				continue
			}

			if methodOptions.Kind == "mutation" || methodOptions.Kind == "Mutation" {
				mutation := new(descriptor.Mutation)
				mutation.Options = methodOptions
				mutation.Name = method.Name
				mutation.Input = getGqlInputType(methodOptions.GqlInput, method.InputType)
				// mutation.Payload = getGqlOutputType(options.GqlOutput, method.OutputType)
				schema.mutations = append(schema.mutations, mutation)
			} else {
				query := new(descriptor.Query)
				query.Options = methodOptions
				query.Name = method.Name
				query.Input = getGqlInputType(query.Options.GqlInput, method.InputType)
				// query.Payload =
				schema.queries = append(schema.queries, query)
			}
		}
	}
}

// Construct enums
func (schema *Schema) Enums() {
	for _, enumType := range schema.protoFile.EnumType {
		enum := new(descriptor.Enumeration)
		enum.Name = enumType.Name
		for _, value := range enumType.Value {
			enum.Values = append(enum.Values, enumValues(value))
		}
		schema.enums = append(schema.enums, enum)
	}
}

// Creates new Schema
func CreateSchema(protoFile *descriptorpb.FileDescriptorProto) *Schema {
	schema := new(Schema)
	schema.Builder = new(strings.Builder)
	schema.protoFile = protoFile

	// Crete the filename for the output file
	schema.FileName(protoFile.Name)

	// Construct Object types
	schema.makeObjectTypes(protoFile.MessageType)

	schema.Enums()

	schema.AddQueriesAndMutations()

	// TODO: Generate Input types
	return schema
}

// Puts a new line in the generated content
func (schema *Schema) NewLine(length ...int) {
	if len(length) == 0 {
		schema.Write("\n")
		return
	}
	for i := 0; i < length[0]; i++ {
		schema.Write("\n")
	}
}

// Adds a space to the generated content
func (schema *Schema) Space(length ...int) {
	if len(length) == 0 {
		schema.Write(" ")
		return
	}
	for i := 0; i < length[0]; i++ {
		schema.Write(" ")
	}
}

// Puts a graphql comment in the generated content
func (schema *Schema) Comment(s string) {
	schema.Write("#")
	schema.Space()
	schema.Write(s)
}

// Write writes a string to the string builder
func (schema *Schema) Write(s string) {
	if len(s) == 0 {
		return
	}
	schema.WriteString(s)
}

// Creates a file name based on the given proto file name
func (schema *Schema) FileName(filename *string) {
	ext := filepath.Ext(*filename)
	schema.fileName = utils.String(strings.TrimSuffix(*filename, ext) + ".graphql")
}

// Prints a message
func (schema *Schema) Print(msg ...string) {
	s := strings.Join(msg, " ")
	log.Print(s)
}

// Write the header content
func (schema *Schema) WriteHeader() {
	schema.NewLine()
	schema.Comment("Auto-generated by protoc-gen-graphql. DO NOT EDIT\n")
	schema.Comment(NAME + " " + VERSION)
	schema.NewLine()
	schema.NewLine()
}
