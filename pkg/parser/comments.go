package parser

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// PathKey converts a slice of int32 path indices to a string key.
func PathKey(path []int32) string {
	var strs []string
	for _, p := range path {
		strs = append(strs, fmt.Sprintf("%d", p))
	}
	return strings.Join(strs, ",")
}

func cleanComment(s string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "//")
		line = strings.TrimPrefix(line, "/*")
		line = strings.TrimSuffix(line, "*/")
		line = strings.TrimPrefix(line, "*")
		cleaned = append(cleaned, strings.TrimSpace(line))
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

// ExtractComments converts SourceCodeInfo_Location into Comments struct.
func ExtractComments(loc *descriptorpb.SourceCodeInfo_Location) Comments {
	if loc == nil {
		return Comments{}
	}

	leading := cleanComment(loc.GetLeadingComments())
	trailing := cleanComment(loc.GetTrailingComments())
	var detached []string
	for _, d := range loc.GetLeadingDetachedComments() {
		detached = append(detached, cleanComment(d))
	}

	fullDesc := leading
	if trailing != "" {
		if fullDesc != "" {
			fullDesc += "\n\n" + trailing
		} else {
			fullDesc = trailing
		}
	}

	return Comments{
		Leading:     leading,
		Trailing:    trailing,
		Detached:    detached,
		Description: fullDesc,
	}
}
