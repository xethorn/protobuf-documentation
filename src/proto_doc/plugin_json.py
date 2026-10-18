"""
protoc plugin for generating JSON documentation from proto files.
"""

import sys
import json
from google.protobuf.compiler import plugin_pb2
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from proto_doc.parser import ProtoDocParser


def main():
    # Read CodeGeneratorRequest from stdin
    data = sys.stdin.buffer.read()
    if not data:
        return

    request = plugin_pb2.CodeGeneratorRequest()
    request.ParseFromString(data)

    # Convert request.proto_file (list of FileDescriptorProto) into FileDescriptorSet
    fds = FileDescriptorSet()
    fds.file.extend(request.proto_file)

    parser = ProtoDocParser(fds)
    doc_data = parser.parse_all()

    # Create CodeGeneratorResponse
    response = plugin_pb2.CodeGeneratorResponse()

    # Enable proto3 optional support
    response.supported_features = plugin_pb2.CodeGeneratorResponse.FEATURE_PROTO3_OPTIONAL

    # Output filename defaults to proto_docs.json or parameter
    output_filename = request.parameter if request.parameter else "proto_docs.json"

    f = response.file.add()
    f.name = output_filename
    f.content = json.dumps(doc_data, indent=2)

    sys.stdout.buffer.write(response.SerializeToString())


if __name__ == "__main__":
    main()
