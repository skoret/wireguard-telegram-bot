package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"golang.zx2c4.com/wireguard/wgctrl"
)

const wg0 = "wg0"

func main() {
	flag.Parse()

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

	var maxIp *net.IPNet
	for _, peer := range device.Peers {
		for _, ipNet := range peer.AllowedIPs {
			if maxIp == nil {
				maxIp = &ipNet
				continue
			}
			if ipNet.IP.To4()[3] > maxIp.IP.To4()[3] {
				maxIp = &ipNet
			}
		}
	}
	if maxIp == nil {
		panic("puk")
	}

	if err := WgShow(); err != nil {
		panic(err)
	}

	log.Printf("max ip now is: %v", maxIp.String())
	maxIp.IP.To4()[3] += 1
	log.Printf("next ip is: %v", maxIp.String())
}

func WgShow() error {
	fmt.Println("------ WgShow -----")
	cmd := exec.Command("wg", "show")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println("-------------------")
	return err
}
