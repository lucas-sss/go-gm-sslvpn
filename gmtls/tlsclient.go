package tls

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gmvpn/common/cache"
	"gmvpn/common/config"
	"gmvpn/common/counter"
	"gmvpn/tun"

	"github.com/golang/snappy"
	"github.com/net-byte/water"
	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
)

var gloablTunCfg tun.TunConfig

// StartClient starts the tls client
func StartClient(config config.Config) {
	log.Println("tls client started")

	// 信任的根证书
	certPool := x509.NewCertPool()
	cacert, err := ioutil.ReadFile(config.CaPath)
	if err != nil {
		log.Fatal(err)
	}
	ok := certPool.AppendCertsFromPEM(cacert)
	if !ok {
		return
	}
	signCert, err := gmtls.LoadX509KeyPair(config.SignCertPath, config.SignKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	encCert, err := gmtls.LoadX509KeyPair(config.EncCertPath, config.EncKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	cf := &gmtls.Config{
		GMSupport:          &gmtls.GMSupport{},
		RootCAs:            certPool,
		Certificates:       []gmtls.Certificate{signCert, encCert},
		InsecureSkipVerify: config.TLSInsecureSkipVerify,
		MinVersion:         tls.VersionTLS11,
		MaxVersion:         gmtls.VersionGMSSL,
		CipherSuites:       []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3},
	}

	if config.TLSSni != "" {
		cf.ServerName = config.TLSSni
	}
	for {
		conn, err := gmtls.Dial("tcp", config.RemoteAddr, cf)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("connect to server success")
		cache.GetCache().Set("tlsconn", conn, 24*time.Hour)
		tlsToTun(config, conn)
		cache.GetCache().Delete("tlsconn")
	}
}

// tunToTLS sends packets from tun to tls
func tunToTLS(appCfg config.Config, iface *water.Interface) {
	defer iface.Close()

	packet := make([]byte, appCfg.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			log.Println("tun iface read err, ", err)
			break
		}
		if v, ok := cache.GetCache().Get("tlsconn"); ok {
			b := packet[:n]
			if appCfg.Compress {
				b = snappy.Encode(nil, b)
			}
			//数据封包（添加类型和长度）
			b = Enpack(b, RECORD_TYPE_DATA)

			conn := v.(*gmtls.Conn)
			_, err = conn.Write(b)
			if err != nil {
				log.Println("tls coon write err, ", err)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

// tlsToTun sends packets from tls to tun
func tlsToTun(appCfg config.Config, tlsconn *gmtls.Conn) {
	var iface *water.Interface
	defer func() {
		log.Println("close tls connection and tun interface")
		tlsconn.Close()
		if iface != nil {
			iface.Close()
		}
	}()
	packet := make([]byte, appCfg.BufferSize)
	//存放未读取完整数据
	tmpBuffer := make([]byte, 0)

	for {
		n, err := tlsconn.Read(packet)
		if err != nil {
			log.Println("tls conn read err, ", err)
			break
		}
		//解包处理
		var records [][]byte
		tmpBuffer, records = Depack(append(tmpBuffer, packet[:n]...))
		if len(records) == 0 {
			continue
		}
		for i := 0; i < len(records); i++ {
			record := records[i]
			//解析数据类型
			b := record[HEADER_LEN:]
			t := record[:RECORD_TYPE_LEN]
			if bytes.Equal(t, RECORD_TYPE_DATA) {
				if iface == nil {
					log.Println("iface is nil")
					continue
				}
				if appCfg.Compress {
					b, err = snappy.Decode(nil, b)
					if err != nil {
						log.Println("snappy decode read err, ", err)
						break
					}
				}
				// netutil.PrintEthernetFrame(b)
				_, err = iface.Write(b)
				if err != nil {
					log.Println("tun iface write err, ", err)
					break
				}
				counter.IncrReadBytes(n)
			} else if bytes.Equal(t, RECORD_TYPE_CONTROL) {
				fmt.Println("control record type", b[:2])
				ifaceNew := processCtlMsg(appCfg, b)
				if ifaceNew != nil {
					iface = ifaceNew
				}
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

func processCtlMsg(appCfg config.Config, ctlMsg []byte) *water.Interface {
	ctlType := ctlMsg[:2]
	data := ctlMsg[2:]
	log.Printf("control record: %s", string(data))
	if bytes.Equal(ctlType, RECORD_TYPE_CONTROL_TUN_CONFIG) {
		tunCfg := &tun.TunConfig{}
		err := json.Unmarshal(data, tunCfg)
		if err != nil {
			log.Panic("server push tun config data error")
		}
		tunCfg.Device = appCfg.Device
		tunCfg.ServerMode = false
		tunCfg.LocalGateway = appCfg.LocalGateway
		tunCfg.LocalGateway6 = appCfg.LocalGateway6
		//创建tun网卡设备
		iface := tun.CreatAndConfigTun(tunCfg)
		//读取虚拟网卡数据
		go tunToTLS(appCfg, iface)
		gloablTunCfg = *tunCfg
		return iface
	} else if bytes.Equal(ctlType, RECORD_TYPE_CONTROL_PUSH_ROUTE) {
		routeCfg := &tun.RouteConfig{}
		err := json.Unmarshal(data, routeCfg)
		if err != nil {
			log.Panic("server push route config error")
		}
		routeCfg.Device = gloablTunCfg.Device
		routeCfg.SVip = gloablTunCfg.SVip
		routeCfg.SVip6 = gloablTunCfg.SVip6
		go tun.ConfigRoute(*routeCfg)
	} else {
		log.Println("unknown control msg type: ", ctlType)
	}
	return nil
}
