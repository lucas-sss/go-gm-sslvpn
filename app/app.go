/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-07 22:09:30
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-14 12:28:57
 * @FilePath: /gmvpn/app/app.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package app

import (
	"log"
	"net"
	"strconv"

	"gmvpn/tun"

	"gmvpn/common/config"
	"gmvpn/common/netutil"
	tls "gmvpn/gmtls"

	"github.com/net-byte/water"
)

var _banner = `
_                 
__ __  ___   _  _   _ _    __
\ V / |  _| | || | | ' \  ||_||
 \_/  |_|    \_,_| |_||_| || ||
                             //		 
A simple SSL VPN for gm written in Go.
%s
`
var _srcUrl = "xxxx"

// vtun app struct
type App struct {
	Version string
	Config  *config.Config
	Iface   *water.Interface
}

func NewApp(config *config.Config, version string) *App {
	return &App{
		Config:  config,
		Version: version,
	}
}

// InitConfig initializes the config
func (app *App) InitConfig() {
	log.Printf(_banner, _srcUrl)
	log.Printf("gmvpn version %s", app.Version)
	app.Config.LocalGateway = netutil.DiscoverGateway(true)
	app.Config.LocalGateway6 = netutil.DiscoverGateway(false)
	app.Config.BufferSize = 64 * 1024

	log.Printf("initialized app config: %+v", app.Config)
	//校验配置参数
	if !checkConfig(app.Config) {
		log.Panicln("illegal app config")
	}

	//TODO 如果时服务端生产ip池
	if app.Config.ServerMode {

	}
}

func serverCreateTun(appCfg *config.Config) (*water.Interface, *tun.TunConfig) {
	//生成tun网卡配置
	tunCfg := generateTunConfig(*appCfg)
	tunCfg.ServerMode = true
	//创建虚拟网卡
	iface := tun.CreatAndConfigTun(tunCfg)
	return iface, tunCfg
}

func generateTunConfig(appCfg config.Config) *tun.TunConfig {
	_, net, _ := net.ParseCIDR(appCfg.CIDR)
	prefixLen, _ := net.Mask.Size()
	tunCfg := &tun.TunConfig{}
	first, _ := netutil.IpNetRange(net)
	svip, _ := netutil.IntToIPv4(netutil.IPv4ToInt(first) + 1)

	tunCfg.Device = appCfg.Device
	tunCfg.PrefixLen = prefixLen
	tunCfg.SVip = svip //第一个地址默认为服务端虚拟ip
	tunCfg.Cidr = svip + "/" + strconv.Itoa(prefixLen)
	tunCfg.Cidr = svip + "/" + strconv.Itoa(prefixLen)
	tunCfg.Mtu = appCfg.MTU
	tunCfg.LocalGateway = appCfg.LocalGateway
	tunCfg.LocalGateway6 = appCfg.LocalGateway6

	return tunCfg
}

// 校验所有配置参数
func checkConfig(cfg *config.Config) bool {
	if cfg.Route != nil {
		for _, route := range cfg.Route {
			_, _, err := net.ParseCIDR(route)
			if err != nil {
				log.Println("route format error: ", route)
				return false
			}
		}
	}
	return true
}

// StartApp starts the app
func (app *App) StartApp() {
	if app.Config.ServerMode {
		iface, tunCfg := serverCreateTun(app.Config)
		tls.StartServer(iface, *app.Config, *tunCfg)
	} else {
		tls.StartClient(*app.Config)
	}
}

// StopApp stops the app
func (app *App) StopApp() {
	tun.ResetRoute(*app.Config)
	log.Println("gmvpn stopped")
}
