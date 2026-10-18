package parser

import (
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

var fieldTypeNames = map[descriptorpb.FieldDescriptorProto_Type]string{
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:   "double",
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT:    "float",
	descriptorpb.FieldDescriptorProto_TYPE_INT64:    "int64",
	descriptorpb.FieldDescriptorProto_TYPE_UINT64:   "uint64",
	descriptorpb.FieldDescriptorProto_TYPE_INT32:   "int32",
	descriptorpb.FieldDescriptorProto_TYPE_FIXED64: "fixed64",
	descriptorpb.FieldDescriptorProto_TYPE_FIXED32: "fixed32",
	descriptorpb.FieldDescriptorProto_TYPE_BOOL:    "bool",
	descriptorpb.FieldDescriptorProto_TYPE_STRING:  "string",
	descriptorpb.FieldDescriptorProto_TYPE_GROUP:   "group",
	descriptorpb.FieldDescriptorProto_TYPE_MESSAGE: "message",
	descriptorpb.FieldDescriptorProto_TYPE_BYTES:   "bytes",
	descriptorpb.FieldDescriptorProto_TYPE_UINT32:  "uint32",
	descriptorpb.FieldDescriptorProto_TYPE_ENUM:    "enum",
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED32: "sfixed32",
	descriptorpb.FieldDescriptorProto_TYPE_SFIXED64: "sfixed64",
	descriptorpb.FieldDescriptorProto_TYPE_SINT32:  "sint32",
	descriptorpb.FieldDescriptorProto_TYPE_SINT64:  "sint64",
}

var labelNames = map[descriptorpb.FieldDescriptorProto_Label]string{
	descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL: "optional",
	descriptorpb.FieldDescriptorProto_LABEL_REQUIRED: "required",
	descriptorpb.FieldDescriptorProto_LABEL_REPEATED: "repeated",
}

type Parser struct {
	fileSet *descriptorpb.FileDescriptorSet
}

func NewParser(fileSet *descriptorpb.FileDescriptorSet) *Parser {
	return &Parser{fileSet: fileSet}
}

func (p *Parser) ParseAll() *RootResult {
	var files []*FileDoc
	for _, f := range p.fileSet.GetFile() {
		files = append(files, p.ParseFile(f))
	}
	return &RootResult{
		Version: "1.0",
		Files:   files,
	}
}

func (p *Parser) ParseFile(fileProto *descriptorpb.FileDescriptorProto) *FileDoc {
	commentsMap := make(map[string]*descriptorpb.SourceCodeInfo_Location)
	if fileProto.SourceCodeInfo != nil {
		for _, loc := range fileProto.SourceCodeInfo.Location {
			key := PathKey(loc.Path)
			commentsMap[key] = loc
		}
	}

	fileComments := ExtractComments(commentsMap[""])

	var messages []*MessageDoc
	// FileDescriptorProto message_type field tag is 4
	for i, msg := range fileProto.GetMessageType() {
		messages = append(messages, p.parseMessage(msg, []int32{4, int32(i)}, commentsMap, fileProto.GetPackage(), ""))
	}

	var enums []*EnumDoc
	// FileDescriptorProto enum_type field tag is 5
	for i, enum := range fileProto.GetEnumType() {
		enums = append(enums, p.parseEnum(enum, []int32{5, int32(i)}, commentsMap, fileProto.GetPackage(), ""))
	}

	var services []*ServiceDoc
	// FileDescriptorProto service field tag is 6
	for i, service := range fileProto.GetService() {
		services = append(services, p.parseService(service, []int32{6, int32(i)}, commentsMap, fileProto.GetPackage()))
	}

	var extensions []*ExtensionDoc
	// FileDescriptorProto extension field tag is 7
	for i, ext := range fileProto.GetExtension() {
		extensions = append(extensions, p.parseField(ext, []int32{7, int32(i)}, commentsMap, nil))
	}

	syntax := fileProto.GetSyntax()
	if syntax == "" {
		syntax = "proto2"
	}

	return &FileDoc{
		Name:         fileProto.GetName(),
		Package:      fileProto.GetPackage(),
		Syntax:       syntax,
		Dependencies: fileProto.GetDependency(),
		Comments:     fileComments,
		Services:     services,
		Messages:     messages,
		Enums:        enums,
		Extensions:   extensions,
	}
}

func (p *Parser) parseMessage(
	msg *descriptorpb.DescriptorProto,
	path []int32,
	commentsMap map[string]*descriptorpb.SourceCodeInfo_Location,
	pkg string,
	parentPrefix string,
) *MessageDoc {
	fullName := msg.GetName()
	if parentPrefix != "" {
		fullName = parentPrefix + "." + msg.GetName()
	} else if pkg != "" {
		fullName = pkg + "." + msg.GetName()
	}

	comments := ExtractComments(commentsMap[PathKey(path)])
	isMapEntry := msg.GetOptions().GetMapEntry()

	var oneofs []*OneofDoc
	// Oneof_decl tag in DescriptorProto is 8
	for i, oneof := range msg.GetOneofDecl() {
		oneofPath := append(path, 8, int32(i))
		oneofs = append(oneofs, &OneofDoc{
			Name:     oneof.GetName(),
			Comments: ExtractComments(commentsMap[PathKey(oneofPath)]),
		})
	}

	var fields []*FieldDoc
	// Field tag in DescriptorProto is 2
	for i, field := range msg.GetField() {
		fieldPath := append(path, 2, int32(i))
		fields = append(fields, p.parseField(field, fieldPath, commentsMap, msg))
	}

	var nestedMessages []*MessageDoc
	// Nested_type tag in DescriptorProto is 3
	for i, nestedMsg := range msg.GetNestedType() {
		nestedPath := append(path, 3, int32(i))
		nestedMessages = append(nestedMessages, p.parseMessage(nestedMsg, nestedPath, commentsMap, pkg, fullName))
	}

	var nestedEnums []*EnumDoc
	// Enum_type tag in DescriptorProto is 4
	for i, nestedEnum := range msg.GetEnumType() {
		nestedPath := append(path, 4, int32(i))
		nestedEnums = append(nestedEnums, p.parseEnum(nestedEnum, nestedPath, commentsMap, pkg, fullName))
	}

	var extensions []*ExtensionDoc
	// Extension tag in DescriptorProto is 6
	for i, ext := range msg.GetExtension() {
		extPath := append(path, 6, int32(i))
		extensions = append(extensions, p.parseField(ext, extPath, commentsMap, msg))
	}

	return &MessageDoc{
		Name:           msg.GetName(),
		FullName:       fullName,
		Comments:       comments,
		IsMapEntry:     isMapEntry,
		Oneofs:         oneofs,
		Fields:         fields,
		NestedMessages: nestedMessages,
		NestedEnums:    nestedEnums,
		Extensions:     extensions,
	}
}

func (p *Parser) parseField(
	field *descriptorpb.FieldDescriptorProto,
	path []int32,
	commentsMap map[string]*descriptorpb.SourceCodeInfo_Location,
	parentMsg *descriptorpb.DescriptorProto,
) *FieldDoc {
	comments := ExtractComments(commentsMap[PathKey(path)])
	typeName := fieldTypeNames[field.GetType()]
	if typeName == "" {
		typeName = "unknown"
	}

	typeDisplay := typeName
	rawTypeName := strings.TrimPrefix(field.GetTypeName(), ".")
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE || field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_ENUM {
		typeDisplay = rawTypeName
	}

	labelName := labelNames[field.GetLabel()]
	if labelName == "" {
		labelName = "optional"
	}

	jsonName := field.GetJsonName()
	if jsonName == "" {
		jsonName = field.GetName()
	}

	var oneofName string
	if field.OneofIndex != nil && parentMsg != nil && int(field.GetOneofIndex()) < len(parentMsg.GetOneofDecl()) {
		decl := parentMsg.GetOneofDecl()[field.GetOneofIndex()]
		if !strings.HasPrefix(decl.GetName(), "_") {
			oneofName = decl.GetName()
		}
	}

	extendee := strings.TrimPrefix(field.GetExtendee(), ".")

	return &FieldDoc{
		Name:         field.GetName(),
		Number:       field.GetNumber(),
		Type:         typeName,
		TypeName:     rawTypeName,
		TypeDisplay:  typeDisplay,
		Label:        labelName,
		DefaultValue: field.GetDefaultValue(),
		JSONName:     jsonName,
		Extendee:     extendee,
		OneofName:    oneofName,
		Comments:     comments,
	}
}

func (p *Parser) parseEnum(
	enum *descriptorpb.EnumDescriptorProto,
	path []int32,
	commentsMap map[string]*descriptorpb.SourceCodeInfo_Location,
	pkg string,
	parentPrefix string,
) *EnumDoc {
	fullName := enum.GetName()
	if parentPrefix != "" {
		fullName = parentPrefix + "." + enum.GetName()
	} else if pkg != "" {
		fullName = pkg + "." + enum.GetName()
	}

	comments := ExtractComments(commentsMap[PathKey(path)])

	var values []*EnumValueDoc
	// Value tag in EnumDescriptorProto is 2
	for i, val := range enum.GetValue() {
		valPath := append(path, 2, int32(i))
		values = append(values, &EnumValueDoc{
			Name:     val.GetName(),
			Number:   val.GetNumber(),
			Comments: ExtractComments(commentsMap[PathKey(valPath)]),
		})
	}

	return &EnumDoc{
		Name:     enum.GetName(),
		FullName: fullName,
		Comments: comments,
		Values:   values,
	}
}

func (p *Parser) parseService(
	service *descriptorpb.ServiceDescriptorProto,
	path []int32,
	commentsMap map[string]*descriptorpb.SourceCodeInfo_Location,
	pkg string,
) *ServiceDoc {
	fullName := service.GetName()
	if pkg != "" {
		fullName = pkg + "." + service.GetName()
	}

	comments := ExtractComments(commentsMap[PathKey(path)])

	var methods []*MethodDoc
	// Method tag in ServiceDescriptorProto is 2
	for i, method := range service.GetMethod() {
		methodPath := append(path, 2, int32(i))
		methods = append(methods, &MethodDoc{
			Name:            method.GetName(),
			InputType:       strings.TrimPrefix(method.GetInputType(), "."),
			OutputType:      strings.TrimPrefix(method.GetOutputType(), "."),
			ClientStreaming: method.GetClientStreaming(),
			ServerStreaming: method.GetServerStreaming(),
			Comments:        ExtractComments(commentsMap[PathKey(methodPath)]),
		})
	}

	return &ServiceDoc{
		Name:     service.GetName(),
		FullName: fullName,
		Comments: comments,
		Methods:  methods,
	}
}
