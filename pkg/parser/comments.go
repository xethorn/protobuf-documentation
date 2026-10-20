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

// CleanComment cleans multi-line comments, removing JSDoc/Javadoc delimiters (/**, */, *) and inline // prefixes.
func CleanComment(s string) string {
	if s == "" {
		return ""
	}

	lines := strings.Split(s, "\n")
	var cleaned []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Strip block comment openers/closers
		if strings.HasPrefix(line, "/**") {
			line = strings.TrimPrefix(line, "/**")
			line = strings.TrimSpace(line)
		} else if strings.HasPrefix(line, "/*") {
			line = strings.TrimPrefix(line, "/*")
			line = strings.TrimSpace(line)
		}

		if strings.HasSuffix(line, "*/") {
			line = strings.TrimSuffix(line, "*/")
			line = strings.TrimSpace(line)
		}

		// Strip leading JSDoc asterisks
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
		}

		// Strip inline // prefixes
		if strings.HasPrefix(line, "//") {
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)
		}

		if line != "" {
			cleaned = append(cleaned, line)
		}
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

// ExtractComments converts SourceCodeInfo_Location into Comments struct.
func ExtractComments(loc *descriptorpb.SourceCodeInfo_Location) Comments {
	if loc == nil {
		return Comments{}
	}

	leading := CleanComment(loc.GetLeadingComments())
	trailing := CleanComment(loc.GetTrailingComments())
	var detached []string
	for _, d := range loc.GetLeadingDetachedComments() {
		detached = append(detached, CleanComment(d))
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
