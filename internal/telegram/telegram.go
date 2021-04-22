package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	wg        *sync.WaitGroup
	wireguard struct{} // TODO
}

// NewBot creates new Bot instance
// TODO: NewBot should register available bot commands also
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("bot user: %+v", api.Self)
	return &Bot{
		api: api,
		wg:  &sync.WaitGroup{},
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	// wait all running handlers to finish
	defer b.wg.Wait()

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 30

	// Start polling Telegram for updates
	// TODO: someday it should be webhook instead of updates pulling
	updates, err := b.api.GetUpdatesChan(config)
	if err != nil {
		return err
	}

	for {
		select {
		case update := <-updates:
			b.wg.Add(1)
			go func() {
				defer b.wg.Done()
				if err := b.handle(&update); err != nil {
					log.Printf("uups, it's error: %s", err.Error())
				}
			}()
		case <-ctx.Done():
			log.Printf("stopping bot: %v", ctx.Err())
			b.api.StopReceivingUpdates()
			return nil
		}
	}
}

// TODO: handle different commands from user
func (b *Bot) handle(update *tgbotapi.Update) error {
	log.Printf("new update: %+v", update)
	if update.Message == nil {
		return nil
	}
	if update.Message.Text == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "meeeeh, send me smth")
		if err := b.send(msg); err != nil {
			return err
		}
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	if err := b.send(msg); err != nil {
		return err
	}
	return nil
}

// TODO: send and files too
func (b *Bot) send(c tgbotapi.Chattable) error {
	msg, err := b.api.Send(c)
	log.Printf("send msg: %+v", msg)
	if err != nil {
		return err
	}
	return nil
}
