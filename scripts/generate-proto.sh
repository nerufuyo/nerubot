#!/bin/bash
# Script to generate gRPC code from proto files

echo "Generating gRPC code from proto files..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc not found. Please install Protocol Buffers compiler."
    echo "  macOS: brew install protobuf"
    echo "  Ubuntu: sudo apt-get install protobuf-compiler"
    echo "  Windows: Download from https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate code for each proto file
for proto in api/proto/*.proto; do
    echo "Processing $proto..."
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        "$proto"
done

echo "gRPC code generation completed!"
echo "Generated files are in api/proto/"
