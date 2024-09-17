#!/bin/bash
go get
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/ruijiesms_linux_amd64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o release/ruijiesms_linux_arm64 main.go
CGO_ENABLED=0 GOOS=linux GOARCH=loong64 go build -o release/ruijiesms_linux_loong64 main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o release/ruijiesms_windows_amd64.exe main.go
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o release/ruijiesms_windows_arm64.exe main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release/ruijiesms_darwin_amd64 main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o release/ruijiesms_darwin_arm64 main.go