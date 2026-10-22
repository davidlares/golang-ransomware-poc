package infosys

import (
	"net"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func GetUserName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "Error getting hostname"
	}
	return hostname
}

func GetHostName() string {
	var username [256]uint16
	size := uint32(len(username))
	err := windows.GetUserNameEx(windows.NameSamCompatible, &username[0], &size)
	if err != nil {
		return "Error getting username"
	}
	return syscall.UTF16ToString(username[:])
}

func GetIPAddress() ([]string, error) {
	// array of IPs
	var ips []string
	// listing addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err // Returns nil slice and the error
	}
	// looping
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipNet.IP.String())
		}
	}
	// returning IP addresses
	return ips, nil
}
