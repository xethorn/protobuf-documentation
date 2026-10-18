package parser

// RootResult represents the root document output.
type RootResult struct {
	Version string     `json:"version"`
	Files   []*FileDoc `json:"files"`
}

// Comments holds extracted source code documentation comments.
type Comments struct {
	Leading     string   `json:"leading"`
	Trailing    string   `json:"trailing"`
	Detached    []string `json:"detached"`
	Description string   `json:"description"`
}

// FileDoc holds file-level documentation and AST declarations.
type FileDoc struct {
	Name         string          `json:"name"`
	Package      string          `json:"package"`
	Syntax       string          `json:"syntax"`
	Dependencies []string        `json:"dependencies"`
	Comments     Comments        `json:"comments"`
	Services     []*ServiceDoc   `json:"services"`
	Messages     []*MessageDoc   `json:"messages"`
	Enums        []*EnumDoc      `json:"enums"`
	Extensions   []*ExtensionDoc `json:"extensions"`
}

// ServiceDoc represents a gRPC service definition.
type ServiceDoc struct {
	Name     string       `json:"name"`
	FullName string       `json:"full_name"`
	Comments Comments     `json:"comments"`
	Methods  []*MethodDoc `json:"methods"`
}

// MethodDoc represents an RPC method definition.
type MethodDoc struct {
	Name            string   `json:"name"`
	InputType       string   `json:"input_type"`
	OutputType      string   `json:"output_type"`
	ClientStreaming bool     `json:"client_streaming"`
	ServerStreaming bool     `json:"server_streaming"`
	Comments        Comments `json:"comments"`
}

// MessageDoc represents a Protobuf message definition.
type MessageDoc struct {
	Name           string          `json:"name"`
	FullName       string          `json:"full_name"`
	Comments       Comments        `json:"comments"`
	IsMapEntry     bool            `json:"is_map_entry"`
	Oneofs         []*OneofDoc     `json:"oneofs"`
	Fields         []*FieldDoc     `json:"fields"`
	NestedMessages []*MessageDoc   `json:"nested_messages"`
	NestedEnums    []*EnumDoc      `json:"nested_enums"`
	Extensions     []*ExtensionDoc `json:"extensions"`
}

// OneofDoc represents a oneof declaration inside a message.
type OneofDoc struct {
	Name     string   `json:"name"`
	Comments Comments `json:"comments"`
}

// FieldDoc represents a field or extension in a message/file.
type FieldDoc struct {
	Name         string   `json:"name"`
	Number       int32    `json:"number"`
	Type         string   `json:"type"`
	TypeName     string   `json:"type_name"`
	TypeDisplay  string   `json:"type_display"`
	Label        string   `json:"label"`
	DefaultValue string   `json:"default_value,omitempty"`
	JSONName     string   `json:"json_name"`
	Extendee     string   `json:"extendee,omitempty"`
	OneofName    string   `json:"oneof_name,omitempty"`
	Comments     Comments `json:"comments"`
}

// ExtensionDoc is an alias for FieldDoc for top-level or nested extensions.
type ExtensionDoc = FieldDoc

// EnumDoc represents an Enum definition.
type EnumDoc struct {
	Name     string          `json:"name"`
	FullName string          `json:"full_name"`
	Comments Comments        `json:"comments"`
	Values   []*EnumValueDoc `json:"values"`
}

// EnumValueDoc represents a single value inside an Enum.
type EnumValueDoc struct {
	Name     string   `json:"name"`
	Number   int32    `json:"number"`
	Comments Comments `json:"comments"`
}
