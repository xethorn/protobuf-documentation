# Proto Doc

A portable, high-performance Protobuf documentation tool written in Go.

`proto-doc` parses Protobuf `.proto` files into structured JSON schemas with full documentation comment extraction, and generates modern static HTML documentation inspired by clean API reference layouts.

---

## Features

- **Project 1: Protobuf to JSON Extractor** (`proto-doc parse` & `protoc-gen-docjson`)
  - Parses `CodeGeneratorRequest` from `protoc` or directly processes `.proto` files in CLI mode.
  - Extracts full `SourceCodeInfo` comments across services, RPC methods, messages, fields, oneofs, enums, enum values, and extensions.
  - Cleanly strips `/**`, `*/`, `*`, and `//` comment delimiters while preserving multi-line documentation context.

- **Project 2: Static HTML Documentation Generator** (`proto-doc html` & `protoc-gen-dochtml`)
  - Generates standalone, self-contained HTML documentation.
  - Modern, fluid, non-boxy responsive layout.
  - **⌘K Search**: Instant client-side search modal searching across services, methods, messages, fields, and enums.
  - **Customization Options**: Configurable title, light/dark theme toggle, font family, primary accent color, and custom logo URL.

---

## Installation

### From Source (Go 1.21+)

Clone the repository and build the binaries:

```bash
# Build the CLI executable
go build -o proto-doc ./cmd/proto-doc

# Build the protoc plugins
go build -o protoc-gen-docjson ./cmd/protoc-gen-docjson
go build -o protoc-gen-dochtml ./cmd/protoc-gen-dochtml
```

To make the binaries globally accessible, move them to your `$PATH` (e.g. `/usr/local/bin`):

```bash
cp proto-doc protoc-gen-docjson protoc-gen-dochtml /usr/local/bin/
```

---

## Usage

### 1. Command Line Interface (CLI)

#### Parse Protobuf Files to JSON

```bash
proto-doc parse -I . -o docs.json path/to/service.proto
```

#### Generate Static HTML Documentation

```bash
proto-doc html \
  -I . \
  -o docs.html \
  --title "Payment Refunds API" \
  --primary-color "#10b981" \
  --theme light \
  path/to/service.proto
```

You can also pass a previously generated JSON file directly to `proto-doc html`:

```bash
proto-doc html -o docs.html -title "Payment Refunds API" docs.json
```

---

### 2. Protoc Plugin Usage

Because `protoc-gen-docjson` and `protoc-gen-dochtml` implement the standard Protobuf compiler plugin interface, you can invoke them directly via `protoc`:

```bash
# Generate JSON documentation
protoc -I . --docjson_out=. path/to/service.proto

# Generate customizable HTML documentation
protoc -I . --dochtml_out=title=Payment Refunds API,primary_color=#10b981:. path/to/service.proto
```

---

## Development & Testing

Run the test suite:

```bash
go test -v ./...
```

The test suite includes:
- **Parser Unit Tests**: Verifying AST traversal, multi-directory import resolution, and multi-line comment cleaning (`/**`, `*/`, `*`, `//`).
- **Generator Unit Tests**: Verifying HTML template rendering, custom styling variables, and `⌘K` search modal script generation.
- **Factory Stress Tests**: Generating dynamic descriptor sets to stress test the parser.
