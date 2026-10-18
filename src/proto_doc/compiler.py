"""
Compiler and AST builder utilities for compile-time or runtime proto parsing.
"""

import os
import sys
import tempfile
from typing import List, Dict, Any, Optional
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from grpc_tools import protoc
from proto_doc.parser import ProtoDocParser


def parse_proto_files(proto_files: List[str], include_paths: Optional[List[str]] = None) -> Dict[str, Any]:
    """Compiles proto files using grpc_tools.protoc into a FileDescriptorSet in memory and parses it."""
    if include_paths is None:
        include_paths = ["."]

    with tempfile.NamedTemporaryFile(suffix=".pb", delete=False) as descriptor_file:
        descriptor_path = descriptor_file.name

    try:
        args = ["protoc"]
        for inc in include_paths:
            args.append(f"-I{inc}")

        # Include standard google protobuf include dir if present in site-packages
        import grpc_tools
        grpc_tools_path = os.path.dirname(grpc_tools.__file__)
        proto_include = os.path.join(grpc_tools_path, "_proto")
        if os.path.exists(proto_include):
            args.append(f"-I{proto_include}")

        args.append(f"--descriptor_set_out={descriptor_path}")
        args.append("--include_imports")
        args.append("--include_source_info")
        args.extend(proto_files)

        res = protoc.main(args)
        if res != 0:
            raise RuntimeError(f"protoc failed with exit code {res} for files: {proto_files}")

        with open(descriptor_path, "rb") as f:
            fds = FileDescriptorSet()
            fds.ParseFromString(f.read())

        parser = ProtoDocParser(fds)
        return parser.parse_all()
    finally:
        if os.path.exists(descriptor_path):
            os.remove(descriptor_path)
