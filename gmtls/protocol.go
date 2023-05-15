/*
 * @Author: liuwei lyy9645@163.com
 * @Date: 2023-05-06 23:15:05
 * @LastEditors: liuwei lyy9645@163.com
 * @LastEditTime: 2023-05-16 00:16:46
 * @FilePath: /vtun/tls/protocol.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tls

import (
	"bytes"
	"encoding/binary"
)

const (
	RECORD_TYPE_LEN = 2
	RECORD_DATA_LEN = 4
	HEADER_LEN      = RECORD_TYPE_LEN + RECORD_DATA_LEN
)

var (
	RECORD_TYPE_DATA               = []byte{0x11, 0x11}
	RECORD_TYPE_CONTROL            = []byte{0x12, 0x11}
	RECORD_TYPE_CONTROL_TUN_CONFIG = []byte{0x12, 0x12}
	RECORD_TYPE_CONTROL_PUSH_ROUTE = []byte{0x12, 0x13}
	RECORD_TYPE_ALARM              = []byte{0x13, 0x11}
	RECORD_TYPE_AUTH               = []byte{0x14, 0x11}
)

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

//封包
func Enpack(body, bodyType []byte) []byte {
	return append(append(append([]byte(nil), bodyType...), IntToBytes(len(body))...), body...)
}

//解包
func Depack(buffer []byte) ([]byte, [][]byte) {
	length := len(buffer)
	records := make([][]byte, 0)

	var i int
	for i = 0; i < length; i++ {
		if length < i+HEADER_LEN {
			break
		}
		t := buffer[i : i+RECORD_TYPE_LEN]
		if bytes.Equal(RECORD_TYPE_DATA, t) || bytes.Equal(RECORD_TYPE_CONTROL, t) || bytes.Equal(RECORD_TYPE_ALARM, t) ||
			bytes.Equal(RECORD_TYPE_AUTH, t) {
			byteLen := buffer[i+RECORD_TYPE_LEN : i+HEADER_LEN]
			l := BytesToInt(byteLen)
			if length < i+HEADER_LEN+l {
				break
			}
			record := append([]byte(nil), buffer[i:i+HEADER_LEN+l]...)
			records = append(records, record)
			i += HEADER_LEN + l
		}
	}
	i--
	if i == length {
		return make([]byte, 0), records
	}
	if i < 0 {
		i = 0
	}
	return buffer[i:], records
}
