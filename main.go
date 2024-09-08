package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fverse/protoc-graphql/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading proto: %v\n", err)
		os.Exit(1)
	}

	var request pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(data, &request); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing proto: %v\n", err)
		os.Exit(1)
	}

	p := plugin.New(&request)

	p.SetSupportOptionalField()

	p.Logger.Log("p %v", p.Args)

	output, err := proto.Marshal(p.Response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing output: %v\n", err)
		os.Exit(1)
	}

	_, err = os.Stdout.Write(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}
}
