package wireguard

import (
	"io"
	"log"

	cfgs "github.com/skoret/wireguard-bot/internal/wireguard/configs"
)

type devwireguard struct{}

func NewDevWireguard() (Wireguard, error) {
	log.Println("--- create dummy dev wireguard client ---")
	return &devwireguard{}, nil
}

func (d *devwireguard) Close() error {
	log.Println("dev wireguard closed")
	return nil
}

func (d *devwireguard) CreateConfigForNewKeys() (io.Reader, error) {
	log.Println("dev wireguard creates dummy config")
	return cfgs.ProcessClientConfig(cfg)
}

func (d *devwireguard) CreateConfigForPublicKey(string) (io.Reader, error) {
	return d.CreateConfigForNewKeys()
}

var cfg = cfgs.ClientConfig{
	Address:    "<peer_ip>",
	PrivateKey: "<private_key>",
	DNS:        []string{"<dns>"},

	PublicKey:  "<public_key>",
	AllowedIPs: []string{"<allowed_ip>"},
	Endpoint:   "<server_endpoint>",
}
