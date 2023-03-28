package ip

import (
	"math"
	"net"
	"strings"
)

// InternalIP return internal ip.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

func String2Long(ip string) int64 {
	b := net.ParseIP(ip).To4()
	if b == nil {
		return 0
	}
	return int64(b[3]) | int64(b[2])<<8 | int64(b[1])<<16 | int64(b[0])<<24
}

// Long2String 把数值转为ip字符串
func Long2String(i int64) string {
	if i > math.MaxUint32 {
		return ""
	}
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)
	return ip.String()
}
