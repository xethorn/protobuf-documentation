"""
CLI interface for proto-doc.
"""

import sys
import argparse
import json
from proto_doc.compiler import parse_proto_files
from proto_doc.generator import HTMLGenerator


def main():
    parser = argparse.ArgumentParser(prog="proto-doc", description="Protobuf documentation tool")
    subparsers = parser.add_subparsers(dest="command", required=True)

    # Subcommand: parse (Project 1)
    parse_parser = subparsers.add_parser("parse", help="Parse proto files into structured JSON")
    parse_parser.add_argument("proto_files", nargs="+", help="Protobuf source files")
    parse_parser.add_argument("-I", "--proto_path", action="append", default=[], help="Include paths for imports")
    parse_parser.add_argument("-o", "--output", default="-", help="Output JSON file (default: stdout)")

    # Subcommand: html (Project 2)
    html_parser = subparsers.add_parser("html", help="Generate HTML documentation from proto files or JSON")
    html_parser.add_argument("inputs", nargs="+", help="Protobuf source files OR parsed JSON documentation file")
    html_parser.add_argument("-I", "--proto_path", action="append", default=[], help="Include paths for imports")
    html_parser.add_argument("-o", "--output", default="docs.html", help="Output HTML file path (default: docs.html)")
    html_parser.add_argument("--title", default="API Documentation", help="Document title")
    html_parser.add_argument("--primary-color", default="#006aff", help="Primary accent color CSS hex code")
    html_parser.add_argument("--font-family", default="system-ui, -apple-system, sans-serif", help="Custom font family")
    html_parser.add_argument("--theme", choices=["light", "dark"], default="light", help="Default theme")
    html_parser.add_argument("--logo-url", default="", help="Custom logo URL for header")

    args = parser.parse_args()

    if args.command == "parse":
        include_paths = args.proto_path if args.proto_path else ["."]
        result = parse_proto_files(args.proto_files, include_paths=include_paths)
        json_str = json.dumps(result, indent=2)

        if args.output == "-":
            sys.stdout.write(json_str)
            sys.stdout.write("\n")
        else:
            with open(args.output, "w", encoding="utf-8") as f:
                f.write(json_str)

    elif args.command == "html":
        # Check if first input is JSON file or proto file
        first_input = args.inputs[0]
        if len(args.inputs) == 1 and first_input.endswith(".json"):
            with open(first_input, "r", encoding="utf-8") as f:
                doc_data = json.load(f)
        else:
            include_paths = args.proto_path if args.proto_path else ["."]
            doc_data = parse_proto_files(args.inputs, include_paths=include_paths)

        config = {
            "title": args.title,
            "primary_color": args.primary_color,
            "font_family": args.font_family,
            "theme": args.theme,
            "logo_url": args.logo_url,
        }

        generator = HTMLGenerator(config=config)
        html_out = generator.render(doc_data)

        with open(args.output, "w", encoding="utf-8") as f:
            f.write(html_out)
        print(f"Successfully generated HTML documentation at: {args.output}")


if __name__ == "__main__":
    main()
