package main

import (
	"fmt"
	"io"
	"os"

	cfgs "github.com/skoret/wireguard-bot/internal/wireguard/configs"
)

func main() {
	clientConfig := cfgs.ClientConfig{
		Address: "10.8.0.2/24",
		//PrivateKey: "aGsGuo9ODki0ZpS1U3c28tsI6UWjCW1Gbn8lIYRamXA=",
		DNS: []string{"8.8.8.8", "8.8.4.4"},

		PublicKey:  "G8naBU85RGmh2iZBi2KL3qomJOGKy5jvU97bO2I5tQ4=",
		AllowedIPs: []string{"0.0.0.0/0"},
	}

	serverConfig := cfgs.ServerConfig{
		Address:      "10.8.0.1/24",
		SaveConfig:   true,
		ListenPort:   "35053",
		PrivateKey:   "SDZnNuMWQz+cKlr6f7Vu+Q98R+sl1D9EJPmDWWJZaUM=",
		NetInterface: "eth",
		Peers: []cfgs.PeerConfig{
			{
				PublicKey:  "KQwNg8z7nSgD23nHga8PKeSrh2GupEstDkQ3Jww5eg4=",
				AllowedIPs: []string{"10.8.0.2/32"},
			},
			{
				PublicKey:  "KQwNg8z7nSgD23nHga8PKeSrh2GupEstDkQ3Jww5eg4=",
				AllowedIPs: []string{"10.8.0.3/32", "10.8.0.4/32"},
			},
		},
	}

	/// Processing client config
	fmt.Println("------ client cfg 1 ------")
	cfg, err := cfgs.ProcessClientConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cfg.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(os.Stdout, cfg); err != nil {
		panic(err)
	}

	fmt.Println("--------------------------")
	fmt.Println("------ client cfg 2 ------")
	clientConfig.PrivateKey = "aGsGuo9ODki0ZpS1U3c28tsI6UWjCW1Gbn8lIYRamXA="
	cfg, err = cfgs.ProcessClientConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cfg.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(os.Stdout, cfg); err != nil {
		panic(err)
	}

	fmt.Println("--------------------------")
	/// Processing server config
	fmt.Println("------ server cfg 0 ------")
	cfg, err = cfgs.ProcessServerConfig(serverConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cfg.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := io.Copy(os.Stdout, cfg); err != nil {
		panic(err)
	}
	fmt.Println("------------------------")
}
