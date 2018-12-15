package helper

import (
	"net"

	"github.com/govindarajan/laserproxy/logger"
)

//GetLocalIPs returns various interface ips for the hostname passed
func GetLocalIPs() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.LogError("Error getting IP addresses: %+v")
		return nil, err
	}
	var IPs []string
	for _, ip := range ifaces {
		addrs, err := ip.Addrs()
		if err != nil {
			logger.LogError("Error fetching IP Addresses: %+v" + err.Error())
			return IPs, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					IPs = append(IPs, ipnet.IP.String())
				}
			}
		}
	}
	return IPs, nil
}
