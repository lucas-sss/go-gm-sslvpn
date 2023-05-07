/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-07 21:55:40
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-07 23:50:59
 * @FilePath: /gmvpn/common/config/config.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package config

// The config struct
type Config struct {
	Device                    string //虚拟网卡名称
	LocalAddr                 string //服务端绑定地址
	RemoteAddr                string //客户端连接地址
	CIDR                      string //ipv4虚拟网络地址
	CIDRv6                    string //ipv6虚拟网络地址
	ServerMode                bool   //服务端模式
	GlobalMode                bool   //客户端是开启全局转发
	Compress                  bool   //是否开启压缩
	MTU                       int    //虚拟网卡mtu
	Timeout                   int
	CaPath                    string //ca证书路径
	TLSCertificateFilePath    string //rsa证书文件路径
	TLSCertificateKeyFilePath string //rsa私钥文件路径
	SignCertPath              string
	SignKeyPath               string
	EncCertPath               string
	EncKeyPath                string
	TLSSni                    string //是否开启tls sni标识（多域名识别）
	TLSInsecureSkipVerify     bool   //是否跳过证书验证
	TLSCipher                 string //tls 密码套件

	BufferSize int
	Verbose    bool
	//*********** 二次生成配置 ***********//
	//
	ServerIP   string
	ServerIPv6 string
	//网关配置
	LocalGateway   string
	LocalGatewayv6 string
	//*********** 二次生成配置 ***********//
}
