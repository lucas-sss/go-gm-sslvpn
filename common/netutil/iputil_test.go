/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-13 16:33:51
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-14 23:16:20
 * @FilePath: /gmvpn/common/netutil/ip_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package netutil

import (
	"fmt"
	"net"
	"testing"
)

func TestIp(t *testing.T) {
	cidr := "10.8.8.0/24"
	ip, net, _ := net.ParseCIDR(cidr)
	fmt.Println("ip: ", ip.To4().String())

	ones, mask := net.Mask.Size()
	fmt.Println("mask: ", mask, ", onse: ", ones)

	begin, end := IpNetRange(net)
	fmt.Println("start ip: ", begin)
	fmt.Println("end ip: ", end)

	intBegin := IPv4ToInt(begin)
	intEnd := IPv4ToInt(end)
	fmt.Println("all host num: ", intEnd-intBegin)

}
