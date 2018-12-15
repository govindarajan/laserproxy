package helper

import (
	"net"
	"os/exec"
	"strings"

	"github.com/govindarajan/laserproxy/logger"
)

type LocalIPAddr struct {
	IFace   string
	IP      string
	Gateway string
}

//GetLocalIPs returns various interface ips for the hostname passed
// TODO: Get it from cache
func GetLocalIPs() ([]LocalIPAddr, error) {
	LIPs := make([]LocalIPAddr, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.LogError("Error getting IP addresses: %+v")
		return nil, err
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			logger.LogError("Error fetching IP Addresses: %+v" + err.Error())
			return LIPs, err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					lip := LocalIPAddr{IP: ipnet.IP.String(), IFace: iface.Name}
					lip.Gateway = getGateway(iface.Name)
					if lip.Gateway == "" {
						continue
					}
					LIPs = append(LIPs, lip)
				}
			}
		}
	}
	return LIPs, nil
}

func getGateway(iface string) string {

	out, err := exec.Command("ip", "route", "show").Output()
	if err != nil {
		logger.LogError(err.Error())
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 4 && words[0] == "default" &&
			words[4] == iface {
			return words[2]
		}
	}
	//logger.LogError("Not able to get Gateway")
	return ""
}
