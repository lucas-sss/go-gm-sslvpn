package tls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	"gmvpn/common/cache"
	"gmvpn/common/config"
	"gmvpn/common/counter"
	"gmvpn/common/netutil"
	"gmvpn/tun"

	"github.com/golang/snappy"
	"github.com/net-byte/water"
	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
)

// StartServer starts the tls server
func StartServer(iface *water.Interface, appCfg config.Config, tunCfg tun.TunConfig) {
	defer func() {
		log.Println("server tun interface close")
		iface.Close()
	}()
	log.Printf("vtun tls server started on %v", appCfg.LocalAddr)

	//TLS
	p := x509.NewCertPool()
	//TODO 同时加载rsa证书
	ca, err := ioutil.ReadFile(appCfg.CaPath)
	if err != nil {
		fmt.Println("read gm ca err: ", err)
		return
	}
	p.AppendCertsFromPEM(ca)
	sigCert, err := gmtls.LoadX509KeyPair(appCfg.SignCertPath, appCfg.SignKeyPath)
	if err != nil {
		fmt.Println(err)
	}
	encCert, err := gmtls.LoadX509KeyPair(appCfg.EncCertPath, appCfg.EncKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	cf := &gmtls.Config{
		GMSupport:    &gmtls.GMSupport{},
		Certificates: []gmtls.Certificate{sigCert, encCert},
		ClientAuth:   gmtls.RequireAndVerifyClientCert,
		ClientCAs:    p,
		MinVersion:   gmtls.VersionGMSSL,
		MaxVersion:   gmtls.VersionGMSSL,
		CipherSuites: []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3},
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	ln, err := gmtls.Listen("tcp", appCfg.LocalAddr, cf)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		log.Println("server tls socket close")
		ln.Close()
	}()

	// server -> client
	go toClient(appCfg, iface)
	// client -> server
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept client connect fail", err)
			continue
		}
		log.Println("receive client tls connect")
		//向客户端推送路由配置
		pushTunConfig(tunCfg, conn)
		//推送路由
		go pushRouteConfig(appCfg, tunCfg, conn)
		go toServer(appCfg, conn, iface)
	}
}

func pushTunConfig(serTunCfg tun.TunConfig, coon net.Conn) error {
	log.Printf("push client tun config")
	//TODO 为客户端分配虚拟ip
	cvip := "10.9.0.10"
	cidr := cvip + "/" + strconv.Itoa(serTunCfg.PrefixLen)
	clientTunCfg := tun.TunConfig{
		PrefixLen:  serTunCfg.PrefixLen,
		Cidr:       cidr,
		Mtu:        serTunCfg.Mtu,
		SVip:       serTunCfg.SVip,
		CVip:       cvip,
		GlobalMode: false,
	}
	cfg, err := json.Marshal(clientTunCfg)
	if err != nil {
		log.Println("...")
		return err
	}
	log.Printf("push ctl msg: %s", string(cfg))
	cfg = append(append([]byte(nil), RECORD_TYPE_CONTROL_TUN_CONFIG...), cfg...)
	//数据封包（添加类型和长度）
	cfg = Enpack(cfg, RECORD_TYPE_CONTROL)
	_, err = coon.Write(cfg)
	if err != nil {
		log.Println("write tun config to conn fail")
		return err
	}
	return nil
}

func pushRouteConfig(appCfg config.Config, serTunCfg tun.TunConfig, coon net.Conn) error {
	log.Printf("push route config")
	routeCfg := tun.RouteConfig{
		Ipv4Route: appCfg.Route,
		Ipv6Route: appCfg.Route6,
		SVip:      serTunCfg.SVip,
		SVip6:     serTunCfg.SVip6,
	}
	cfg, err := json.Marshal(routeCfg)
	if err != nil {
		log.Println("...")
		return err
	}
	log.Printf("push ctl msg: %s", string(cfg))
	cfg = append(append([]byte(nil), RECORD_TYPE_CONTROL_PUSH_ROUTE...), cfg...)
	//数据封包（添加类型和长度）
	cfg = Enpack(cfg, RECORD_TYPE_CONTROL)
	_, err = coon.Write(cfg)
	if err != nil {
		log.Println("write route config to conn fail")
		return err
	}
	return nil
}

// toClient sends packets from iface to tlsconn
func toClient(config config.Config, iface *water.Interface) {
	log.Println("listening tun data")
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		b := packet[:n]
		if key := netutil.GetDstKey(b); key != "" {
			if v, ok := cache.GetCache().Get(key); ok {
				if config.Compress {
					b = snappy.Encode(nil, b)
				}
				//数据封包（添加类型和长度）
				b = Enpack(b, RECORD_TYPE_DATA)
				//写入数据
				_, err := v.(net.Conn).Write(b)
				if err != nil {
					cache.GetCache().Delete(key)
					continue
				}
				counter.IncrWrittenBytes(n)
			}
		}
	}
}

// toServer sends packets from tlsconn to iface
func toServer(config config.Config, tlsconn net.Conn, iface *water.Interface) {
	defer tlsconn.Close()
	packet := make([]byte, config.BufferSize)
	//存放未读取完整数据
	tmpBuffer := make([]byte, 0)
	for {
		n, err := tlsconn.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		//解包
		var record []byte
		tmpBuffer, record = Depack(append(tmpBuffer, packet[:n]...))
		if bytes.Equal(record, nil) {
			continue
		}
		// 解析数据类型
		b := record[HEADER_LEN:]
		t := record[:RECORD_TYPE_LEN]
		if bytes.Equal(t, RECORD_TYPE_DATA) {
			// fmt.Println("data record")
			if config.Compress {
				b, err = snappy.Decode(nil, b)
				if err != nil {
					netutil.PrintErr(err, config.Verbose)
					break
				}
			}
			if key := netutil.GetSrcKey(b); key != "" {
				cache.GetCache().Set(key, tlsconn, 24*time.Hour)
				_, err := iface.Write(b)
				if err != nil {
					netutil.PrintErr(err, config.Verbose)
					break
				}
				counter.IncrReadBytes(len(b))
			}
		} else if bytes.Equal(t, RECORD_TYPE_CONTROL) {
			fmt.Println("control record")
		} else if bytes.Equal(t, RECORD_TYPE_AUTH) {
			fmt.Println("auth record")
		} else if bytes.Equal(t, RECORD_TYPE_ALARM) {
			fmt.Println("alarm record")
		} else {
			fmt.Println("unknown record")
		}
	}
}
