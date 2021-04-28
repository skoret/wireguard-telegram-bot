package configs

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type ClientConfig struct {
	Address    string
	PrivateKey string
	DNS        string

	PublicKey  string
	AllowedIPs string
	Endpoint   string
}

type ServerConfig struct {
	Address      string
	ListenPort   string
	PrivateKey   string
	NetInterface string

	PublicKey  string
	AllowedIPs string
}

func Handle_client_config() {

	clientConfig := ClientConfig{
		Address:    "10.8.0.2/24",
		PrivateKey: "aGsGuo9ODki0ZpS1U3c28tsI6UWjCW1Gbn8lIYRamXA=",
		DNS:        "8.8.8.8",

		PublicKey:  "G8naBU85RGmh2iZBi2KL3qomJOGKy5jvU97bO2I5tQ4=",
		AllowedIPs: "0.0.0.0/0",
		Endpoint:   "34.91.35.38:35053",
	}

	serverConfig := ServerConfig{
		Address:      "10.8.0.1/24",
		ListenPort:   "35053",
		PrivateKey:   "SDZnNuMWQz+cKlr6f7Vu+Q98R+sl1D9EJPmDWWJZaUM=",
		NetInterface: "eth",
		PublicKey:    "KQwNg8z7nSgD23nHga8PKeSrh2GupEstDkQ3Jww5eg4=",
		AllowedIPs:   "10.8.0.2/32",
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// TODO: Relative path to root dir of project
	parent := filepath.Dir(wd)
	grandParent := filepath.Dir(parent)

	/// Processing client config
	clientTemplatePath := grandParent + "/internal/utils/client.template"
	t, err := template.ParseFiles(clientTemplatePath)
	if err != nil {
		panic(err)
	}

	clientConfigFile, err := os.Create(grandParent + "/internal/utils/client.conf")
	if err != nil {
		log.Println("create file: ", err)
		panic(err)
	}

	err = t.Execute(clientConfigFile, clientConfig)
	if err != nil {
		panic(err)
	}

	clientConfigFile.Close()

	/// Processing server config
	serverTemplatePath := grandParent + "/internal/utils/server.template"
	t, err = template.ParseFiles(serverTemplatePath)
	if err != nil {
		panic(err)
	}

	serverConfigFile, err := os.Create(grandParent + "/internal/utils/server.conf")
	if err != nil {
		log.Println("create file: ", err)
		panic(err)
	}

	err = t.Execute(serverConfigFile, serverConfig)
	if err != nil {
		panic(err)
	}

	serverConfigFile.Close()
}
