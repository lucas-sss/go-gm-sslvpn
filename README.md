## 一款基于go实现的支持国密算法的sslvpn



# 特性
* 支持gmtls

# 用法

```
Usage of ./gmvpn-linux-amd64:
  -autonat
    	open iptables auto snat on server mode (default true)
  -ca string
    	gmtls ca file path (default "./certs/ca.crt")
  -certificate string
    	tls certificate file path (default "./certs/server.pem")
  -cidr string
    	tun interface cidr (default "10.8.8.0/24")
  -cidr6 string
    	tun interface ipv6 cidr (default "fced:9999::9999/64")
  -cipher string
    	tls cipher suites (default "SM2_WITH_SM4_SM3")
  -compress
    	enable data compression
  -dev string
    	device name
  -enccert string
    	gmtls enc cert file path (default "./certs/enccert.crt")
  -enckey string
    	gmtls enc key file path (default "./certs/enckey.key")
  -g	client global mode
  -isv
    	tls insecure skip verify
  -local string
    	bind to local address (default ":3001")
  -mtu int
    	tun mtu (default 1500)
  -privatekey string
    	tls certificate key file path (default "./certs/server.key")
  -remote string
    	remote server address (default ":3001")
  -route value
    	push ipv4 route to client
  -route6 value
    	push ipv6 route to client
  -s	server mode
  -signcert string
    	gmtls sign cert file path (default "./certs/signcert.crt")
  -signkey string
    	gmtls sign key file path (default "./certs/signkey.key")
  -sni string
    	tls handshake sni
  -t int
    	dial timeout in seconds (default 30)
```

## 编译

```
./build.sh
```

## 运行

### Linux客户端

```
sudo ./vtun-linux-amd64 -remote server-addr:port -isv

```

### MacOS客户端

```
sudo ./vtun-drawin-amd64 -remote server-addr:port -isv

```

### Windows客户端
在windows上使用，你必须下载[wintun.dll](https://www.wintun.net/)文件并且把它放到当前应用目录下。  
用管理员权限打开powershell并运行命令:
```
.\vtun-drawin-amd64 -remote server-addr:port -isv

```

### Linux服务端

```
sudo ./vtun-linux-amd64 -s -cidr :3001 

```

