package tls

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gmvpn/common/cache"
	"gmvpn/common/config"
	"gmvpn/common/counter"
	"gmvpn/common/netutil"

	"github.com/golang/snappy"
	"github.com/net-byte/water"
	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
)

// StartClient starts the tls client
func StartClient(iface *water.Interface, config config.Config) {
	log.Println("vtun tls client started")
	go tunToTLS(config, iface)
	// tlsconfig := &tls.Config{
	// 	InsecureSkipVerify: config.TLSInsecureSkipVerify,
	// 	MinVersion:         tls.VersionTLS13,
	// 	CurvePreferences:   []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	// 		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	// 		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	// 	},
	// }
	// if config.TLSSni != "" {
	// 	tlsconfig.ServerName = config.TLSSni
	// }
	// for {
	// 	conn, err := tls.Dial("tcp", config.ServerAddr, tlsconfig)
	// 	if err != nil {
	// 		time.Sleep(3 * time.Second)
	// 		netutil.PrintErr(err, config.Verbose)
	// 		continue
	// 	}
	// 	cache.GetCache().Set("tlsconn", conn, 24*time.Hour)
	// 	tlsToTun(config, conn, iface)
	// 	cache.GetCache().Delete("tlsconn")
	// }

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
		MinVersion:         gmtls.VersionGMSSL,
		MaxVersion:         gmtls.VersionGMSSL,
		CipherSuites:       []uint16{gmtls.GMTLS_SM2_WITH_SM4_SM3},
	}

	// conn, err := gmtls.Dial("tcp", config.ServerAddr, cf)
	// if err != nil {
	// 	fmt.Printf("%s\r\n", err)
	// 	return
	// }
	// defer conn.Close()

	if config.TLSSni != "" {
		cf.ServerName = config.TLSSni
	}
	for {
		conn, err := gmtls.Dial("tcp", config.RemoteAddr, cf)
		if err != nil {
			time.Sleep(3 * time.Second)
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		cache.GetCache().Set("tlsconn", conn, 24*time.Hour)
		tlsToTun(config, conn, iface)
		cache.GetCache().Delete("tlsconn")
	}
}

// tunToTLS sends packets from tun to tls
func tunToTLS(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if v, ok := cache.GetCache().Get("tlsconn"); ok {
			b := packet[:n]
			if config.Compress {
				b = snappy.Encode(nil, b)
			}
			//数据封包（添加类型和长度）
			b = Enpack(b, RECORD_TYPE_DATA)

			conn := v.(*gmtls.Conn)
			_, err = conn.Write(b)
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

// tlsToTun sends packets from tls to tun
func tlsToTun(config config.Config, tlsconn *gmtls.Conn, iface *water.Interface) {
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

		// b := packet[:n]
		// if config.Compress {
		// 	b, err = snappy.Decode(nil, b)
		// 	if err != nil {
		// 		netutil.PrintErr(err, config.Verbose)
		// 		break
		// 	}
		// }
		// if config.Obfs {
		// 	b = cipher.XOR(b)
		// }
		// _, err = iface.Write(b)
		// if err != nil {
		// 	netutil.PrintErr(err, config.Verbose)
		// 	break
		// }
		// counter.IncrReadBytes(n)

		//解包
		var record []byte
		tmpBuffer, record = Depack(append(tmpBuffer, packet[:n]...))
		if bytes.Equal(record, []byte{}) {
			continue
		}

		//解析数据类型
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
			_, err = iface.Write(b)
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				break
			}
			counter.IncrReadBytes(n)
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
