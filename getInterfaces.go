package main

import (
	"fmt"
	"net"
)

// Check if the Given IP is valid or not
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

// List all the network interfaces which are present in server
func LocalInterfaces() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("Local interfaces: %+v\n", err.Error()))
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("Local interfaces: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				fmt.Printf("Name %s\n", v)
			}
		}
	}
}

func main() {
	LocalInterfaces()
}
