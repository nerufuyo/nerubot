# Proto Code Generation Guide

This guide explains how to generate Go code from Protocol Buffer definitions for the NeruBot microservices.

## Prerequisites

### 1. Install Protocol Buffer Compiler (protoc)

#### Windows
Download the latest protoc compiler:
1. Visit: https://github.com/protocolbuffers/protobuf/releases
2. Download: `protoc-<version>-win64.zip`
3. Extract to `C:\protoc\`
4. Add `C:\protoc\bin` to your PATH environment variable

Verify installation:
```powershell
protoc --version
# Should output: libprotoc 3.x.x or higher
```

#### Linux/macOS
```bash
# Ubuntu/Debian
sudo apt-get install -y protobuf-compiler

# macOS
brew install protobuf

# Verify
protoc --version
```

### 2. Install Go Plugins (Already Installed)

The following Go plugins are required and have been installed:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Verify Go plugins are in PATH:
```powershell
# Windows PowerShell
$env:PATH -split ';' | Select-String "go\\bin"

# Should show path containing Go binaries
# Typically: C:\Users\<username>\go\bin
```

If not in PATH, add `%USERPROFILE%\go\bin` (Windows) or `$HOME/go/bin` (Linux/macOS) to your PATH.

## Generating Proto Code

### Option 1: Using the Generation Script (Recommended)

#### Windows PowerShell
```powershell
.\scripts\generate-proto.bat
```

#### Linux/macOS
```bash
chmod +x scripts/generate-proto.sh
./scripts/generate-proto.sh
```

### Option 2: Manual Generation

Generate code for all services:

```powershell
# Create output directory
New-Item -ItemType Directory -Force -Path api/proto/gen

# Generate Music service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/music.proto

# Generate Confession service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/confession.proto

# Generate Roast service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/roast.proto

# Generate Chatbot service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/chatbot.proto

# Generate News service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/news.proto

# Generate Whale service
protoc --go_out=. --go_opt=paths=source_relative `
       --go-grpc_out=. --go-grpc_opt=paths=source_relative `
       api/proto/whale.proto
```

## Generated Files

After successful generation, you should see:
```
api/proto/
├── music.proto
├── music.pb.go           # Generated message types
├── music_grpc.pb.go      # Generated gRPC service code
├── confession.proto
├── confession.pb.go
├── confession_grpc.pb.go
├── roast.proto
├── roast.pb.go
├── roast_grpc.pb.go
├── chatbot.proto
├── chatbot.pb.go
├── chatbot_grpc.pb.go
├── news.proto
├── news.pb.go
├── news_grpc.pb.go
├── whale.proto
├── whale.pb.go
└── whale_grpc.pb.go
```

## Troubleshooting

### Error: "protoc: command not found"
- Ensure protoc is installed and in your PATH
- Restart your terminal after adding to PATH
- On Windows, check: `$env:PATH -split ';' | Select-String protoc`

### Error: "protoc-gen-go: program not found or is not executable"
- Install Go plugins: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- Ensure `$GOPATH/bin` or `$HOME/go/bin` is in your PATH
- On Windows: `$env:PATH += ";$env:USERPROFILE\go\bin"`

### Error: "protoc-gen-go-grpc: program not found or is not executable"
- Install gRPC plugin: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- Verify plugin location: `Get-Command protoc-gen-go-grpc` (Windows) or `which protoc-gen-go-grpc` (Linux)

### Error: Import path not found
- Ensure you're running protoc from the project root directory
- Check that proto files have correct `option go_package` declarations

## Next Steps

After generating proto code:

1. **Implement gRPC Servers**: Add server implementations in each service
2. **Add gRPC Clients**: Update API Gateway to call backend services
3. **Test Services**: Run individual services and test gRPC communication
4. **Deploy**: Use Docker Compose or Railway for deployment

See `docs/DEVELOPMENT.md` for complete implementation guide.

## Resources

- Protocol Buffers Documentation: https://protobuf.dev/
- gRPC Go Tutorial: https://grpc.io/docs/languages/go/
- NeruBot Architecture: See `ARCHITECTURE.md` and `projects/project_plan.md`
