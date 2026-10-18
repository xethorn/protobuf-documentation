package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/user/proto-doc/pkg/parser"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil || len(inputBytes) == 0 {
		return
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(inputBytes, req); err != nil {
		return
	}

	fds := &descriptorpb.FileDescriptorSet{
		File: req.GetProtoFile(),
	}

	p := parser.NewParser(fds)
	docData := p.ParseAll()

	jsonBytes, err := json.MarshalIndent(docData, "", "  ")
	if err != nil {
		return
	}

	filename := "proto_docs.json"
	if req.GetParameter() != "" {
		filename = req.GetParameter()
	}

	featureSupported := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	res := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &featureSupported,
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    proto.String(filename),
				Content: proto.String(string(jsonBytes)),
			},
		},
	}

	outputBytes, err := proto.Marshal(res)
	if err != nil {
		return
	}

	os.Stdout.Write(outputBytes)
}
