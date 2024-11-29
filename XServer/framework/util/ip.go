package util

import (
	"net"
)

// tcp/ip协议中，专门保留了三个IP地址区域作为私有地址，其地址范围如下：
// 10.0.0.0/8：10.0.0.0～10.255.255.255
// 172.16.0.0/12：172.16.0.0～172.31.255.255
// 192.168.0.0/16：192.168.0.0～192.168.255.255
// IsPublicIP ...
func IsPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	switch true {
	case ip4[0] == 10:
		return false
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
		return false
	case ip4[0] == 192 && ip4[1] == 168:
		return false
	default:
		return true
	}
}

// GetIntranetIp ...
// TODO 多网卡的情况下获取IP不准确(目前返回第一个有效值)
func GetIntranetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	addr := ""
	for _, v := range addrs {
		// 检查ip地址判断是否回环地址
		ipNet, ok := v.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.IsLoopback() {
			continue
		}
		if ipNet.IP.To4() == nil {
			continue
		}
		addr = ipNet.IP.String()
		break
	}

	return addr, nil
}
