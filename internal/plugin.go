package internal

import (
	"log"
	"os"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin struct {
	Request    *pluginpb.CodeGeneratorRequest
	Response   *pluginpb.CodeGeneratorResponse
	ProtoFiles *descriptorpb.FieldDescriptorProto

	Args   *Args
	Logger *Logger

	schema *Schema
}

func (p *Plugin) SetSupportOptionalField() {
	o := uint64(pluginpb.CodeGeneratorResponse_Feature_value["FEATURE_PROTO3_OPTIONAL"])
	p.Response.SupportedFeatures = &o
}

// New creates a new Plugin
func New(request *pluginpb.CodeGeneratorRequest) *Plugin {
	logger := NewLogger()
	args := ParseArgs(request.GetParameter(), logger)
	schema := NewSchema()

	return &Plugin{
		Request:  request,
		Response: new(pluginpb.CodeGeneratorResponse),
		Args:     args,
		Logger:   logger,
		schema:   schema,
	}
}

func (p *Plugin) Version() string {
	v := "protoc-gen-graphql"
	return v
}

// Prints an error, and exits.
func (p *Plugin) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ": " + err.Error()
	log.Print("protoc-gen-graohql: error: ", s)
	os.Exit(1)
}

// Prints a message
func (p *Plugin) Info(msg ...string) {
	s := strings.Join(msg, " ")
	log.Print("prpotoc-gen-graphql: ", s)
}