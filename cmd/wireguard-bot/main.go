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

	done := make(chan struct{})
	go func() {
		if err := tg.Run(ctx); err != nil {
			log.Fatalf("error running telegram: %s", err.Error())
		}
		close(done)
	}()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Printf("graceful shutdown with signal %v", sig)
		cancel()
		<-done
	}()
	<-done
}
