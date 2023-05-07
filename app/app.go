/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-07 22:09:30
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-07 23:57:34
 * @FilePath: /gmvpn/app/app.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package app

import (
	"log"
	"os"

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
	Config  *config.Config
	Version string
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
	if !app.Config.ServerMode {
		app.Config.LocalGateway = netutil.DiscoverGateway(true)
		app.Config.LocalGatewayv6 = netutil.DiscoverGateway(false)
	}
	app.Config.BufferSize = 64 * 1024
	if !parseConfig(app.Config) {
		os.Exit(1)
	}

	app.Iface = tun.CreateTun(*app.Config)
	log.Printf("initialized config: %+v", app.Config)
	netutil.PrintStats(app.Config.Verbose, app.Config.ServerMode)
}

func parseConfig(cfg *config.Config) bool {

	return true
}

// StartApp starts the app
func (app *App) StartApp() {
	if app.Config.ServerMode {
		tls.StartServer(app.Iface, *app.Config)
	} else {
		tls.StartClient(app.Iface, *app.Config)
	}
}

// StopApp stops the app
func (app *App) StopApp() {
	tun.ResetRoute(*app.Config)
	app.Iface.Close()
	log.Println("gmvpn stopped")
}
