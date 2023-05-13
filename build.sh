#!/bin/bash

# Build for Windows
echo "Building for Windows"
GOOS=windows GOARCH=amd64 go build -o load_windows_amd64 main.go

# Build for macOS
echo "Building for macOS"
GOOS=darwin GOARCH=amd64 go build -o load_darwin_amd64 main.go

# Build for Linux
echo "Building for Linux"
GOOS=linux GOARCH=amd64 go build -o load_linux_amd64 main.go
