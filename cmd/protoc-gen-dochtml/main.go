package main

import (
	"io"
	"os"
	"strings"

	"github.com/user/proto-doc/pkg/generator"
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

	filename := "index.html"
	cfg := generator.DefaultConfig

	if req.GetParameter() != "" {
		parts := strings.Split(req.GetParameter(), ",")
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
				switch k {
				case "out":
					filename = v
				case "title":
					cfg.Title = v
				case "primary_color":
					cfg.PrimaryColor = v
				case "font_family":
					cfg.FontFamily = v
				case "theme":
					cfg.Theme = v
				case "logo_url":
					cfg.LogoURL = v
				}
			}
		}
	}

	gen := generator.NewHTMLGenerator(cfg)
	htmlStr, err := gen.Render(docData)
	if err != nil {
		return
	}

	featureSupported := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	res := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &featureSupported,
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    proto.String(filename),
				Content: proto.String(htmlStr),
			},
		},
	}

	outputBytes, err := proto.Marshal(res)
	if err != nil {
		return
	}

	os.Stdout.Write(outputBytes)
}
