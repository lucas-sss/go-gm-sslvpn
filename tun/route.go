/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-13 12:33:20
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-14 10:44:29
 * @FilePath: /gmvpn/tun/route.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tun

import (
	"gmvpn/common/netutil"
	"log"
	"runtime"
)

type RouteConfig struct {
	Device    string   `json:"device"`
	Ipv4Route []string `json:"ipv4Route"`
	Ipv6Route []string `json:"ipv6Route"`
	SVip      string   `json:"sVip"`
	SVip6     string   `json:"sVip6"`
}

func ConfigRoute(cfg RouteConfig) {
	if cfg.Ipv4Route != nil {
		for _, route4 := range cfg.Ipv4Route {
			addIpv4Route(cfg.Device, route4, cfg.SVip)
		}
	}
	if cfg.Ipv6Route != nil {
		for _, route6 := range cfg.Ipv6Route {
			addIpv6Route(cfg.Device, route6, cfg.SVip6)
		}
	}
}

func addIpv4Route(device, route, sVip string) {
	if sVip == "" {
		log.Println("add ipv4 route fail, sVip is empty")
		return
	}
	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "route", "add", route, "dev", device)
	} else if os == "darwin" {
		netutil.ExecCmd("route", "add", route, "-interface", device)
	} else if os == "windows" {
		netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", sVip, "metric", "6")
	} else {
		log.Printf("not support os %v", os)
	}
}

func addIpv6Route(device, route, sVip6 string) {
	if sVip6 == "" {
		return
	}

	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "-6", "route", "add", route, "dev", device)
	} else if os == "darwin" {
		netutil.ExecCmd("route", "add", "-inet6", route, "-interface", device)
	} else if os == "windows" {
		//netutil.ExecCmd("cmd", "/C", "route", "-6", "add", "::/0", "mask", "::/0", sVip6, "metric", "6")
	} else {
		log.Printf("not support os %v", os)
	}
}
