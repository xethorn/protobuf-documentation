package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/user/proto-doc/pkg/generator"
	"github.com/user/proto-doc/pkg/parser"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func parseProtoWithProtoc(protoFiles []string, includePaths []string) (*parser.RootResult, error) {
	tempFile, err := os.CreateTemp("", "proto_descriptor_*.pb")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	args := []string{"--descriptor_set_out=" + tempFile.Name(), "--include_imports", "--include_source_info"}
	for _, inc := range includePaths {
		args = append(args, "-I"+inc)
	}
	args = append(args, protoFiles...)

	cmd := exec.Command("protoc", args...)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("protoc failed: %v, stderr: %s", err, errBuf.String())
	}

	descriptorData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(descriptorData, fds); err != nil {
		return nil, err
	}

	p := parser.NewParser(fds)
	return p.ParseAll(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: proto-doc <parse|html> [options]")
		os.Exit(1)
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "parse":
		parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
		var protoPaths arrayFlags
		parseCmd.Var(&protoPaths, "I", "Include paths")
		parseCmd.Var(&protoPaths, "proto_path", "Include paths")
		output := parseCmd.String("o", "-", "Output JSON file")
		parseCmd.Parse(os.Args[2:])

		files := parseCmd.Args()
		if len(files) == 0 {
			fmt.Println("Error: no proto files provided")
			os.Exit(1)
		}

		if len(protoPaths) == 0 {
			protoPaths = []string{"."}
		}

		result, err := parseProtoWithProtoc(files, protoPaths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
			os.Exit(1)
		}

		if *output == "-" {
			fmt.Println(string(jsonBytes))
		} else {
			if err := os.WriteFile(*output, jsonBytes, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Write file error: %v\n", err)
				os.Exit(1)
			}
		}

	case "html":
		htmlCmd := flag.NewFlagSet("html", flag.ExitOnError)
		var protoPaths arrayFlags
		htmlCmd.Var(&protoPaths, "I", "Include paths")
		htmlCmd.Var(&protoPaths, "proto_path", "Include paths")
		output := htmlCmd.String("o", "docs.html", "Output HTML file")
		title := htmlCmd.String("title", "API Documentation", "Document title")
		primaryColor := htmlCmd.String("primary-color", "#006aff", "Primary CSS accent color")
		fontFamily := htmlCmd.String("font-family", "system-ui, -apple-system, sans-serif", "Custom font family")
		theme := htmlCmd.String("theme", "light", "Default theme (light or dark)")
		logoURL := htmlCmd.String("logo-url", "", "Logo URL")

		htmlCmd.Parse(os.Args[2:])

		inputs := htmlCmd.Args()
		if len(inputs) == 0 {
			fmt.Println("Error: no input files provided")
			os.Exit(1)
		}

		var docData *parser.RootResult
		var err error

		if len(inputs) == 1 && filepath.Ext(inputs[0]) == ".json" {
			jsonBytes, err := os.ReadFile(inputs[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading JSON: %v\n", err)
				os.Exit(1)
			}
			docData = &parser.RootResult{}
			if err := json.Unmarshal(jsonBytes, docData); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
				os.Exit(1)
			}
		} else {
			if len(protoPaths) == 0 {
				protoPaths = []string{"."}
			}
			docData, err = parseProtoWithProtoc(inputs, protoPaths)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

		cfg := generator.Config{
			Title:        *title,
			PrimaryColor: *primaryColor,
			FontFamily:   *fontFamily,
			Theme:        *theme,
			LogoURL:      *logoURL,
		}

		gen := generator.NewHTMLGenerator(cfg)
		htmlStr, err := gen.Render(docData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "HTML render error: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(*output, []byte(htmlStr), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Write HTML error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully generated HTML documentation at: %s\n", *output)

	default:
		fmt.Printf("Unknown command: %s\n", subcommand)
		os.Exit(1)
	}
}
