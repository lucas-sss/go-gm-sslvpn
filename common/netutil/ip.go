package netutil

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func PrintEthernetFrame(packet []byte) {
	size := len(packet)
	//获取以太网帧内部数据类型
	// 0x0800代表IP协议帧
	// 0x0806代表ARP协议帧
	// 0x8864代表PPPoE
	// 0x86dd代表IPv6
	fmt.Println("以太网帧: ", hex.EncodeToString(packet))
	mt := fmt.Sprintf("%x", MACType(packet))

	fmt.Println("帧类型：", mt)
	//动态获取ipv4头部长度
	printMACHeader(packet)

	//16进制0800代表ip报文
	if strings.EqualFold(mt, "0800") {
		b := packet[14:]
		printIPv4Header(b)
		hl := getIPv4HeaderLen(b)
		if b[9] == 1 { //icmp
			icmpPacket := b[hl:]
			printICMPHeader(icmpPacket)
		}

		if b[9] == 6 { //tcp
			tcpPacket := b[hl:]
			printTCPHeader(tcpPacket)
		}

		if b[9] == 17 { //udp
			udpPacket := b[hl:]
			printUDPHeader(udpPacket)
		}
	}

	//arp报文
	if strings.EqualFold(mt, "0806") {
		b := packet[14:size]
		printARPHeader(b)
	}

	srcIP := getSrcIP(packet)
	dstIP := getDstIP(packet)
	fmt.Printf("Msg Protocol type: %v(1=ICMP, 6=TCP, 17=UDP)\tsrcIP:%v ---> dstIP:%v\n", packet[9], srcIP, dstIP)
}

func PrintEthernetFrameData(packet []byte) {
	b := packet[14:]
	printIPv4Header(b)
	hl := getIPv4HeaderLen(b)
	if b[9] == 1 { //icmp
		icmpPacket := b[hl:]
		printICMPHeader(icmpPacket)
	}

	if b[9] == 6 { //tcp
		tcpPacket := b[hl:]
		printTCPHeader(tcpPacket)
	}

	if b[9] == 17 { //udp
		udpPacket := b[hl:]
		printUDPHeader(udpPacket)
	}

	srcIP := getSrcIP(packet)
	dstIP := getDstIP(packet)
	fmt.Printf("Msg Protocol type: %v(1=ICMP, 6=TCP, 17=UDP)\tsrcIP:%v ---> dstIP:%v\n", packet[9], srcIP, dstIP)
}

// 读取以太网帧的第13、14个字节，这两个字节代表数据类型
func MACType(macFrame []byte) []byte {
	return macFrame[12:14]
}

func MACTypeARP(macFrame []byte) []byte {
	return macFrame[2:4]
}

func printMACHeader(packet []byte) {
	fmt.Println("----->MAC Header<----")
	printSrcMACEther(packet)
	printDstMACEther(packet)
	printTypeEther(packet)
	fmt.Println("----->MAC Header<----")
}

func printDstMACEther(packet []byte) {
	fmt.Printf("Ether Header--->DstMac:%x:%x:%x:%x:%x:%x\n", packet[0], packet[1], packet[2], packet[3], packet[4], packet[5])
}

func printSrcMACEther(packet []byte) {
	fmt.Printf("Ether Header--->SrcMac:%x:%x:%x:%x:%x:%x\n", packet[6], packet[7], packet[8], packet[9], packet[10], packet[11])
}
func printTypeEther(packet []byte) {
	fmt.Printf("Ether Header--->Type:%04x\n", uint16(packet[12])<<8|uint16(packet[13]))
}

func MACDestination(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[:6])
}

func MACSource(macFrame []byte) net.HardwareAddr {
	return net.HardwareAddr(macFrame[6:12])
}

// ///////////////ipv4////////////////
// 打印ip报文头的详情
func printIPv4Header(packet []byte) {
	fmt.Println()
	fmt.Println("----->IP Header<----")
	printVersionIPv4(packet)
	printHeaderLenIPv4(packet)
	printServiceTypeIPv4(packet)
	printAllLenIPv4(packet)
	printIdentificationIPv4(packet)
	printFlagsIPv4(packet)
	printFragmentOffsetIPv4(packet)
	printTTLIPv4(packet)
	printProtocolIPv4(packet)
	printChecksumIPv4(packet)
	printSrcIPv4(packet)
	printDstIPv4(packet)
	fmt.Println("----->IP Header<---End---")
}

func getIPv4HeaderLen(packet []byte) int {
	header := packet[0]
	headerLen := header & 0x0f * 4
	hl, _ := strconv.Atoi(fmt.Sprintf("%d", headerLen))
	return hl
}

func printVersionIPv4(packet []byte) {
	header := packet[0]
	fmt.Printf("IPv4 Header--->Version:%d\n", header>>4)
}

func printHeaderLenIPv4(packet []byte) {
	header := packet[0]
	fmt.Printf("IPv4 Header--->HeaderLen:%d byte\n", header&0x0f*4)
}

func printServiceTypeIPv4(packet []byte) {
	st := packet[1]
	fmt.Printf("IPv4 Header--->ServiceType:%v\n", st)
}

func printAllLenIPv4(packet []byte) {
	fmt.Printf("IPv4 Header--->AllLen:%d\n", uint16(packet[2])<<8|uint16(packet[3]))
}

func printIdentificationIPv4(packet []byte) {
	id := uint16(packet[4])<<8 | uint16(packet[5])
	fmt.Printf("IPv4 Header--->Identification:%x\t%d\n", id, id)
}

func printFlagsIPv4(packet []byte) {
	// 向右移动5位，相当于去掉后5位。数值降低
	fmt.Printf("IPv4 Header--->Flags:%03b\n", packet[6]>>5)
}
func printFragmentOffsetIPv4(packet []byte) {
	// 向左移动3位，相当于将前3位去掉。数值可能增加
	fmt.Printf("IPv4 Header--->FragmentOffset:%013b\n", uint16(packet[6])<<3|uint16(packet[7]))
}

func printTTLIPv4(packet []byte) {
	fmt.Printf("IPv4 Header--->TTL:%d\n", packet[8])
}

func printProtocolIPv4(packet []byte) {
	fmt.Printf("IPv4 Header--->ProtocolType:%d\n", packet[9])
}

func printChecksumIPv4(packet []byte) {
	fmt.Printf("IPv4 Header--->Checksum:%d\n", uint16(packet[10])<<8|uint16(packet[11]))
}

func printSrcIPv4(packet []byte) net.IP {
	fmt.Printf("IPv4 Header--->SrcIP:%d.%d.%d.%d\n", packet[12], packet[13], packet[14], packet[15])
	return net.IPv4(packet[12], packet[13], packet[14], packet[15])
}

func printDstIPv4(packet []byte) net.IP {
	fmt.Printf("IPv4 Header--->DstIP:%d.%d.%d.%d\n", packet[16], packet[17], packet[18], packet[19])
	return net.IPv4(packet[16], packet[17], packet[18], packet[19])
}

///////////////icmp////////////////

func printICMPHeader(packet []byte) {
	fmt.Println()
	fmt.Println("----->ICMP Header<----")
	printICMPType(packet)
	printICMPCode(packet)
	printICMPCheckSum(packet)
	printICMPIdentification(packet)
	printICMPSeqNum(packet)
	printICMPTimestamp(packet)
	printICMPData(packet)
	fmt.Println("----->ICMP Header<----End----")
}

func printICMPType(packet []byte) {
	fmt.Printf("ICMP Header--->Type:%d\n", packet[0])
}

func printICMPCode(packet []byte) {
	fmt.Printf("ICMP Header--->Code:%d\n", packet[1])
}

func printICMPCheckSum(packet []byte) {
	fmt.Printf("ICMP Header--->Checksum:%04x\n", uint16(packet[2])<<8|uint16(packet[3]))
}

func printICMPIdentification(packet []byte) {
	fmt.Printf("TICMP Header--->Identification(process):%d\n", uint16(packet[4])<<8|uint16(packet[5]))
}

func printICMPSeqNum(packet []byte) {
	fmt.Printf("ICMP Header--->SeqNum:%d\n", uint16(packet[6])<<8|uint16(packet[7]))
}

func printICMPTimestamp(packet []byte) {
	fmt.Printf("ICMP Header --->timestamp:%02x %02x %02x %02x %02x %02x %02x %02x\n", packet[8], packet[9], packet[10], packet[11], packet[12], packet[13], packet[14], packet[15])
}

func printICMPData(packet []byte) {
	fmt.Printf("ICMP Header--->data:%v\n", string(packet[20:]))
}

////////////////tcp///////////////

func printTCPHeader(packet []byte) {
	fmt.Println()
	fmt.Printf("----->TCP Header<----len:%d\n", len(packet[:]))
	printTCPsrcPort(packet)
	printTCPdstPort(packet)
	printTCPSequenceNumber(packet)
	printTCPACKNumber(packet)
	printTCPHeaderLen(packet)
	printTCPFlagCWR(packet)
	printTCPFlagECE(packet)
	printTCPFlagUrgent(packet)
	printTCPFlagACK(packet)
	printTCPFlagPSH(packet)
	printTCPFlagRST(packet)
	printTCPFlagSYN(packet)
	printTCPFlagFIN(packet)
	printTCPWindowSize(packet)
	printTCPCheckSum(packet)
	printTCPUrgentPointer(packet)
	printTCPData(packet)
	fmt.Println("----->TCP Header<---END---")
}

func printTCPsrcPort(packet []byte) {
	fmt.Printf("TCP Header--->SrcPort:%d\n", int64(uint16(packet[0])<<8|uint16(packet[1])))
}

func printTCPdstPort(packet []byte) {
	fmt.Printf("TCP Header--->DstPort:%d\n", int64(uint16(packet[2])<<8|uint16(packet[3])))
}

func printTCPSequenceNumber(packet []byte) {
	fmt.Printf("TCP Header--->SequenceNumber:%d\n", uint32(packet[4])<<24|uint32(packet[5])<<16|uint32(packet[6])<<8|uint32(packet[7]))
}

func printTCPACKNumber(packet []byte) {
	fmt.Printf("TCP Header--->ACKNum:%d\n", uint32(packet[8])<<24|uint32(packet[9])<<16|uint32(packet[10])<<8|uint32(packet[11]))
}

func printTCPHeaderLen(packet []byte) {
	fmt.Printf("TCP Header--->HeaderLen:%d\n", packet[12]>>4*4)
}

func printTCPFlagCWR(packet []byte) {
	fmt.Printf("TCP Header--->FlagCWR:%d\n", packet[13]&0x80>>7)
}
func printTCPFlagECE(packet []byte) {
	fmt.Printf("TCP Header--->FlagEcho:%d\n", packet[13]&0x40>>6)
}

func printTCPFlagUrgent(packet []byte) {
	fmt.Printf("TCP Header--->FlagUrgent:%d\n", packet[13]&0x20>>5)
}

func printTCPFlagACK(packet []byte) {
	fmt.Printf("TCP Header--->FlagACK:%d\n", packet[13]>>4&0b0001)
	fmt.Printf("TCP Header--->FlagACK:%d\n", packet[13]>>4&0b1)
	//fmt.Printf("TCP Header--->FlagACK:%d\tpacket:%v\n", packet[13]>>4, packet[13])
	fmt.Printf("TCP Header--->FlagACK:%d\n", packet[13]&0x10>>4)
}

func printTCPFlagPSH(packet []byte) {
	fmt.Printf("TCP Header--->FlagPSH:%d\n", packet[13]&0x08>>3)
}

func printTCPFlagRST(packet []byte) {
	fmt.Printf("TCP Header--->FlagRST:%d\n", packet[13]&0x04>>2)
}

func printTCPFlagSYN(packet []byte) {
	fmt.Printf("TCP Header--->FlagSYN:%d\n", packet[13]&0x02>>1)
}

func printTCPFlagFIN(packet []byte) {
	fmt.Printf("TCP Header--->FlagFIN:%d\n", packet[13]&0x01)
}

func printTCPWindowSize(packet []byte) {
	fmt.Printf("TCP Header--->WindowSize:%d\n", uint16(packet[14])<<8|uint16(packet[15]))
}

func printTCPCheckSum(packet []byte) {
	fmt.Printf("TCP Header--->Checksum:%04x\n", uint16(packet[16])<<8|uint16(packet[17]))
}

func printTCPUrgentPointer(packet []byte) {
	fmt.Printf("TCP Header--->UrgentPointer:%d\n", uint16(packet[18])<<8|uint16(packet[19]))
}

func printTCPData(packet []byte) {
	headerLen := packet[12] >> 4 * 4
	p := packet[headerLen:]
	dataLen := len(p)
	if dataLen > 1 {
		fmt.Printf("TCP Header--->Data--->:%v\theaderLen:%v\tdataLen:%d\n", string(p[:dataLen-1]), headerLen, dataLen)
	}
}

// ////////////udp/////////////////
func printUDPHeader(packet []byte) {
	fmt.Println()
	fmt.Println("----->UDP Header<----")
	printUDPSrcPort(packet)
	printUDPDstPort(packet)
	printUDPAllLen(packet)
	printUDPChecksum(packet)
	printUDPData(packet)
	fmt.Println("----->UDP Header<---End--")
}

func printUDPSrcPort(packet []byte) {
	fmt.Printf("UDP Header--->SrcPort:%d\n", uint16(packet[0])<<8|uint16(packet[1]))
}

func printUDPDstPort(packet []byte) {
	fmt.Printf("UDP Header--->DstPort:%d\n", uint16(packet[2])<<8|uint16(packet[3]))
}

func printUDPAllLen(packet []byte) {
	fmt.Printf("UDP Header--->AllLen:%d\n", uint16(packet[4])<<8|uint16(packet[5]))
}

func printUDPChecksum(packet []byte) {
	fmt.Printf("UDP Header--->Checksum:%d\n", uint16(packet[6])<<8|uint16(packet[7]))
}

func printUDPData(packet []byte) {
	fmt.Printf("UDP Header--->data:%v\n", string(packet[8:]))
}

//////////////arrp///////////////

func printARPHeader(packet []byte) {
	fmt.Println()
	fmt.Println("----->ARP Header<----")
	printARPTYPE(packet)
	printARPProTYPE(packet)
	printARPHardwareLen(packet)
	printARPPLen(packet)
	printARPOp(packet)
	printARPSrcHardwareMAC(packet)
	printARPSrcIP(packet)
	printARPDstHardwareMAC(packet)
	printARPDstIP(packet)
	fmt.Println("----->ARP Header<---END---")

}
func printARPTYPE(packet []byte) {
	fmt.Printf("ARP Header--->Type:%d\n", uint16(packet[0])<<8|uint16(packet[1]))
}

func printARPProTYPE(packet []byte) {
	fmt.Printf("ARP Header--->ProtocolType:%04x\n", uint16(packet[2])<<8|uint16(packet[3]))
}

func printARPHardwareLen(packet []byte) {
	fmt.Printf("ARP Header--->HardwareLen:%d\n", packet[4])
}

func printARPPLen(packet []byte) {
	fmt.Printf("ARP Header--->ProtocolLen:%d\n", packet[5])
}

func printARPOp(packet []byte) {
	fmt.Printf("ARP Header--->op:%d\n", uint16(packet[6])<<8|uint16(packet[7]))
}

func printARPSrcHardwareMAC(packet []byte) {
	fmt.Printf("ARP Header--->SrcHardwarMAC:%x:%x:%x:%x:%x:%x\n", packet[8], packet[9], packet[10], packet[11], packet[12], packet[13])
}

func printARPSrcIP(packet []byte) {
	fmt.Printf("ARP Header--->SrcIP:%d.%d.%d.%d\n", packet[14], packet[15], packet[16], packet[17])
}

func printARPDstHardwareMAC(packet []byte) {
	fmt.Printf("ARP Header--->DstHardwareMAC:%x:%x:%x:%x:%x:%x\n", packet[18], packet[19], packet[20], packet[21], packet[22], packet[23])
}

func printARPDstIP(packet []byte) {
	fmt.Printf("ARP Header--->DstIP:%d.%d.%d.%d\n", packet[24], packet[25], packet[26], packet[27])
}

func GetIPv4SrcARP(packet []byte) net.IP {
	return net.IPv4(packet[14], packet[15], packet[16], packet[17])
}

func GetIPv4DstARP(packet []byte) net.IP {
	return net.IPv4(packet[24], packet[25], packet[26], packet[27])
}

func getSrcIP(packet []byte) string {
	key := ""
	if isIPv4Pkt(packet) && len(packet) >= 20 {
		key = GetIPv4Src(packet).To4().String()
	}
	return key
}

func getDstIP(packet []byte) string {
	key := ""
	if isIPv4Pkt(packet) && len(packet) >= 20 {
		key = GetIPv4Dst(packet).To4().String()
	}
	return key
}
func isIPv4Pkt(packet []byte) bool {
	flag := packet[0] >> 4
	return flag == 4
}
