#!/bin/bash

#Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/gmvpn-linux-amd64 ./main.go
#Linux arm64
#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/gmvpn-linux-arm64 ./main.go
echo "build linux success"

#Mac amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/gmvpn-darwin-amd64 ./main.go
#Mac arm64
#CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/gmvpn-darwin-arm64 ./main.go
echo "build macos success"

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/gmvpn-win-amd64.exe ./main.go
echo "build windows success"


echo "build finished!!!"
