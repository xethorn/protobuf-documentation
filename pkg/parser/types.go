package parser

// Represents the root JSON documentation output result structure.
type RootResult struct {
	Version string     `json:"version"`
	Files   []*FileDoc `json:"files"`
}

// Represents extracted source code documentation comments (leading, trailing, detached).
type Comments struct {
	Leading     string   `json:"leading"`
	Trailing    string   `json:"trailing"`
	Detached    []string `json:"detached"`
	Description string   `json:"description"`
}

// Represents file-level documentation and AST declarations for a .proto file.
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

// Represents a gRPC service definition in a .proto file.
type ServiceDoc struct {
	Name     string       `json:"name"`
	FullName string       `json:"full_name"`
	Comments Comments     `json:"comments"`
	Methods  []*MethodDoc `json:"methods"`
}

// Represents an RPC method definition inside a gRPC service.
type MethodDoc struct {
	Name            string   `json:"name"`
	InputType       string   `json:"input_type"`
	OutputType      string   `json:"output_type"`
	ClientStreaming bool     `json:"client_streaming"`
	ServerStreaming bool     `json:"server_streaming"`
	Comments        Comments `json:"comments"`
}

// Represents a Protobuf message definition.
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

// Represents a oneof declaration inside a message.
type OneofDoc struct {
	Name     string   `json:"name"`
	Comments Comments `json:"comments"`
}

// Represents a field or extension definition inside a message or file.
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

// Represents an extension definition in a file or message scope.
type ExtensionDoc = FieldDoc

// Represents an Enum definition.
type EnumDoc struct {
	Name     string          `json:"name"`
	FullName string          `json:"full_name"`
	Comments Comments        `json:"comments"`
	Values   []*EnumValueDoc `json:"values"`
}

// Represents a single value declaration inside an Enum.
type EnumValueDoc struct {
	Name     string   `json:"name"`
	Number   int32    `json:"number"`
	Comments Comments `json:"comments"`
}
