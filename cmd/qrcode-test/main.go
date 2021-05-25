package main

import (
	"fmt"
	"os"

	"github.com/yeqown/go-qrcode"
)

func main() {
	conf, err := os.ReadFile("testdata/client.conf")
	if err != nil {
		fmt.Printf("could not read file: %v", err)
	}
	options := []qrcode.ImageOption{
		qrcode.WithLogoImageFilePNG("assets/logo-min.png"),
		qrcode.WithQRWidth(7),
		qrcode.WithBuiltinImageEncoder(qrcode.PNG_FORMAT),
	}
	qrc, err := qrcode.New(string(conf), options...)

	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	if err := qrc.Save("testdata/conf-qrcode.png"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
