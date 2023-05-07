package tls

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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

// StartServer starts the tls server
func StartServer(iface *water.Interface, config config.Config) {
	log.Printf("vtun tls server started on %v", config.LocalAddr)
	// cert, err := tls.LoadX509KeyPair(config.TLSCertificateFilePath, config.TLSCertificateKeyFilePath)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// tlsConfig := &tls.Config{
	// 	Certificates:     []tls.Certificate{cert},
	// 	MinVersion:       tls.VersionTLS13,
	// 	CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
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
	// ln, err := tls.Listen("tcp", config.LocalAddr, tlsConfig)
	// if err != nil {
	// 	log.Panic(err)
	// }

	p := x509.NewCertPool()
	//TODO 同时加载rsa证书

	ca, err := ioutil.ReadFile(config.CaPath)
	if err != nil {
		fmt.Println("read gm ca err: ", err)
		return
	}
	p.AppendCertsFromPEM(ca)

	sigCert, err := gmtls.LoadX509KeyPair(config.SignCertPath, config.SignKeyPath)
	if err != nil {
		fmt.Println(err)
	}

	encCert, err := gmtls.LoadX509KeyPair(config.EncCertPath, config.EncKeyPath)
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
	ln, err := gmtls.Listen("tcp", config.LocalAddr, cf)
	if err != nil {
		fmt.Println(err)
	}
	defer ln.Close()

	// server -> client
	go toClient(config, iface)
	// client -> server
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		sniffConn := NewPeekPreDataConn(conn)
		switch sniffConn.Type {
		case TypeHttp:
			if sniffConn.Handle() {
				continue
			}
		case TypeHttp2:
			if sniffConn.Handle() {
				continue
			}
		}
		go toServer(config, sniffConn, iface)
	}
}

// toClient sends packets from iface to tlsconn
func toClient(config config.Config, iface *water.Interface) {
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
		// if key := netutil.GetSrcKey(b); key != "" {
		// 	cache.GetCache().Set(key, tlsconn, 24*time.Hour)
		// 	iface.Write(b)
		// 	counter.IncrReadBytes(n)
		// }

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
