package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

// ExtractHostPort from address
func ExtractHostPort(addr string) (host string, port uint64, err error) {
	var ports string
	host, ports, err = net.SplitHostPort(addr)
	if err != nil {
		return
	}
	port, err = strconv.ParseUint(ports, 10, 16) //nolint:gomnd
	return
}

func isValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

func isPrivateIP(addr string) bool {
	ip := net.ParseIP(addr)
	if ip4 := ip.To4(); ip4 != nil {
		// Following RFC 4193, Section 3. Local IPv6 Unicast Addresses which says:
		//   The Internet Assigned Numbers Authority (IANA) has reserved the
		//   following three blocks of the IPv4 address space for private internets:
		//     10.0.0.0        -   10.255.255.255  (10/8 prefix)
		//     172.16.0.0      -   172.31.255.255  (172.16/12 prefix)
		//     192.168.0.0     -   192.168.255.255 (192.168/16 prefix)
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// Following RFC 4193, Section 3. Private Address Space which says:
	//   The Internet Assigned Numbers Authority (IANA) has reserved the
	//   following block of the IPv6 address space for local internets:
	//     FC00::  -  FDFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF (FC00::/7 prefix)
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}

// Port return a real port.
func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

// Extract returns a private addr and port.
func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil && lis == nil {
		return "", err
	}
	if lis != nil {
		p, ok := Port(lis)
		if !ok {
			return "", fmt.Errorf("failed to extract port: %v", lis.Addr())
		}
		port = strconv.Itoa(p)
	}
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	minIndex := int(^uint(0) >> 1)
	ips := make([]net.IP, 0)
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index >= minIndex && len(ips) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for i, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isValidIP(ip.String()) {
				minIndex = iface.Index
				if i == 0 {
					ips = make([]net.IP, 0, 1)
				}
				ips = append(ips, ip)
				if ip.To4() != nil {
					break
				}
			}
		}
	}
	if len(ips) != 0 {
		return net.JoinHostPort(ips[len(ips)-1].String(), port), nil
	}
	return "", nil
}

// ExtractV1 returns a private addr and port.
func ExtractV1(hostPort string) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", err
	}
	if len(addr) > 0 && (addr != "0.0.0.0" &&
		addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isPrivateIP(ip.String()) {
				return net.JoinHostPort(ip.String(), port), nil
			}
		}
	}
	return "", nil
}

// IsLocalIP 检查请求的ip是否是本地ip
func IsLocalIP(host string) bool {
	ip, _, err := ExtractHostPort(host)
	if err != nil {
		return false
	}

	allowIps := []string{"localhost", "127.0.0.1"}
	for _, item := range allowIps {
		if ip == item {
			return true
		}
	}
	return false
}

// GetHostname get hostname
func GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	return name
}

// GetLocalHTTPAddrPairs get available http server and request address
func GetLocalHTTPAddrPairs() (serverAddr string, requestAddr string) {
	port, err := GetAvailablePort()
	if err != nil {
		fmt.Printf("GetAvailablePort error: %v\n", err)
		return "", ""
	}
	serverAddr = fmt.Sprintf(":%d", port)
	requestAddr = fmt.Sprintf("http://127.0.0.1:%d", port)
	return serverAddr, requestAddr
}

// GetAvailablePort get available port
func GetAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}
