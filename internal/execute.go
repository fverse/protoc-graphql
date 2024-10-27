package internal

import (
	"github.com/fverse/protoc-graphql/pkg/utils"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Checks if the proto files is explicitly passed in the command line
func (plugin *Plugin) isFileExplicit(protoFile *descriptorpb.FileDescriptorProto) bool {
	for _, file := range plugin.Request.FileToGenerate {
		if file == *protoFile.Name {
			return true
		}
	}
	return false
}

// Generates the protoc response
func (plugin *Plugin) Execute() {
	for _, protoFile := range plugin.Request.ProtoFile {
		if plugin.isFileExplicit(protoFile) {
			schema := CreateSchema(protoFile)
			schema.args = plugin.args

			schema.generate()

			plugin.Response.File = append(plugin.Response.File, &pluginpb.CodeGeneratorResponse_File{
				Name:    schema.fileName,
				Content: utils.String(schema.String()),
			})
			plugin.schema = append(plugin.schema, schema)
		}
	}
}
