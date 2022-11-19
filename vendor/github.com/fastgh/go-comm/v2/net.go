package comm

import (
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
)

func HostnameP() string {
	r, err := Hostname()
	if err != nil {
		panic(errors.Wrapf(err, "get hostname"))
	}
	return r
}

func Hostname() (string, error) {
	return os.Hostname()
}

func BroadcastInterfacesP(dump bool) []net.Interface {
	r, err := BroadcastInterfaces(dump)
	if err != nil {
		panic(err)
	}
	return r
}

func BroadcastInterfaces(dump bool) ([]net.Interface, error) {
	netIfs, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "get network interfaces")
	}

	r := make([]net.Interface, 0, len(netIfs))
	for _, netIf := range netIfs {
		flag := netIf.Flags
		if (flag | net.FlagUp) == 0 {
			// ignore because it is down
			continue
		}
		if (flag | net.FlagBroadcast) == 0 {
			// ignore non-broadcast network interface
			continue
		}

		if dump {
			// TODO: single log
			fmt.Printf("candidate interface: %s\n", netIf.Name)
		}

		r = append(r, netIf)
	}

	return r, nil
}

func BroadcastIpWithInterfaceP(intf net.Interface) net.IP {
	r, err := BroadcastIpWithInterface(intf)
	if err != nil {
		panic(err)
	}
	return r
}

func BroadcastIpWithInterface(intf net.Interface) (net.IP, error) {
	// intf.MulticastAddrs()
	addrs, err := intf.Addrs()
	if err != nil {
		return nil, errors.Wrapf(err, "get addresses for interface: %s", intf.Name)
	}

	for _, addr := range addrs {
		if ipAddr, isIpAddr := addr.(*net.IPNet); isIpAddr {
			ip := ipAddr.IP
			if !ip.IsLoopback() && ip.To4() != nil {
				return ip, nil
			}
		}
	}

	return nil, nil
}

func ResolveBroadcastIpP(interfaces []net.Interface, interfaceName string) (net.IP, net.IP) {
	localIp, broadcastIp, err := ResolveBroadcastIp(interfaces, interfaceName)
	if err != nil {
		panic(err)
	}
	return localIp, broadcastIp
}

func ResolveBroadcastIp(interfaces []net.Interface, interfaceName string) (net.IP, net.IP, error) {
	for _, intF := range interfaces {
		if intF.Name == interfaceName {
			localIp, err := BroadcastIpWithInterface(intF)
			if localIp == nil {
				return nil, nil, errors.Wrapf(err, "cannot get a broadcast ip for interface %s", interfaceName)
			}

			broadcastIp := make(net.IP, len(localIp))
			copy(broadcastIp, localIp)
			broadcastIp[len(broadcastIp)-1] = 255
			return localIp, broadcastIp, nil
		}
	}

	return nil, nil, fmt.Errorf("interface %s is not found, or down, or not supports broadcast", interfaceName)
}
