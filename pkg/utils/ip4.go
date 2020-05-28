package utils

import (
	"net"
)

type IPv4 struct {
	ip net.IP
}

func newIPv4(ip net.IP) (result *IPv4) {
	ip = ip.To4()
	if ip != nil {
		result = &IPv4{ip}
	}
	return
}

func (ip4 *IPv4) String() string {
	return ip4.ip.String()
}

func (ip4 *IPv4) IsPrivate() bool {
	return ip4.IsPrivateClassA() || ip4.IsPrivateClassB() || ip4.IsPrivateClassC()
}

func (ip4 *IPv4) IsPrivateClassA() bool {
	return privateClassA.member(ip4)
}

func (ip4 *IPv4) IsPrivateClassB() bool {
	return privateClassB.member(ip4)
}

func (ip4 *IPv4) IsPrivateClassC() bool {
	return privateClassC.member(ip4)
}

//noinspection GoUnusedExportedFunction
func IPv4s() ([]IPv4, error) {
	return IPv4sFrom(net.InterfaceAddrs())
}

func IPv4sFrom(addrs []net.Addr, existingErr error) (ip4s []IPv4, err error) {
	if err = existingErr; err != nil {
		return
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok {
			ip := newIPv4(ipNet.IP)
			if ip != nil {
				ip4s = append(ip4s, *ip)
			}
		}
	}
	return
}

type privateClass struct {
	Byte1, Byte2Low, Byte2High byte
}

var privateClassA = privateClass{Byte1: 10, Byte2Low: 0, Byte2High: 255}
var privateClassB = privateClass{Byte1: 172, Byte2Low: 16, Byte2High: 31}
var privateClassC = privateClass{Byte1: 192, Byte2Low: 168, Byte2High: 168}

func (pC *privateClass) member(ipv4 *IPv4) bool {
	if ipv4 != nil {
		if ipv4.ip[0] == pC.Byte1 {
			b2 := ipv4.ip[1]
			return (pC.Byte2Low <= b2) && (b2 <= pC.Byte2High)
		}
	}
	return false
}
