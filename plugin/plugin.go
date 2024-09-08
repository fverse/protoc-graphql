package plugin

import (
	"github.com/fverse/protoc-graphql/internal"
	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin struct {
	Response *pluginpb.CodeGeneratorResponse
	Args     *Args
	Logger   *Logger
}

func (p *Plugin) SetSupportOptionalField() {
	o := uint64(pluginpb.CodeGeneratorResponse_Feature_value["FEATURE_PROTO3_OPTIONAL"])
	p.Response.SupportedFeatures = &o
}

func New(request *pluginpb.CodeGeneratorRequest) *Plugin {
	logger := NewLogger()
	args := ParseArgs(request.GetParameter(), logger)

	response := internal.Generate()

	return &Plugin{
		Response: response,
		Args:     args,
		Logger:   logger,
	}
}
