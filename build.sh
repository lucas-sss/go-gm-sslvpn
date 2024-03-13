#!/bin/bash
###
 # @Author: liuwei lyy9645@163.com
 # @Date: 2023-05-16 20:20:23
 # @LastEditors: liuwei lyy9645@163.com
 # @LastEditTime: 2023-07-04 10:06:26
 # @FilePath: /gmvpn/build.sh
 # @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
### 

#Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gmvpn-linux-amd64 ./main.go
#Linux arm64
#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/gmvpn-linux-arm64 ./main.go
echo "build linux success"

#Mac amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gmvpn-darwin-amd64 ./main.go
#Mac arm64
#CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o gmvpn-darwin-arm64 ./main.go
echo "build macos success"

CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o gmvpn-win-amd64.exe ./main.go
echo "build windows success"


echo "build finished!!!"
