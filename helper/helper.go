package helper

import (
	"fmt"
	"net"
)

//GetLocalIPs returns various interface ips for the hostname passed
func GetLocalIPs() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("Error getting IP addresses: %+v\n", err.Error()))
		return nil
	}
	var IPs []string
	for _, ip := range ifaces {
		addrs, err := ip.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("Error fetching IP Addresses: %+v\n", err.Error()))
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					IPs = append(IPs, ipnet.IP.String())
				}
			}
		}
	}
	return IPs
}
