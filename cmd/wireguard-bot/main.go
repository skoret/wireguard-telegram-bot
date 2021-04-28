package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/skoret/wireguard-bot/internal/telegram"
	configs "github.com/skoret/wireguard-bot/internal/utils"
)

func main() {
	tg, err := telegram.NewBot(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("failed to create telegram bot: %s", err.Error())
	}

	// TODO: Insert templating procedure to apropriate function in telegram
	configs.Handle_client_config()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := tg.Run(ctx); err != nil {
			log.Fatalf("failed to run telegram bot: %s", err.Error())
		}
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
