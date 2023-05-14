/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-13 16:33:04
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-13 18:08:09
 * @FilePath: /gmvpn/common/netutil/ip.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package netutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func IntToIPv4(i uint32) (string, error) {
	ipv4 := uint32(i)

	ip0 := ipv4 >> 24
	ip8 := ipv4 << 8 >> 24
	ip16 := ipv4 << 16 >> 24
	ip24 := ipv4 << 24 >> 24
	if ip0 > 255 || ip8 > 255 || ip16 > 255 || ip24 > 255 {
		return "", errors.New(fmt.Sprintf("bad ipv4 int"))
	}
	return fmt.Sprintf("%d.%d.%d.%d", ipv4>>24, ipv4<<8>>24, ipv4<<16>>24, ipv4<<24>>24), nil
}

// IPv4ByLong
// 将 uint32 长整型转换成IPV4 地址
// converts a uint32 represented by a string into an ipv4 address string
// 168427779 => "10.10.1.2"
func IntStrToIPv4(str string) (string, error) {
	ipv4Int, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return "", errors.New(fmt.Sprintf("fail to convert string to Int32 :%s", err.Error()))
	}
	return IntToIPv4(uint32(ipv4Int))
}

// IPv6ByLong
// 将 big.Int 长整型转换成IPV6 地址
// converts a big integer represented by a string into an IPv6 address string
// 53174336768213711679990085974688268287=> "2801:0137:0000:0000:0000:ffff:ffff:ffff"
func IPv6ByLong(ipv6long string) (string, error) {
	bi, ok := new(big.Int).SetString(ipv6long, 10)
	if !ok {
		return "", errors.New("fail to convert string to big.Int")
	}

	b255 := new(big.Int).SetBytes([]byte{255})
	var buf = make([]byte, 2)
	p := make([]string, 8)
	j := 0
	var i uint
	tmpint := new(big.Int)
	for i = 0; i < 16; i += 2 {
		tmpint.Rsh(bi, 120-i*8).And(tmpint, b255)
		bytes := tmpint.Bytes()
		if len(bytes) > 0 {
			buf[0] = bytes[0]
		} else {
			buf[0] = 0
		}
		tmpint.Rsh(bi, 120-(i+1)*8).And(tmpint, b255)
		bytes = tmpint.Bytes()
		if len(bytes) > 0 {
			buf[1] = bytes[0]
		} else {
			buf[1] = 0
		}
		p[j] = hex.EncodeToString(buf)
		j++
	}

	return strings.Join(p, ":"), nil
}

// 将IPV4 转换成 uint32 长整型
func IPv4ToInt(ipv4 string) (ip uint32) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return 0
	}
	ips := reg.FindStringSubmatch(ipv4)
	if ips == nil {
		return 0
	}

	//上面正则做了判断，这里就不报错了
	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return 0
	}

	ip += uint32(ip1 * 0x1000000) // 左移24位
	ip += uint32(ip2 * 0x10000)   // 左移16位
	ip += uint32(ip3 * 0x100)     // 左移8位
	ip += uint32(ip4)             // 左移0位

	return ip
}

// 将IPV6 转换成 big.Int 长整型
func IPv6ToInt(ipv6 string) (*big.Int, error) {
	ip := net.ParseIP(ipv6)
	return NetIpv6ToInt(ip)
}

// 将net.IP 类型 转换成 big.Int 长整型
func NetIpv6ToInt(ip net.IP) (*big.Int, error) {
	if ip == nil {
		return nil, errors.New("invalid ipv6")
	}
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ip)
	return IPv6Int, nil
}

func IsIPv4(ip string) bool {
	if net.ParseIP(ip).To4() != nil {
		return true
	}
	return false
}

// IpNetRange 返回网段的起始IP、结束IP
func IpNetRange(ipNet *net.IPNet) (start, end string) {
	mask := ipNet.Mask
	broadcast := copyIp(ipNet.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(broadcast) - i - 1
		broadcast[ipIdx] = ipNet.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	return ipNet.IP.String(), broadcast.String()
}

// IncreaseIP IP地址自增
func IncreaseIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

// DecreaseIP IP地址自减
func DecreaseIP(ip net.IP) {
	length := len(ip)
	for i := length - 1; i >= 0; i-- {
		ip[length-1]--
		if ip[length-1] < 0xFF {
			break
		}
		for j := 1; j < length; j++ {
			ip[length-j-1]--
			if ip[length-j-1] < 0xFF {
				return
			}
		}
	}
}

// Format ipv6最简格式，示例：240e:f7:c000:103:13::f4
func Format(ip string) string {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return ip
	}
	return netIp.String()
}

// FormatZero ipv6省略前导零格式，示例：240e:f7:c000:103:13:0:0:f4
func FormatZero(ip string) string {
	p := net.ParseIP(ip)
	if p == nil || p.To4() != nil || len(p) != net.IPv6len {
		return ip
	}

	const maxLen = len("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
	b := make([]byte, 0, maxLen)

	for i := 0; i < net.IPv6len; i += 2 {
		if i > 0 {
			b = append(b, ':')
		}
		b = appendHex(b, (uint32(p[i])<<8)|uint32(p[i+1]))
	}
	return string(b)
}

const hexDigit = "0123456789abcdef"

// Convert i to a hexadecimal string. Leading zeros are not printed.
func appendHex(dst []byte, i uint32) []byte {
	if i == 0 {
		return append(dst, '0')
	}
	for j := 7; j >= 0; j-- {
		v := i >> uint(j*4)
		if v > 0 {
			dst = append(dst, hexDigit[v&0xf])
		}
	}
	return dst
}

func copyIp(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
