package module

import (
	"Project/LearnGo/Learn/webcrawler/errors"
	"fmt"
	"net"
	"strconv"
)

// 组件网络地址的类型
type mAddr struct {
	// 网络协议
	network string
	// 网络地址
	address string
}

// 获取访问组件时要遵循的网络协议
func (maddr *mAddr) Network() string {
	return maddr.network
}

// 获取组件的网络地址
func (maddr *mAddr) String() string {
	return maddr.address
}

// 根据参数创建并返回一个网络地址值
func NewAddr(network string, ip string, port uint64) (net.Addr, error) {
	if network != "http" && network != "https" {
		errMsg := fmt.Sprintf("illegal network for module address: %s", network)
		return nil, errors.NewIllegalParameterError(errMsg)
	}

	if parsedIP := net.ParseIP(ip); parsedIP == nil {
		errMsg := fmt.Sprintf("illegal IP for module address: %s", ip)
		return nil, errors.NewIllegalParameterError(errMsg)
	}
	return &mAddr{network: network, address: ip + ":" + strconv.Itoa(int(port))}, nil
}
