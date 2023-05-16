/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-07 22:09:30
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-16 00:23:28
 * @FilePath: /gmvpn/app/app.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package app

import (
	"encoding/json"
	"log"
	"net"
	"strconv"

	"gmvpn/common"
	"gmvpn/tun"

	"gmvpn/common/config"
	"gmvpn/common/netutil"
	tls "gmvpn/gmtls"

	"github.com/net-byte/water"
)

var _banner = `
 _ _                    _ __          
/ _| |  _ __    __ __   | |_ \  _ _    
\__, | | '  \   \ V /   | .__/ | ' \   
|___/  |_|_|_|  _\_/_   |_|__  |_||_|   
A simple SSLVPN for gm written in Go.
%s
`
var _srcUrl = "https://github.com/lucas-sss/gmvpn"

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

	s, _ := json.MarshalIndent(app.Config, "", "\t")
	log.Printf("initialized app config: \n%s", s)
	//校验配置参数
	if !checkConfig(app.Config) {
		log.Panicln("illegal app config")
	}

	// 如果时服务端生产ip池
	if app.Config.ServerMode {
		//解析虚拟网络地址
		_, net, _ := net.ParseCIDR(app.Config.CIDR)
		first, end := netutil.IpNetRange(net)
		ipv4FirstInt := netutil.IPv4ToInt(first)
		ipv4EndInt := netutil.IPv4ToInt(end)
		sVip, _ := netutil.IntToIPv4(ipv4FirstInt + 1)
		vipPoolBegin, _ := netutil.IntToIPv4(ipv4FirstInt + 2)
		vipPoolEnd, _ := netutil.IntToIPv4(ipv4EndInt - 1)
		vipPoolSize := ipv4EndInt - ipv4FirstInt - 2
		log.Printf("vip info: svip: %s, vip pool begin: %s, vipPool end: %s, vip pool size: %d", sVip, vipPoolBegin, vipPoolEnd, vipPoolSize)

		//初始化虚拟ip池
		vipList := make([]string, vipPoolSize)
		app.Config.VipPool = common.NewRWMutexMap(int(vipPoolSize))

		for i := ipv4FirstInt + 2; i < netutil.IPv4ToInt(end); i++ {
			vip, _ := netutil.IntToIPv4(i)
			vipInfo := new(config.VipInfo)
			vipInfo.Used = false

			vipList = append(vipList, vip)
			app.Config.VipPool.Set(vip, *vipInfo)
		}
		app.Config.VipList = vipList
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
	mask, _ := net.Mask.Size()
	tunCfg := &tun.TunConfig{}
	first, _ := netutil.IpNetRange(net)
	svip, _ := netutil.IntToIPv4(netutil.IPv4ToInt(first) + 1)

	tunCfg.Device = appCfg.Device
	tunCfg.Mask = mask
	tunCfg.SVip = svip //第一个地址默认为服务端虚拟ip
	tunCfg.Cidr = svip + "/" + strconv.Itoa(mask)
	tunCfg.PrefixLen = mask
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

// 添加iptable自动snat转换
func addIptablesSnat(cidr, cidr6 string, rs, r6s []string) {
	if rs != nil {
		existSnat := netutil.CheckExistSNat()
		for _, route := range rs {
			if _, ok := existSnat[route]; ok {
				log.Printf("already exist ipv4 snat %s, ignore processing\n", route)
				continue
			}
			//iptables -t nat ${op} POSTROUTING -s ${s} ${d} -j MASQUERADE
			netutil.ExecCmd("iptables", "-t", "nat", "-I", "POSTROUTING", "-s", cidr, route, "-j", "MASQUERADE")
		}
	}
}

// StartApp starts the app
func (app *App) StartApp() {
	if app.Config.ServerMode {
		iface, tunCfg := serverCreateTun(app.Config)
		if app.Config.AutoSnat { //配置了自动snat转换
			addIptablesSnat(tunCfg.Cidr, tunCfg.Cidr6, app.Config.Route, app.Config.Route6)
		}
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
