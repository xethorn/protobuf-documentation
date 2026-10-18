package generator_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/proto-doc/pkg/factory"
	"github.com/user/proto-doc/pkg/generator"
	"github.com/user/proto-doc/pkg/parser"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestHTMLGenerator(t *testing.T) {
	t.Run("Context: HTML Rendering", func(t *testing.T) {
		t.Run("Context: Custom Color, Theme, and Title Configuration", func(t *testing.T) {
			var config generator.Config
			var rootDoc *parser.RootResult

			// beforeEach fresh initialization
			setup := func() {
				config = generator.Config{
					Title:        "Payments Refund Service API",
					PrimaryColor: "#10b981",
					Theme:        "dark",
					FontFamily:   "Inter, sans-serif",
				}

				fct := factory.NewDescriptorFactory(99)
				fileProto := fct.CreateFileDescriptor("payments.test", 1, 2, 1)

				fds := &descriptorpb.FileDescriptorSet{
					File: []*descriptorpb.FileDescriptorProto{fileProto},
				}

				p := parser.NewParser(fds)
				rootDoc = p.ParseAll()
			}

			t.Run("Should confirm custom title and CSS properties in HTML output", func(t *testing.T) {
				setup()
				gen := generator.NewHTMLGenerator(config)
				htmlStr, err := gen.Render(rootDoc)
				require.NoError(t, err)

				assert.Contains(t, htmlStr, "<title>Payments Refund Service API</title>")
				assert.Contains(t, htmlStr, "--primary-color: #10b981")
				assert.Contains(t, htmlStr, "data-theme=\"dark\"")
				assert.Contains(t, htmlStr, "Inter, sans-serif")
			})

			t.Run("Should confirm Cmd + K modal search engine presence and indexing script", func(t *testing.T) {
				setup()
				gen := generator.NewHTMLGenerator(config)
				htmlStr, err := gen.Render(rootDoc)
				require.NoError(t, err)

				assert.True(t, strings.Contains(htmlStr, "id=\"search-modal\"") || strings.Contains(htmlStr, "search-modal"))
				assert.Contains(t, htmlStr, "⌘K")
				assert.Contains(t, htmlStr, "toggleSearchModal")
				assert.Contains(t, htmlStr, "const searchIndex = [];")
			})
		})
	})
}
