package main

import (
	"context"
	"github.com/skoret/wireguard-bot/internal/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: use some config library to load token from env file
	tg, err := telegram.NewBot(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		if err := tg.Run(ctx); err != nil {
			log.Fatalf("error running telegram: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("graceful shutdown")
}
