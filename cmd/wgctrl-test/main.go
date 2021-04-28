package main

import (
	"flag"
	"fmt"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"log"
	"net"
	"os"
	"strings"
)

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

	device, err := c.Device("wg0")
	if err != nil {
		log.Fatalf("failed to get device %q: %v", device, err)
	}
	show(device)

	fmt.Println("---- add new peer ----")
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}
	pub := pri.PublicKey()
	_, ipNet, err := net.ParseCIDR("10.8.0.3/32")
	if err != nil {
		log.Fatalf("failed to parse ip with mask: %v", err)
	}
	log.Printf("private key: %s", pri.String())
	log.Printf("public  key: %s", pub.String())

	cfg := wgtypes.Config{
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:                   pub,
				Remove:                      false,
				UpdateOnly:                  false,
				PresharedKey:                nil,
				Endpoint:                    nil,
				PersistentKeepaliveInterval: nil,
				ReplaceAllowedIPs:           false,
				AllowedIPs:                  []net.IPNet{*ipNet},
			},
		},
	}

	if err := c.ConfigureDevice("wg0", cfg); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(err)
		} else {
			log.Fatalf("Unknown config error: %v", err)
		}
	}

	show(device)

	fmt.Println("we need to run 'wg-quick save wg0' to dump updated interface config to config file")
	//cmd := exec.Command("wg-quick", "save")
}

func show(d *wgtypes.Device) {
	fmt.Println("---- wg device ----")
	printDevice(d)
	for _, p := range d.Peers {
		printPeer(p)
	}
	fmt.Println("-------------------")
}

func printDevice(d *wgtypes.Device) {
	const f = `interface: %s (%s)
  public key: %s
  private key: (hidden)
  listening port: %d

`

	fmt.Printf(
		f,
		d.Name,
		d.Type.String(),
		d.PublicKey.String(),
		d.ListenPort)
}

func printPeer(p wgtypes.Peer) {
	const f = `peer: %s
  endpoint: %s
  allowed ips: %s
  latest handshake: %s
  transfer: %d B received, %d B sent

`

	fmt.Printf(
		f,
		p.PublicKey.String(),
		// TODO(mdlayher): get right endpoint with getnameinfo.
		p.Endpoint.String(),
		ipsString(p.AllowedIPs),
		p.LastHandshakeTime.String(),
		p.ReceiveBytes,
		p.TransmitBytes,
	)
}

func ipsString(ipns []net.IPNet) string {
	ss := make([]string, 0, len(ipns))
	for _, ipn := range ipns {
		ss = append(ss, ipn.String())
	}

	return strings.Join(ss, ", ")
}
