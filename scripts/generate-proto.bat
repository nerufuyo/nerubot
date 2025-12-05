@echo off
REM Script to generate gRPC code from proto files (Windows)

echo Generating gRPC code from proto files...

REM Check if protoc is installed
where protoc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: protoc not found. Please install Protocol Buffers compiler.
    echo   Download from: https://github.com/protocolbuffers/protobuf/releases
    echo   Add to PATH after installation
    exit /b 1
)

REM Check if Go plugins are installed
where protoc-gen-go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Installing protoc-gen-go...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
)

where protoc-gen-go-grpc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Installing protoc-gen-go-grpc...
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
)

REM Generate code for each proto file
for %%f in (api\proto\*.proto) do (
    echo Processing %%f...
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative %%f
)

echo gRPC code generation completed!
echo Generated files are in api\proto\
