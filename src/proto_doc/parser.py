"""
Protobuf SourceCodeInfo parser and JSON generator.

Parses FileDescriptorSet or FileDescriptorProto objects and extracts complete AST elements,
types, labels, options, extensions, and documentation comments into a structured JSON schema.
"""

from typing import Any, Dict, List, Optional, Tuple
from google.protobuf.descriptor_pb2 import (
    FileDescriptorProto,
    FileDescriptorSet,
    DescriptorProto,
    FieldDescriptorProto,
    EnumDescriptorProto,
    EnumValueDescriptorProto,
    ServiceDescriptorProto,
    MethodDescriptorProto,
    OneofDescriptorProto,
)


FIELD_TYPE_NAMES = {
    FieldDescriptorProto.TYPE_DOUBLE: "double",
    FieldDescriptorProto.TYPE_FLOAT: "float",
    FieldDescriptorProto.TYPE_INT64: "int64",
    FieldDescriptorProto.TYPE_UINT64: "uint64",
    FieldDescriptorProto.TYPE_INT32: "int32",
    FieldDescriptorProto.TYPE_FIXED64: "fixed64",
    FieldDescriptorProto.TYPE_FIXED32: "fixed32",
    FieldDescriptorProto.TYPE_BOOL: "bool",
    FieldDescriptorProto.TYPE_STRING: "string",
    FieldDescriptorProto.TYPE_GROUP: "group",
    FieldDescriptorProto.TYPE_MESSAGE: "message",
    FieldDescriptorProto.TYPE_BYTES: "bytes",
    FieldDescriptorProto.TYPE_UINT32: "uint32",
    FieldDescriptorProto.TYPE_ENUM: "enum",
    FieldDescriptorProto.TYPE_SFIXED32: "sfixed32",
    FieldDescriptorProto.TYPE_SFIXED64: "sfixed64",
    FieldDescriptorProto.TYPE_SINT32: "sint32",
    FieldDescriptorProto.TYPE_SINT64: "sint64",
}

LABEL_NAMES = {
    FieldDescriptorProto.LABEL_OPTIONAL: "optional",
    FieldDescriptorProto.LABEL_REQUIRED: "required",
    FieldDescriptorProto.LABEL_REPEATED: "repeated",
}


def _clean_comment(comment: Optional[str]) -> str:
    if not comment:
        return ""
    lines = comment.splitlines()
    cleaned_lines = []
    for line in lines:
        line = line.strip()
        if line.startswith("//"):
            line = line[2:].strip()
        elif line.startswith("/*"):
            line = line[2:].strip()
        elif line.endswith("*/"):
            line = line[:-2].strip()
        elif line.startswith("*"):
            line = line[1:].strip()
        cleaned_lines.append(line)
    return "\n".join(cleaned_lines).strip()


def _extract_comments(location: Optional[Any]) -> Dict[str, Any]:
    if not location:
        return {"leading": "", "trailing": "", "detached": [], "description": ""}

    leading = _clean_comment(location.leading_comments) if location.HasField("leading_comments") else ""
    trailing = _clean_comment(location.trailing_comments) if location.HasField("trailing_comments") else ""
    detached = [
        _clean_comment(c) for c in location.leading_detached_comments
    ] if location.leading_detached_comments else []

    full_description = leading
    if trailing:
        if full_description:
            full_description += "\n\n" + trailing
        else:
            full_description = trailing

    return {
        "leading": leading,
        "trailing": trailing,
        "detached": detached,
        "description": full_description,
    }


class ProtoDocParser:
    """Parses Protobuf descriptors into structured documentation dictionaries."""

    def __init__(self, file_descriptor_set: FileDescriptorSet):
        self.file_descriptor_set = file_descriptor_set

    def parse_all(self) -> Dict[str, Any]:
        files = []
        for file_proto in self.file_descriptor_set.file:
            files.append(self.parse_file(file_proto))
        return {
            "files": files,
            "version": "1.0",
        }

    def parse_file(self, file_proto: FileDescriptorProto) -> Dict[str, Any]:
        comments_map: Dict[Tuple[int, ...], Any] = {}
        if file_proto.HasField("source_code_info"):
            for loc in file_proto.source_code_info.location:
                path_key = tuple(loc.path)
                comments_map[path_key] = loc

        file_comments = _extract_comments(comments_map.get((), None))

        # Index nested messages & enums
        messages = []
        # Message_type field in FileDescriptorProto is tag 4
        for idx, msg in enumerate(file_proto.message_type):
            messages.append(self._parse_message(msg, (4, idx), comments_map, package=file_proto.package))

        enums = []
        # Enum_type field in FileDescriptorProto is tag 5
        for idx, enum in enumerate(file_proto.enum_type):
            enums.append(self._parse_enum(enum, (5, idx), comments_map, package=file_proto.package))

        services = []
        # Service field in FileDescriptorProto is tag 6
        for idx, service in enumerate(file_proto.service):
            services.append(self._parse_service(service, (6, idx), comments_map, package=file_proto.package))

        extensions = []
        # Extension field in FileDescriptorProto is tag 7
        for idx, ext in enumerate(file_proto.extension):
            extensions.append(self._parse_field(ext, (7, idx), comments_map, parent_msg=None))

        return {
            "name": file_proto.name,
            "package": file_proto.package,
            "syntax": file_proto.syntax or "proto2",
            "dependencies": list(file_proto.dependency),
            "comments": file_comments,
            "services": services,
            "messages": messages,
            "enums": enums,
            "extensions": extensions,
        }

    def _parse_message(
        self,
        msg: DescriptorProto,
        path: Tuple[int, ...],
        comments_map: Dict[Tuple[int, ...], Any],
        package: str,
        parent_prefix: str = ""
    ) -> Dict[str, Any]:
        full_name = f"{parent_prefix}.{msg.name}" if parent_prefix else (f"{package}.{msg.name}" if package else msg.name)
        comments = _extract_comments(comments_map.get(path, None))

        # Check if map entry message
        is_map_entry = msg.options.map_entry if msg.HasField("options") else False

        # Oneofs (tag 8)
        oneofs = []
        for idx, oneof in enumerate(msg.oneof_decl):
            oneof_path = path + (8, idx)
            oneof_comments = _extract_comments(comments_map.get(oneof_path, None))
            oneofs.append({
                "name": oneof.name,
                "comments": oneof_comments,
            })

        # Fields (tag 2)
        fields = []
        for idx, field in enumerate(msg.field):
            field_path = path + (2, idx)
            fields.append(self._parse_field(field, field_path, comments_map, parent_msg=msg))

        # Nested messages (tag 3)
        nested_messages = []
        for idx, nested_msg in enumerate(msg.nested_type):
            nested_path = path + (3, idx)
            nested_messages.append(self._parse_message(
                nested_msg, nested_path, comments_map, package, parent_prefix=full_name
            ))

        # Nested enums (tag 4)
        nested_enums = []
        for idx, nested_enum in enumerate(msg.enum_type):
            nested_path = path + (4, idx)
            nested_enums.append(self._parse_enum(
                nested_enum, nested_path, comments_map, package, parent_prefix=full_name
            ))

        # Message extensions (tag 6)
        extensions = []
        for idx, ext in enumerate(msg.extension):
            ext_path = path + (6, idx)
            extensions.append(self._parse_field(ext, ext_path, comments_map, parent_msg=msg))

        return {
            "name": msg.name,
            "full_name": full_name,
            "comments": comments,
            "is_map_entry": is_map_entry,
            "oneofs": oneofs,
            "fields": fields,
            "nested_messages": nested_messages,
            "nested_enums": nested_enums,
            "extensions": extensions,
        }

    def _parse_field(
        self,
        field: FieldDescriptorProto,
        path: Tuple[int, ...],
        comments_map: Dict[Tuple[int, ...], Any],
        parent_msg: Optional[DescriptorProto] = None,
    ) -> Dict[str, Any]:
        comments = _extract_comments(comments_map.get(path, None))
        type_name = FIELD_TYPE_NAMES.get(field.type, "unknown")

        type_display = type_name
        if field.type in (FieldDescriptorProto.TYPE_MESSAGE, FieldDescriptorProto.TYPE_ENUM):
            type_display = field.type_name.lstrip(".")

        oneof_name = None
        if field.HasField("oneof_index") and parent_msg and field.oneof_index < len(parent_msg.oneof_decl):
            # Exclude synthetic oneofs for proto3 optional
            oneof_decl = parent_msg.oneof_decl[field.oneof_index]
            if not oneof_decl.name.startswith("_"):
                oneof_name = oneof_decl.name

        return {
            "name": field.name,
            "number": field.number,
            "type": type_name,
            "type_name": field.type_name.lstrip("."),
            "type_display": type_display,
            "label": LABEL_NAMES.get(field.label, "optional"),
            "default_value": field.default_value if field.HasField("default_value") else None,
            "json_name": field.json_name if field.HasField("json_name") else field.name,
            "extendee": field.extendee.lstrip(".") if field.HasField("extendee") else None,
            "oneof_name": oneof_name,
            "comments": comments,
        }

    def _parse_enum(
        self,
        enum: EnumDescriptorProto,
        path: Tuple[int, ...],
        comments_map: Dict[Tuple[int, ...], Any],
        package: str,
        parent_prefix: str = ""
    ) -> Dict[str, Any]:
        full_name = f"{parent_prefix}.{enum.name}" if parent_prefix else (f"{package}.{enum.name}" if package else enum.name)
        comments = _extract_comments(comments_map.get(path, None))

        values = []
        # Value field in EnumDescriptorProto is tag 2
        for idx, val in enumerate(enum.value):
            val_path = path + (2, idx)
            val_comments = _extract_comments(comments_map.get(val_path, None))
            values.append({
                "name": val.name,
                "number": val.number,
                "comments": val_comments,
            })

        return {
            "name": enum.name,
            "full_name": full_name,
            "comments": comments,
            "values": values,
        }

    def _parse_service(
        self,
        service: ServiceDescriptorProto,
        path: Tuple[int, ...],
        comments_map: Dict[Tuple[int, ...], Any],
        package: str,
    ) -> Dict[str, Any]:
        full_name = f"{package}.{service.name}" if package else service.name
        comments = _extract_comments(comments_map.get(path, None))

        methods = []
        # Method field in ServiceDescriptorProto is tag 2
        for idx, method in enumerate(service.method):
            method_path = path + (2, idx)
            methods.append(self._parse_method(method, method_path, comments_map, package))

        return {
            "name": service.name,
            "full_name": full_name,
            "comments": comments,
            "methods": methods,
        }

    def _parse_method(
        self,
        method: MethodDescriptorProto,
        path: Tuple[int, ...],
        comments_map: Dict[Tuple[int, ...], Any],
        package: str,
    ) -> Dict[str, Any]:
        comments = _extract_comments(comments_map.get(path, None))
        return {
            "name": method.name,
            "input_type": method.input_type.lstrip("."),
            "output_type": method.output_type.lstrip("."),
            "client_streaming": method.client_streaming,
            "server_streaming": method.server_streaming,
            "comments": comments,
        }
