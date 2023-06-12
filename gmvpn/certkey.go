/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-21 15:00:06
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-21 15:01:43
 * @FilePath: /gmvpn/gmvpn/certkey.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import "fyne.io/fyne/v2"

var resourceCaCrt = &fyne.StaticResource{
	StaticName: "ca.crt",
	StaticContent: []byte(
		"-----BEGIN CERTIFICATE-----\nMIICWTCCAf6gAwIBAgIJAKUAeq6Z86Q/MAoGCCqBHM9VAYN1MIGGMQswCQYDVQQG\nEwJDTjELMAkGA1UECAwCWkoxCzAJBgNVBAcMAkhaMQwwCgYDVQQKDANGTEsxEzAR\nBgNVBAsMCkZMS19zc2x2cG4xDzANBgNVBAMMBlJvb3RDQTEpMCcGCSqGSIb3DQEJ\nARYaZnVsYW5rZXNlcnZpY2VzQGZsa3NlYy5jb20wHhcNMjEwODE3MDYzODE3WhcN\nMzEwODE1MDYzODE3WjCBhjELMAkGA1UEBhMCQ04xCzAJBgNVBAgMAlpKMQswCQYD\nVQQHDAJIWjEMMAoGA1UECgwDRkxLMRMwEQYDVQQLDApGTEtfc3NsdnBuMQ8wDQYD\nVQQDDAZSb290Q0ExKTAnBgkqhkiG9w0BCQEWGmZ1bGFua2VzZXJ2aWNlc0BmbGtz\nZWMuY29tMFkwEwYHKoZIzj0CAQYIKoEcz1UBgi0DQgAEvjoGzYxSyb+GoRVbcfwG\nZWv1iAGaeLuOOF0vD8EOoTCDA6Cvgw1W1Cp53as7D6q6xZdEgxgM62+0i/5hPCMU\nGqNTMFEwHQYDVR0OBBYEFPRQbrJI3hFu2O0jSblR+SVWBHH/MB8GA1UdIwQYMBaA\nFPRQbrJI3hFu2O0jSblR+SVWBHH/MA8GA1UdEwEB/wQFMAMBAf8wCgYIKoEcz1UB\ng3UDSQAwRgIhAIDHLYzq8IxkGpwHtaxHOBQcw3EOB7tSpC1VQWM3/OAbAiEAy46a\nJLxbU+G9h0TXY2xQ4Q0J35qr0VZOBiksOlgnc8A=\n-----END CERTIFICATE-----\n"),
}

var resourceEnccertCrt = &fyne.StaticResource{
	StaticName: "enccert.crt",
	StaticContent: []byte(
		"Certificate:\n    Data:\n        Version: 3 (0x2)\n        Serial Number: 1 (0x1)\n    Signature Algorithm: sm2sign-with-sm3\n        Issuer: C=CN, ST=ZJ, L=HZ, O=FLK, OU=FLK_sslvpn, CN=RootCA/emailAddress=fulankeservices@flksec.com\n        Validity\n            Not Before: Aug 17 06:38:17 2021 GMT\n            Not After : Aug 15 06:38:17 2031 GMT\n        Subject: C=CN, ST=ZJ, O=FLK, OU=FLK_sslvpn, CN=ServerEnc/emailAddress=sslvpn@flksec.com\n        Subject Public Key Info:\n            Public Key Algorithm: id-ecPublicKey\n                Public-Key: (256 bit)\n                pub:\n                    04:81:3f:a3:33:32:5f:65:85:c0:e1:2a:3c:8b:19:\n                    30:18:87:ab:0c:7e:31:e7:eb:b6:b9:62:85:80:de:\n                    0c:55:60:d2:5c:28:09:55:8d:5e:be:b5:90:2b:99:\n                    d8:50:ae:6c:56:e4:51:2f:8a:cf:58:cb:53:2b:96:\n                    bf:76:79:67:c8\n                ASN1 OID: sm2p256v1\n                NIST CURVE: SM2\n        X509v3 extensions:\n            X509v3 Basic Constraints: \n                CA:FALSE\n            X509v3 Key Usage: \n                Key Encipherment\n            Netscape Comment: \n                GmSSL Generated Certificate\n            X509v3 Subject Key Identifier: \n                3F:37:67:95:AA:EE:8A:97:56:E6:6F:38:C5:09:2C:61:AF:C5:FA:58\n            X509v3 Authority Key Identifier: \n                keyid:F4:50:6E:B2:48:DE:11:6E:D8:ED:23:49:B9:51:F9:25:56:04:71:FF\n\n    Signature Algorithm: sm2sign-with-sm3\n         30:45:02:20:4b:82:10:d2:f5:65:04:25:74:46:04:3f:5e:5e:\n         42:31:97:74:6a:5a:ec:b5:43:59:20:14:43:51:9e:c7:db:22:\n         02:21:00:b8:d5:71:c7:bb:97:1e:ec:4a:f3:ba:31:ef:1d:45:\n         44:d9:a3:81:7b:88:01:1b:f0:da:4c:c6:fe:e7:29:b4:68\n-----BEGIN CERTIFICATE-----\nMIICcTCCAhegAwIBAgIBATAKBggqgRzPVQGDdTCBhjELMAkGA1UEBhMCQ04xCzAJ\nBgNVBAgMAlpKMQswCQYDVQQHDAJIWjEMMAoGA1UECgwDRkxLMRMwEQYDVQQLDApG\nTEtfc3NsdnBuMQ8wDQYDVQQDDAZSb290Q0ExKTAnBgkqhkiG9w0BCQEWGmZ1bGFu\na2VzZXJ2aWNlc0BmbGtzZWMuY29tMB4XDTIxMDgxNzA2MzgxN1oXDTMxMDgxNTA2\nMzgxN1owczELMAkGA1UEBhMCQ04xCzAJBgNVBAgMAlpKMQwwCgYDVQQKDANGTEsx\nEzARBgNVBAsMCkZMS19zc2x2cG4xEjAQBgNVBAMMCVNlcnZlckVuYzEgMB4GCSqG\nSIb3DQEJARYRc3NsdnBuQGZsa3NlYy5jb20wWTATBgcqhkjOPQIBBggqgRzPVQGC\nLQNCAASBP6MzMl9lhcDhKjyLGTAYh6sMfjHn67a5YoWA3gxVYNJcKAlVjV6+tZAr\nmdhQrmxW5FEvis9Yy1Mrlr92eWfIo4GHMIGEMAkGA1UdEwQCMAAwCwYDVR0PBAQD\nAgUgMCoGCWCGSAGG+EIBDQQdFhtHbVNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUw\nHQYDVR0OBBYEFD83Z5Wq7oqXVuZvOMUJLGGvxfpYMB8GA1UdIwQYMBaAFPRQbrJI\n3hFu2O0jSblR+SVWBHH/MAoGCCqBHM9VAYN1A0gAMEUCIEuCENL1ZQQldEYEP15e\nQjGXdGpa7LVDWSAUQ1Gex9siAiEAuNVxx7uXHuxK87ox7x1FRNmjgXuIARvw2kzG\n/ucptGg=\n-----END CERTIFICATE-----\n"),
}

var resourceEnckeyKey = &fyne.StaticResource{
	StaticName: "enckey.key",
	StaticContent: []byte(
		"-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBG0wawIBAQQg4P5cRT6AImnAiRT2\nsKsmiowlbIcVIzu6qAZ4b0wJXQahRANCAASBP6MzMl9lhcDhKjyLGTAYh6sMfjHn\n67a5YoWA3gxVYNJcKAlVjV6+tZArmdhQrmxW5FEvis9Yy1Mrlr92eWfI\n-----END PRIVATE KEY-----\n"),
}

var resourceSigncertCrt = &fyne.StaticResource{
	StaticName: "signcert.crt",
	StaticContent: []byte(
		"Certificate:\n    Data:\n        Version: 3 (0x2)\n        Serial Number: 2 (0x2)\n    Signature Algorithm: sm2sign-with-sm3\n        Issuer: C=CN, ST=ZJ, L=HZ, O=FLK, OU=FLK_sslvpn, CN=RootCA/emailAddress=fulankeservices@flksec.com\n        Validity\n            Not Before: Aug 17 06:38:17 2021 GMT\n            Not After : Aug 15 06:38:17 2031 GMT\n        Subject: C=CN, ST=ZJ, O=FLK, OU=FLK_sslvpn, CN=ServerSig/emailAddress=sslvpn@flksec.com\n        Subject Public Key Info:\n            Public Key Algorithm: id-ecPublicKey\n                Public-Key: (256 bit)\n                pub:\n                    04:7e:40:6d:58:c0:f7:e7:47:67:e6:a8:7a:f2:5a:\n                    93:69:55:6e:43:f7:fd:75:16:1c:f7:73:02:81:93:\n                    73:c4:df:bc:a6:f0:b1:d0:b3:4f:58:8a:4c:84:9d:\n                    90:11:08:f2:2a:f0:e3:2a:7b:23:bd:54:73:55:9a:\n                    0d:d5:22:ce:8e\n                ASN1 OID: sm2p256v1\n                NIST CURVE: SM2\n        X509v3 extensions:\n            X509v3 Basic Constraints: \n                CA:FALSE\n            X509v3 Key Usage: \n                Digital Signature\n            Netscape Comment: \n                GmSSL Generated Certificate\n            X509v3 Subject Key Identifier: \n                C5:AA:B1:F2:0E:0F:4D:F7:6F:60:CB:3F:9D:F6:C2:60:61:E9:49:D2\n            X509v3 Authority Key Identifier: \n                keyid:F4:50:6E:B2:48:DE:11:6E:D8:ED:23:49:B9:51:F9:25:56:04:71:FF\n\n            X509v3 Extended Key Usage: \n                TLS Web Server Authentication, TLS Web Client Authentication\n    Signature Algorithm: sm2sign-with-sm3\n         30:44:02:20:2c:e5:73:2e:37:37:5f:d4:e7:fb:70:a6:07:d1:\n         bd:2f:3a:dc:71:2c:a8:6b:49:13:9e:c4:6f:43:ef:61:76:47:\n         02:20:4a:e6:02:34:87:c2:ff:3f:54:49:ca:bc:dd:96:26:12:\n         24:e5:8f:f3:20:47:47:de:56:37:71:f2:89:3d:09:7e\n-----BEGIN CERTIFICATE-----\nMIICjzCCAjagAwIBAgIBAjAKBggqgRzPVQGDdTCBhjELMAkGA1UEBhMCQ04xCzAJ\nBgNVBAgMAlpKMQswCQYDVQQHDAJIWjEMMAoGA1UECgwDRkxLMRMwEQYDVQQLDApG\nTEtfc3NsdnBuMQ8wDQYDVQQDDAZSb290Q0ExKTAnBgkqhkiG9w0BCQEWGmZ1bGFu\na2VzZXJ2aWNlc0BmbGtzZWMuY29tMB4XDTIxMDgxNzA2MzgxN1oXDTMxMDgxNTA2\nMzgxN1owczELMAkGA1UEBhMCQ04xCzAJBgNVBAgMAlpKMQwwCgYDVQQKDANGTEsx\nEzARBgNVBAsMCkZMS19zc2x2cG4xEjAQBgNVBAMMCVNlcnZlclNpZzEgMB4GCSqG\nSIb3DQEJARYRc3NsdnBuQGZsa3NlYy5jb20wWTATBgcqhkjOPQIBBggqgRzPVQGC\nLQNCAAR+QG1YwPfnR2fmqHryWpNpVW5D9/11Fhz3cwKBk3PE37ym8LHQs09YikyE\nnZARCPIq8OMqeyO9VHNVmg3VIs6Oo4GmMIGjMAkGA1UdEwQCMAAwCwYDVR0PBAQD\nAgeAMCoGCWCGSAGG+EIBDQQdFhtHbVNTTCBHZW5lcmF0ZWQgQ2VydGlmaWNhdGUw\nHQYDVR0OBBYEFMWqsfIOD033b2DLP532wmBh6UnSMB8GA1UdIwQYMBaAFPRQbrJI\n3hFu2O0jSblR+SVWBHH/MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAK\nBggqgRzPVQGDdQNHADBEAiAs5XMuNzdf1Of7cKYH0b0vOtxxLKhrSROexG9D72F2\nRwIgSuYCNIfC/z9UScq83ZYmEiTlj/MgR0feVjdx8ok9CX4=\n-----END CERTIFICATE-----\n"),
}

var resourceSignkeyKey = &fyne.StaticResource{
	StaticName: "signkey.key",
	StaticContent: []byte(
		"-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBG0wawIBAQQgWvzCr7SFBTg0MklS\n+JTrdN0g2gNFoZ/TECjObGINTBGhRANCAAR+QG1YwPfnR2fmqHryWpNpVW5D9/11\nFhz3cwKBk3PE37ym8LHQs09YikyEnZARCPIq8OMqeyO9VHNVmg3VIs6O\n-----END PRIVATE KEY-----\n"),
}
