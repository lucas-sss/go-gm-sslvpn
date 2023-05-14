package tun

import (
	"log"
	"runtime"
	"strconv"

	"gmvpn/common/config"
	"gmvpn/common/netutil"

	"github.com/net-byte/water"
)

type TunConfig struct {
	ServerMode    bool   `json:"serverMode"`
	Device        string `json:"device"`
	PrefixLen     int    `json:"prefixLen"`
	Cidr          string `json:"cidr"`  //ipv4网络地址
	Cidr6         string `json:"cidr6"` //ipv6网络地址
	Mtu           int    `json:"mtu"`
	LocalGateway  string `json:"localGateway"`  //本地默认ipv4网关
	LocalGateway6 string `json:"localGateway6"` //本地默认ipv6网关

	//服务端
	SVip  string `json:"sVip"`  //服务端ipv4虚拟ip地址
	SVip6 string `json:"sVip6"` //服务端ipv6虚拟ip地址

	//客户端
	CVip       string `json:"cVip"`
	CVip6      string `json:"cVip6"`
	GlobalMode bool   `json:"globalMode"` //开启全局模式，只对客户端有效
}

func CreatAndConfigTun(tunCfg *TunConfig) (iface *water.Interface) {
	log.Printf("tun config: %+v", tunCfg)

	c := water.Config{DeviceType: water.TUN}
	c.PlatformSpecificParams = water.PlatformSpecificParams{}
	os := runtime.GOOS
	if os == "windows" {
		c.PlatformSpecificParams.Name = "vtun"
		c.PlatformSpecificParams.Network = []string{tunCfg.Cidr, tunCfg.Cidr6}
	}
	if tunCfg.Device != "" {
		c.PlatformSpecificParams.Name = tunCfg.Device
	}
	iface, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to create tun interface:", err)
	}
	log.Printf("interface created %v", iface.Name())
	tunCfg.Device = iface.Name()
	configTun(*tunCfg)
	return iface
}

func configTun(tunCfg TunConfig) {
	if tunCfg.ServerMode {
		configServerTun(tunCfg)
		return
	}
	//client tun config
	configClientTun(tunCfg)
}

func configServerTun(config TunConfig) {
	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", config.Device, "mtu", strconv.Itoa(config.Mtu))
		netutil.ExecCmd("/sbin/ip", "addr", "add", config.Cidr, "dev", config.Device)
		netutil.ExecCmd("/sbin/ip", "-6", "addr", "add", config.Cidr6, "dev", config.Device)
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", config.Device, "up")
	} else {
		log.Printf("not support os %v", os)
		return
	}
	log.Printf("interface configured %v", config.Device)
}

func configClientTun(config TunConfig) {
	physicalIface := netutil.GetInterface()
	log.Println("physicalIface: ", physicalIface)

	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", config.Device, "mtu", strconv.Itoa(config.Mtu))
		netutil.ExecCmd("/sbin/ip", "addr", "add", config.Cidr, "dev", config.Device)
		// netutil.ExecCmd("/sbin/ip", "-6", "addr", "add", config.Cidr6, "dev", config.Device)
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", config.Device, "up")
		if config.GlobalMode && physicalIface != "" {
			if config.LocalGateway != "" {
				netutil.ExecCmd("/sbin/ip", "route", "add", "0.0.0.0/1", "dev", config.Device)
				netutil.ExecCmd("/sbin/ip", "route", "add", "128.0.0.0/1", "dev", config.Device)
			}
			if config.LocalGateway6 != "" {
				netutil.ExecCmd("/sbin/ip", "-6", "route", "add", "::/1", "dev", config.Device)
			}
		}

	} else if os == "darwin" {
		netutil.ExecCmd("ifconfig", config.Device, "inet", config.CVip, config.SVip, "up")
		// netutil.ExecCmd("ifconfig", config.Device, "inet6", config.CVip6, config.SVip6, "up")
		if config.GlobalMode && physicalIface != "" {
			if config.LocalGateway != "" {
				netutil.ExecCmd("route", "add", "default", config.SVip)
				netutil.ExecCmd("route", "change", "default", config.SVip6)
				netutil.ExecCmd("route", "add", "0.0.0.0/1", "-interface", config.Device)
				netutil.ExecCmd("route", "add", "128.0.0.0/1", "-interface", config.Device)
			}
			if config.LocalGateway6 != "" {
				netutil.ExecCmd("route", "add", "-inet6", "default", config.SVip)
				netutil.ExecCmd("route", "change", "-inet6", "default", config.SVip6)
				netutil.ExecCmd("route", "add", "-inet6", "::/1", "-interface", config.Device)
			}
		}
	} else if os == "windows" {
		if config.GlobalMode && physicalIface != "" {
			if config.LocalGateway != "" {
				netutil.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
				netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", config.SVip, "metric", "6")
			}
			if config.LocalGateway6 != "" {
				netutil.ExecCmd("cmd", "/C", "route", "-6", "delete", "::/0", "mask", "::/0")
				netutil.ExecCmd("cmd", "/C", "route", "-6", "add", "::/0", "mask", "::/0", config.SVip6, "metric", "6")
			}
		}
	} else {
		log.Printf("not support os %v", os)
	}
	log.Printf("interface configured %v", config.Device)
}

// ResetRoute resets the system routes
func ResetRoute(config config.Config) {
	log.Printf("reset route before exit")
	if config.ServerMode || !config.GlobalMode {
		return
	}
	os := runtime.GOOS
	if os == "darwin" {
		if config.LocalGateway != "" {
			netutil.ExecCmd("route", "add", "default", config.LocalGateway)
			netutil.ExecCmd("route", "change", "default", config.LocalGateway)
		}
		if config.LocalGateway6 != "" {
			netutil.ExecCmd("route", "add", "-inet6", "default", config.LocalGateway6)
			netutil.ExecCmd("route", "change", "-inet6", "default", config.LocalGateway6)
		}
	} else if os == "windows" {
		serverAddrIP := netutil.LookupServerAddrIP(config.RemoteAddr)
		if serverAddrIP != nil {
			if config.LocalGateway != "" {
				netutil.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
				netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", config.LocalGateway, "metric", "6")
			}
			if config.LocalGateway6 != "" {
				netutil.ExecCmd("cmd", "/C", "route", "-6", "delete", "::/0", "mask", "::/0")
				netutil.ExecCmd("cmd", "/C", "route", "-6", "add", "::/0", "mask", "::/0", config.LocalGateway6, "metric", "6")
			}
		}
	}
}
