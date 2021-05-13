package wireguard

import (
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	cfgs "github.com/skoret/wireguard-bot/internal/wireguard/configs"
)

type Wireguard struct {
	device string
	dns    []string
	client *wgctrl.Client
}

func NewWireguard() (*Wireguard, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	devs, err := client.Devices()
	if err != nil {
		return nil, err
	}
	log.Printf("--- known devices: ---")
	for i, d := range devs {
		log.Printf("#%d device: %+v", i, d)
	}
	log.Printf("----------------------")
	return &Wireguard{
		device: os.Getenv("PUBLIC_INTERFACE"),
		dns:    strings.Split(os.Getenv("DNS_IPS"), ","),
		client: client,
	}, nil
}

func (w *Wireguard) Close() error {
	return w.client.Close()
}

func (w *Wireguard) CreateConfigForNewKeys() (io.Reader, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}
	ipNet, err := w.getNextIPNet()
	if err != nil {
		return nil, err
	}

	cfgFile, err := w.createConfig(pri.String(), ipNet)
	if err != nil {
		panic(err)
	}

	// wg server conf update
	if err := w.updateDevice(pri.PublicKey(), ipNet); err != nil {
		return nil, err
	}
	return cfgFile, nil
}

func (w *Wireguard) CreateConfigForPublicKey(key string) (io.Reader, error) {
	pub, err := wgtypes.ParseKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}
	ipNet, err := w.getNextIPNet()
	if err != nil {
		return nil, err
	}

	cfgFile, err := w.createConfig("", ipNet)
	if err != nil {
		panic(err)
	}

	// wg server conf update
	if err := w.updateDevice(pub, ipNet); err != nil {
		return nil, err
	}
	return cfgFile, nil
}

func (w *Wireguard) createConfig(pri string, ipNet *net.IPNet) (io.Reader, error) {
	device, err := w.client.Device(w.device)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device "+w.device)
	}
	clientConfig := cfgs.ClientConfig{
		Address:    ipNet.String(),
		PrivateKey: pri,
		DNS:        w.dns,

		PublicKey:  device.PublicKey.String(),
		AllowedIPs: []string{"0.0.0.0/0"},
		Endpoint:   os.Getenv("SERVER_ENDPOINT"),
	}
	cfgFile, err := cfgs.ProcessClientConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	return cfgFile, nil
}

func (w *Wireguard) updateDevice(pub wgtypes.Key, ipNet *net.IPNet) error {
	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:  pub,
				AllowedIPs: []net.IPNet{*ipNet},
			},
		},
	}
	if err := w.client.ConfigureDevice(w.device, cfg); err != nil {
		return errors.Wrap(err, "failed to update server configuration")
	}
	cmd := exec.Command("wg-quick", "save", w.device)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to dump server config to conf file")
	}
	return nil
}

func (w *Wireguard) getNextIPNet() (*net.IPNet, error) {
	ip, err := w.getLatestUsedIP()
	if err != nil {
		return nil, err
	}

	return &net.IPNet{
		IP:   nextIP(ip, 1),
		Mask: net.IPv4Mask(255, 255, 255, 255),
	}, nil
}

// Thanks to https://gist.github.com/udhos/b468fbfd376aa0b655b6b0c539a88c03
func nextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

func (w *Wireguard) getLatestUsedIP() (net.IP, error) {
	device, err := w.client.Device(w.device)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get device "+w.device)
	}
	lastIP, err := w.getDeviceAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get address for device "+w.device)
	}
	for _, peer := range device.Peers {
		for _, ipNet := range peer.AllowedIPs {
			if bytes.Compare(ipNet.IP, lastIP) >= 0 {
				lastIP = ipNet.IP
			}
		}
	}
	if lastIP == nil {
		return nil, errors.New("failed to get latest used ip for device " + w.device)
	}
	return lastIP, nil
}

func (w *Wireguard) getDeviceAddress() (net.IP, error) {
	ife, err := net.InterfaceByName(w.device)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get interface "+w.device)
	}
	addrs, err := ife.Addrs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get address for interface "+w.device)
	}
	for _, addr := range addrs {
		if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			return ipv4Addr, nil
		}
	}
	return nil, errors.New("failed to get address for interface " + w.device)
}
