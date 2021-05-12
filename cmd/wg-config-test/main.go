package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/pkg/errors"
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

	_, err = getDeviceAddress(device)
	if err != nil {
		log.Fatalf("failed to get device address: %v", err)
	}
	ip, err := getLastUsedIP(device)
	if err != nil {
		log.Fatalf("failed to get last used ip: %v", err)
	}
	log.Printf("latest used ip: %s", ip.String())
}

func WgShow() error {
	fmt.Println("------ WgShow -----")
	cmd := exec.Command("wg", "show")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println("-------------------")
	return err
}

// TODO WIP
func getDeviceAddress(device *wgtypes.Device) (net.IP, error) {
	fmt.Println("------ getDeviceAddress -----")
	ief, err := net.InterfaceByName(device.Name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get interface "+device.Name)
	}
	log.Printf("interface: %+v", ief)
	fmt.Println("-------------------")
	addrs, err := ief.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get address for interface "+device.Name)
	}
	for _, addr := range addrs {
		log.Printf("%T: %+v", addr, addr)
	}
	return nil, nil
}

func getLastUsedIP(device *wgtypes.Device) (net.IP, error) {
	fmt.Println("------ getLastUsedIP -----")
	lastIP := net.ParseIP("0.0.0.0")
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
