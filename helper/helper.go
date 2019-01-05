package helper

import (
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/govindarajan/laserproxy/logger"
)

type LocalIPAddr struct {
	IFace   string
	IP      net.IP
	Gateway net.IP
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
					lip := LocalIPAddr{IP: ipnet.IP, IFace: iface.Name}
					lip.Gateway = getGateway(iface.Name)
					if lip.Gateway == nil {
						continue
					}
					LIPs = append(LIPs, lip)
				}
			}
		}
	}
	return LIPs, nil
}

func getGateway(iface string) net.IP {

	out, err := exec.Command("ip", "route", "show").Output()
	if err != nil {
		logger.LogError(err.Error())
		return nil
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 4 && words[0] == "default" &&
			words[4] == iface {
			return net.ParseIP(words[2])
		}
	}
	//logger.LogError("Not able to get Gateway")
	return nil
}

//GetHostIPs get the list of IPs of a hostname
func GetHostIPs(hostname string) ([]string, error) {
	ip, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	return ip, nil
}

// WatchNetworkChange will keep checking the network change.
// If detected, it send signal in the change channel.
// This is not thread-safe.
func WatchNetworkChange(checkIntlSec int, change chan bool) {
	var existingIPs []LocalIPAddr
	for {
		time.Sleep(time.Duration(checkIntlSec) * time.Second)
		curIPs, err := GetLocalIPs()
		if err != nil {
			logger.LogError("Error while getting local IPs")
			continue
		}
		if existingIPs == nil {
			// First time. Lets store it.
			existingIPs = curIPs
			continue
		}

		if len(existingIPs) != len(curIPs) {
			change <- true
		} else {
			for i, _ := range curIPs {
				if !curIPs[i].IP.Equal(existingIPs[i].IP) {
					change <- true
					break
				}
			}
		}
		existingIPs = curIPs
	}
}
