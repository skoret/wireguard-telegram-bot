package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

const wg0 = "wg0"

func main() {
	if err := WgShow(); err != nil {
		log.Fatalf("failed to run wg show: %v", err)
	}
	c, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			panic(err)
		}
	}()

	device, err := c.Device(wg0)
	if err != nil {
		log.Fatalf("failed to get device %q: %v", device, err)
	}
	ip, err := getLastUsedIP(device)
	if err != nil {
		log.Fatalf("failed to get last used ip: %v", err)
	}
	log.Printf("latest used ip: %s", ip.String())

	interfacesShow()
	netlinkShow()
}

func WgShow() error {
	log.Println("------ WgShow -----")
	cmd := exec.Command("wg", "show")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	log.Println("-------------------")
	return err
}

func interfacesShow() {
	log.Println("------ interfaces -----")
	ifis, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return
	}
	for _, ifi := range ifis {
		log.Printf("interface: %+v", ifi)
		addr, err := getDeviceAddress(ifi.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("address: %s", addr.String())
		log.Println("-------------------")
	}
}

func netlinkShow() {
	log.Println("------ netlink info -----")
	c, err := genetlink.Dial(nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() { _ = c.Close() }()

	fmls, err := c.ListFamilies()
	if err != nil {
		log.Println(err)
		return
	}
	for _, fml := range fmls {
		log.Printf("%+v", fml)
		log.Println("-------------------")
	}
	ifis, err := rtnlInterfaces()
	if err != nil {
		log.Println(err)
		return
	}
	for _, ifi := range ifis {
		log.Println(ifi)
	}
}

// rtnlInterfaces uses rtnetlink to fetch a list of WireGuard interfaces.
func rtnlInterfaces() ([]string, error) {
	// Use the stdlib's rtnetlink helpers to get ahold of a table of all
	// interfaces, so we can begin filtering it down to just WireGuard devices.
	tab, err := syscall.NetlinkRIB(unix.RTM_GETLINK, unix.AF_UNSPEC)
	if err != nil {
		return nil, fmt.Errorf("wglinux: failed to get list of interfaces from rtnetlink: %v", err)
	}

	msgs, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, fmt.Errorf("wglinux: failed to parse rtnetlink messages: %v", err)
	}

	return parseRTNLInterfaces(msgs)
}

// parseRTNLInterfaces unpacks rtnetlink messages and returns WireGuard
// interface names.
func parseRTNLInterfaces(msgs []syscall.NetlinkMessage) ([]string, error) {
	log.Println("------ parseRTNLInterfaces -----")
	var ifis []string
	for i, m := range msgs {
		// Only deal with link messages, and they must have an ifinfomsg
		// structure appear before the attributes.
		if m.Header.Type != unix.RTM_NEWLINK {
			continue
		}

		if len(m.Data) < unix.SizeofIfInfomsg {
			return nil, fmt.Errorf("wglinux: rtnetlink message is too short for ifinfomsg: %d", len(m.Data))
		}

		ad, err := netlink.NewAttributeDecoder(m.Data[syscall.SizeofIfInfomsg:])
		if err != nil {
			return nil, err
		}

		// Determine the interface's name and if it's a WireGuard device.
		var (
			ifi  string
			isWG bool
		)

		log.Printf("parse #%d msg:", i)
		for ad.Next() {
			log.Printf("type: %d | data: %s", ad.Type(), ad.String())
			switch ad.Type() {
			case unix.IFLA_IFNAME:
				ifi = ad.String()
			case unix.IFLA_LINKINFO:
				ad.Do(isWGKind(&isWG))
			}
		}

		if err := ad.Err(); err != nil {
			return nil, err
		}

		if isWG {
			// Found one; append it to the list.
			log.Println("it's wireguard interface")
			ifis = append(ifis, ifi)
		} else {
			log.Println("it's NOT wireguard interface")
		}
		log.Println("-------------------")
	}

	return ifis, nil
}

// wgKind is the IFLA_INFO_KIND value for WireGuard devices.
const wgKind = "wireguard"

// isWGKind parses netlink attributes to determine if a link is a WireGuard
// device, then populates ok with the result.
func isWGKind(ok *bool) func(b []byte) error {
	return func(b []byte) error {
		ad, err := netlink.NewAttributeDecoder(b)
		if err != nil {
			return err
		}

		for ad.Next() {
			if ad.Type() != unix.IFLA_INFO_KIND {
				continue
			}

			if ad.String() == wgKind {
				*ok = true
				return nil
			}
		}

		return ad.Err()
	}
}

func getDeviceAddress(name string) (net.IP, error) {
	log.Println("------ getDeviceAddress -----")
	ief, err := net.InterfaceByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get interface "+name)
	}
	log.Printf("interface: %+v", ief)
	log.Println("-------------------")
	addrs, err := ief.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get address for interface "+name)
	}
	for _, addr := range addrs {
		log.Printf("%T: %+v", addr, addr)
		if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			return ipv4Addr, nil
		}
	}
	log.Println("-------------------")
	return nil, nil
}

func getLastUsedIP(device *wgtypes.Device) (net.IP, error) {
	log.Println("------ getLastUsedIP -----")
	lastIP, err := getDeviceAddress(device.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device address")
	}
	for _, peer := range device.Peers {
		for _, ipNet := range peer.AllowedIPs {
			if bytes.Compare(ipNet.IP, lastIP) >= 0 {
				lastIP = ipNet.IP
			}
		}
	}
	if lastIP.Equal(net.ParseIP("0.0.0.0")) {
		return nil, errors.New("failed to get la")
	}
	return lastIP, nil
}
