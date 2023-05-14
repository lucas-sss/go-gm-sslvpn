//go:build darwin
// +build darwin

/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-08 00:17:21
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-13 12:09:18
 * @FilePath: /gmvpn/vendor/github.com/net-byte/go-gateway/gateway_darwin.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

package gateway

import (
	"net"
)

func discoverGatewayOSSpecificIPv4() (ip net.IP, err error) {
	ipstr := execCmd("sh", "-c", "route -n get default | grep 'gateway' | awk 'NR==1{print $2}'")
	ipv4 := net.ParseIP(ipstr)
	if ipv4 == nil {
		return nil, errCantParse
	}
	return ipv4, nil
}

func discoverGatewayOSSpecificIPv6() (ip net.IP, err error) {
	ipstr := execCmd("sh", "-c", "route -6 -n get default | grep 'gateway' | awk 'NR==1{print $2}'")
	ipv6 := net.ParseIP(ipstr)
	if ipv6 == nil {
		return nil, errCantParse
	}
	return ipv6, nil
}
