package configs

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/joho/godotenv/autoload"
)

type ClientConfig struct {
	Address    string
	PrivateKey string
	DNS        []string

	PublicKey  string
	AllowedIPs []string
	Endpoint   string
}

type ServerConfig struct {
	Address      string
	SaveConfig   bool
	ListenPort   string
	PrivateKey   string
	NetInterface string

	Peers []PeerConfig
}

type PeerConfig struct {
	PublicKey  string
	AllowedIPs []string
}

const (
	clientTmplFile = "client.tmpl"
	serverTmplFile = "server.tmpl"
)

var (
	tmplFolder = os.Getenv("TEMPLATES_FOLDER")
	clientTmpl = must(
		template.New(clientTmplFile).
			Funcs(template.FuncMap{"join": strings.Join}).
			ParseFiles(filepath.Join(tmplFolder, clientTmplFile)),
	)
	serverTmpl = must(
		template.New(serverTmplFile).
			Funcs(template.FuncMap{"join": strings.Join}).
			ParseFiles(filepath.Join(tmplFolder, serverTmplFile)),
	)
)

func must(t *template.Template, err error) *template.Template {
	log.Printf("env variable: %s", os.Getenv("TEMPLATES_FOLDER"))
	log.Printf("template folder: %s", tmplFolder)
	if err != nil {
		panic(err)
	}
	return t
}

func ProcessClientConfig(cfg ClientConfig) (io.Reader, error) {
	return processConfig(cfg)
}

func ProcessServerConfig(cfg ServerConfig) (io.Reader, error) {
	return processConfig(cfg)
}

func processConfig(cfg interface{}) (io.Reader, error) {
	var err error
	pr, pw := io.Pipe()
	go func() {
		defer func() {
			if werr := pw.Close(); werr != nil {
				err = werr
			}
		}()
		switch c := cfg.(type) {
		case ClientConfig:
			err = clientTmpl.Execute(pw, c)
		case ServerConfig:
			err = serverTmpl.Execute(pw, c)
		default:
			err = errors.New("unsupported type of cfg argument")
		}
	}()
	return pr, err
}
