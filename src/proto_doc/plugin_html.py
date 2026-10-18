"""
protoc plugin for generating HTML documentation directly from protoc.
"""

import sys
import json
from google.protobuf.compiler import plugin_pb2
from google.protobuf.descriptor_pb2 import FileDescriptorSet
from proto_doc.parser import ProtoDocParser
from proto_doc.generator import HTMLGenerator


def main():
    data = sys.stdin.buffer.read()
    if not data:
        return

    request = plugin_pb2.CodeGeneratorRequest()
    request.ParseFromString(data)

    fds = FileDescriptorSet()
    fds.file.extend(request.proto_file)

    parser = ProtoDocParser(fds)
    doc_data = parser.parse_all()

    # Parse parameter options e.g. title=Square API,primary_color=#006aff
    config = {}
    output_filename = "index.html"
    if request.parameter:
        parts = request.parameter.split(",")
        for part in parts:
            if "=" in part:
                k, v = part.split("=", 1)
                k = k.strip()
                v = v.strip()
                if k == "out":
                    output_filename = v
                else:
                    config[k] = v

    generator = HTMLGenerator(config=config)
    html_content = generator.render(doc_data)

    response = plugin_pb2.CodeGeneratorResponse()
    response.supported_features = plugin_pb2.CodeGeneratorResponse.FEATURE_PROTO3_OPTIONAL

    f = response.file.add()
    f.name = output_filename
    f.content = html_content

    sys.stdout.buffer.write(response.SerializeToString())


if __name__ == "__main__":
    main()
