package wireguard

import (
	"io"
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
	return &Wireguard{
		device: os.Getenv("PUBLIC_INTERFACE"),
		dns:    strings.Split(os.Getenv("DNS_IPS"), ","),
		client: client,
	}, nil
}

func (w *Wireguard) Close() error {
	return w.client.Close()
}

func (w *Wireguard) CreateNewConfig() (io.Reader, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}
	// TODO: get ip net dynamically
	ipNet, err := w.getNextIPNet()
	if err != nil {
		return nil, err
	}

	clientConfig := cfgs.ClientConfig{
		Address:    ipNet.String(),
		PrivateKey: pri.String(),
		DNS:        []string{"8.8.8.8", "8.8.4.4"},

		PublicKey:  os.Getenv("SERVER_PUB_KEY"),
		AllowedIPs: []string{"0.0.0.0/0"},
		Endpoint:   os.Getenv("SERVER_ENDPOINT"),
	}
	cfgFile, err := cfgs.ProcessClientConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// wg server conf update
	pub := pri.PublicKey()
	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:  pub,
				AllowedIPs: []net.IPNet{*ipNet},
			},
		},
	}
	if err := w.updateDevice(cfg); err != nil {
		return nil, err
	}
	return cfgFile, nil
}

// TODO: get ip net dynamically
func (w *Wireguard) getNextIPNet() (*net.IPNet, error) {
	address := "10.8.0.3/32"
	_, ipNet, err := net.ParseCIDR(address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ip with mask")
	}
	return ipNet, nil
}

func (w *Wireguard) updateDevice(cfg wgtypes.Config) error {
	if err := w.client.ConfigureDevice(w.device, cfg); err != nil {
		return errors.Wrap(err, "failed to update server configuration")
	}
	cmd := exec.Command("wg-quick", "save", "wg0")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to dump server config to conf file")
	}
	return nil
}
