package tls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	"gmvpn/common"
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

type ClientConn struct {
	Vip string
}

var clientConnCache *common.RWMutexMap

// StartServer starts the tls server
func StartServer(iface *water.Interface, appCfg config.Config, tunCfg tun.TunConfig) {
	defer func() {
		log.Println("server tun interface close")
		iface.Close()
	}()
	log.Printf("vtun tls server started on %v", appCfg.LocalAddr)

	clientConnCache = common.NewRWMutexMap(1024)

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
		err = pushTunConfig(appCfg, tunCfg, conn)
		if err != nil {
			conn.Close()
			continue
		}
		//推送路由
		go pushRouteConfig(appCfg, tunCfg, conn)
		go toServer(appCfg, conn, iface)
	}
}

func pushTunConfig(appCfg config.Config, serTunCfg tun.TunConfig, coon net.Conn) error {
	log.Printf("push tun config")
	// 为客户端分配虚拟ip
	cvip, err := allocateVip(coon, appCfg)
	if err != nil {
		//TODO 推送vip不足消息
		coon.Close()
		return err
	}
	cidr := cvip + "/" + strconv.Itoa(serTunCfg.Mask)
	clientTunCfg := tun.TunConfig{
		Mtu:        serTunCfg.Mtu,
		SVip:       serTunCfg.SVip,
		CVip:       cvip,
		Mask:       serTunCfg.Mask,
		Cidr:       cidr,
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

//分配虚拟ip
func allocateVip(coon net.Conn, appCfg config.Config) (string, error) {
	id := coon.RemoteAddr().String()
	vipPool := appCfg.VipPool
	vipList := appCfg.VipList
	for _, vip := range vipList {
		v, suc := vipPool.Get(vip)
		if !suc {
			continue
		}
		vipInfo := v.(config.VipInfo)
		if vipInfo.Used {
			continue
		}
		ok := vipPool.TrySet(vip, config.VipInfo{
			Used: true,
			Id:   id,
		})
		if !ok {
			continue
		}
		log.Printf("conn[%s] allocate vip: %s", id, vip)
		clientConn := ClientConn{Vip: vip}
		//缓存客户端链接信息
		clientConnCache.Set(id, clientConn)
		return vip, nil
	}
	return "", fmt.Errorf("insufficient virtual ip")
}

//推送路由陪孩子
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
			log.Println("read err, ", err)
			if err == io.EOF {
				//客户端关闭，释放vip
				id := tlsconn.RemoteAddr().String()
				if v, ok := clientConnCache.Get(id); ok {
					clientConn := v.(ClientConn)
					releaseVip(config, clientConn.Vip)
				}
			}
			break
		}
		//解包
		var records [][]byte
		tmpBuffer, records = Depack(append(tmpBuffer, packet[:n]...))
		if len(records) == 0 {
			continue
		}
		for i := 0; i < len(records); i++ {
			record := records[i]
			// 解析数据类型
			b := record[HEADER_LEN:]
			t := record[:RECORD_TYPE_LEN]
			if bytes.Equal(t, RECORD_TYPE_DATA) {
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
						log.Println("server write to tun iface fail, ", err)
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
}

func releaseVip(appCfg config.Config, vip string) {
	log.Println("release vip: ", vip)
	vipInfo := config.VipInfo{
		Used: false,
		Id:   "",
	}
	appCfg.VipPool.Set(vip, vipInfo)
}
