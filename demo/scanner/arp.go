package main

import (
	"context"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	manuf "github.com/timest/gomanuf"
	"net"
	"strings"
	"time"
)

func listenARP(ctx context.Context) {
	handle, err := pcap.OpenLive(iface, 1024, false, 10*time.Second) // 打开handler
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()
	handle.SetBPFFilter("arp")                                // 只过去arp协议的东西
	ps := gopacket.NewPacketSource(handle, handle.LinkType()) // 解析监听到的包
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets(): // ps通道中来了一个包
			arp := p.Layer(layers.LayerTypeARP).(*layers.ARP) // 进行解析
			if arp.Operation == 2 {
				mac := net.HardwareAddr(arp.SourceHwAddress)                  // 解析mac地址
				m := manuf.Search(mac.String())                               // 根据mac地址，找到设备的制造厂商
				pushData(ParseIP(arp.SourceProtAddress).String(), mac, "", m) // 将解析的数据放到data通道，输出结果，并重置计时器
				if strings.Contains(m, "Apple") {
					go sendMdns(ParseIP(arp.SourceProtAddress), mac) // 如果是mac电脑，发送mdns（UDP？）
				} else {
					go sendNbns(ParseIP(arp.SourceProtAddress), mac) // 如果是linux或者windows，发送nbns（也是UDP）
				}
			}
		}
	}
}

// 发送arp包
// ip 目标IP地址
func sendArpPackage(ip IP) {
	srcIp := net.ParseIP(ipNet.IP.String()).To4() // 本机ip
	dstIp := net.ParseIP(ip.String()).To4()       // 目标ip
	if srcIp == nil || dstIp == nil {
		log.Fatal("ip 解析出问题")
	}
	// 以太网首部
	// EthernetType 0x0806  ARP
	ether := &layers.Ethernet{
		SrcMAC:       localHaddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // 0x0001 arp request 0x0002 arp response
		SourceHwAddress:   localHaddr,
		SourceProtAddress: srcIp,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIp,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, ether, a)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(iface, 2048, false, 30*time.Second)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal("发送arp数据包失败..")
	}
}
