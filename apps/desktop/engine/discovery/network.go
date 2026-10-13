package discovery

import (
	"net"
)

// getBroadcastAddresses returns all candidate IPv4 broadcast destinations:
// the global broadcast address (255.255.255.255) plus the directed subnet
// broadcast address for each active, non-loopback network interface.
func getBroadcastAddresses() []net.IP {
	addrs := make([]net.IP, 0)
	seen := make(map[string]bool)

	// Always include universal broadcast
	universal := net.IPv4(255, 255, 255, 255)
	addrs = append(addrs, universal)
	seen[universal.String()] = true

	ifaces, err := net.Interfaces()
	if err != nil {
		return addrs
	}

	for _, iface := range ifaces {
		// Skip interfaces that are down or loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ifAddrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range ifAddrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip4 := ipNet.IP.To4()
			if ip4 == nil {
				continue
			}

			// Compute subnet broadcast: IP | ^Mask
			mask := ipNet.Mask
			if len(mask) == 4 {
				bcast := make(net.IP, 4)
				for i := 0; i < 4; i++ {
					bcast[i] = ip4[i] | ^mask[i]
				}
				key := bcast.String()
				if !seen[key] {
					seen[key] = true
					addrs = append(addrs, bcast)
				}
			}
		}
	}

	return addrs
}

// getPreferredOutboundIP probes the local network to determine the preferred
// non-loopback local IPv4 address. If no active interface is found, returns "".
func getPreferredOutboundIP() string {
	conn, err := net.Dial("udp4", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
		if ok && localAddr.IP.To4() != nil {
			return localAddr.IP.String()
		}
	}

	// Fallback: enumerate interface addresses
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					return ip4.String()
				}
			}
		}
	}

	return ""
}
