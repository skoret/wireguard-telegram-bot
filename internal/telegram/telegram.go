package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
)

type Bot struct {
	api       *tgbotapi.BotAPI
	wg        *sync.WaitGroup
	wireguard struct{} // TODO
}

// NewBot creates new Bot instance
func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("bot user: %+v", api.Self)

	if err := setMyCommands(api); err != nil {
		return nil, err
	}

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
	updates := b.api.GetUpdatesChan(config)

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
	var (
		res tgbotapi.Chattable
		err error
	)
	switch {
	case update.Message != nil:
		res, err = b.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		res, err = b.handleQuery(update.CallbackQuery)
	default:
		return errors.New("unable to handle such update")
	}

	if res != nil {
		if err := b.send(res); err != nil {
			return err
		}
	}
	return err
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
