package parser_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/proto-doc/pkg/factory"
	"github.com/user/proto-doc/pkg/parser"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func findRepoRoot(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func parseProtoFixture(t *testing.T, relPath string) *parser.RootResult {
	root := findRepoRoot(t)
	fullPath := filepath.Join(root, relPath)

	tempFile, err := os.CreateTemp("", "descriptor_*.pb")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	cmd := exec.Command("protoc", "-I"+root, "--descriptor_set_out="+tempFile.Name(), "--include_imports", "--include_source_info", fullPath)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err = cmd.Run()
	require.NoError(t, err, "protoc failed: %s", errBuf.String())

	data, err := os.ReadFile(tempFile.Name())
	require.NoError(t, err)

	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(data, fds)
	require.NoError(t, err)

	p := parser.NewParser(fds)
	return p.ParseAll()
}

func TestParser(t *testing.T) {
	t.Run("Context: Comment Cleaning Unit Tests", func(t *testing.T) {
		t.Run("Should correctly parse complex block comments with /**, */, *, and //", func(t *testing.T) {
			input := "/**\n * // this is cool\n * another detailed line\n */"
			cleaned := parser.CleanComment(input)
			assert.Equal(t, "this is cool\nanother detailed line", cleaned)
		})
	})

	t.Run("Context: Static Proto Fixture Parsing", func(t *testing.T) {
		t.Run("Context: Parsing sample.proto and imported dependencies", func(t *testing.T) {
			var root *parser.RootResult

			// beforeEach setup isolation
			setup := func() {
				root = parseProtoFixture(t, "tests/fixtures/sample.proto")
			}

			t.Run("Should confirm version and file presence", func(t *testing.T) {
				setup()
				assert.Equal(t, "1.0", root.Version)
				require.NotEmpty(t, root.Files)
			})

			t.Run("Should confirm service and RPC method doc extraction", func(t *testing.T) {
				setup()
				var sampleFile *parser.FileDoc
				for _, f := range root.Files {
					if filepath.Base(f.Name) == "sample.proto" {
						sampleFile = f
						break
					}
				}
				require.NotNil(t, sampleFile)
				assert.Equal(t, "payments.core", sampleFile.Package)

				require.Len(t, sampleFile.Services, 1)
				svc := sampleFile.Services[0]
				assert.Equal(t, "RefundsService", svc.Name)
				assert.Contains(t, svc.Comments.Description, "Represents service managing payment refund operations")

				require.Len(t, svc.Methods, 5)
				retrieveMethod := svc.Methods[0]
				assert.Equal(t, "RetrieveRefund", retrieveMethod.Name)
				assert.Equal(t, "payments.core.RetrieveRefundRequest", retrieveMethod.InputType)
				assert.Equal(t, "payments.core.RetrieveRefundResponse", retrieveMethod.OutputType)
				assert.Contains(t, retrieveMethod.Comments.Description, "Represents RPC method to retrieve a payment refund resource")
			})

			t.Run("Should confirm message, fields, and oneof membership extraction from imported resources", func(t *testing.T) {
				setup()
				var refundResourceFile *parser.FileDoc
				for _, f := range root.Files {
					if filepath.Base(f.Name) == "refund.proto" {
						refundResourceFile = f
						break
					}
				}
				require.NotNil(t, refundResourceFile)

				var paymentRefund *parser.MessageDoc
				for _, msg := range refundResourceFile.Messages {
					if msg.Name == "PaymentRefund" {
						paymentRefund = msg
						break
					}
				}
				require.NotNil(t, paymentRefund)
				assert.Contains(t, paymentRefund.Comments.Description, "Represents a refund of a payment transaction")

				// Check field
				var idField *parser.FieldDoc
				for _, f := range paymentRefund.Fields {
					if f.Name == "id" {
						idField = f
						break
					}
				}
				require.NotNil(t, idField)
				assert.Equal(t, "string", idField.Type)
				assert.Equal(t, int32(1), idField.Number)
				assert.Contains(t, idField.Comments.Description, "Represents unique system identifier")

				// Check oneof
				var cardField *parser.FieldDoc
				for _, f := range paymentRefund.Fields {
					if f.Name == "card_last_four" {
						cardField = f
						break
					}
				}
				require.NotNil(t, cardField)
				assert.Equal(t, "destination", cardField.OneofName)
			})

			t.Run("Should confirm Enum and EnumValue extraction", func(t *testing.T) {
				setup()
				var refundResourceFile *parser.FileDoc
				for _, f := range root.Files {
					if filepath.Base(f.Name) == "refund.proto" {
						refundResourceFile = f
						break
					}
				}
				require.NotNil(t, refundResourceFile)

				require.Len(t, refundResourceFile.Enums, 2)
				enumDoc := refundResourceFile.Enums[0]
				assert.Equal(t, "PaymentRefundStatus", enumDoc.Name)
				assert.Len(t, enumDoc.Values, 4)

				pendingVal := enumDoc.Values[1]
				assert.Equal(t, "PENDING", pendingVal.Name)
				assert.Equal(t, int32(1), pendingVal.Number)
				assert.Contains(t, pendingVal.Comments.Description, "Represents refund request currently pending")
			})
		})
	})

	t.Run("Context: Stress Testing via Descriptor Factory", func(t *testing.T) {
		var fct *factory.DescriptorFactory

		setup := func() {
			fct = factory.NewDescriptorFactory(42)
		}

		t.Run("Should successfully parse dynamically generated file descriptor", func(t *testing.T) {
			setup()
			fileProto := fct.CreateFileDescriptor("test.factory.pkg", 3, 5, 2)

			fds := &descriptorpb.FileDescriptorSet{
				File: []*descriptorpb.FileDescriptorProto{fileProto},
			}

			p := parser.NewParser(fds)
			result := p.ParseAll()

			assert.Equal(t, "1.0", result.Version)
			require.Len(t, result.Files, 1)

			fileDoc := result.Files[0]
			assert.Equal(t, "test.factory.pkg", fileDoc.Package)
			assert.Len(t, fileDoc.Services, 3)
			assert.Len(t, fileDoc.Messages, 5)
			assert.Len(t, fileDoc.Enums, 2)

			// Assert field comments inside generated messages
			msg0 := fileDoc.Messages[0]
			require.NotEmpty(t, msg0.Fields)
			assert.Contains(t, msg0.Fields[0].Comments.Description, "Field documentation for field_0")
		})
	})
}
