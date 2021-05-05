package telegram

import (
	"context"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"github.com/skoret/wireguard-bot/internal/wireguard"
)

type Bot struct {
	wg        *sync.WaitGroup
	api       *tgbotapi.BotAPI
	wireguard *wireguard.Wireguard
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

	wguard, err := wireguard.NewWireguard()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wireguard client")
	}

	return &Bot{
		wg:        &sync.WaitGroup{},
		api:       api,
		wireguard: wguard,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	// wait all running handlers to finish and close wg connection
	defer func() {
		b.wg.Wait()
		if err := b.wireguard.Close(); err != nil {
			log.Printf("failed to close wireguard connection: %v", err)
		}
	}()

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
				if errs := b.handle(&update); errs != nil {
					for _, err := range errs {
						log.Printf("error occured: %s", err.Error())
					}
				}
			}()
		case <-ctx.Done():
			log.Printf("stopping bot: %v", ctx.Err())
			b.api.StopReceivingUpdates()
			return nil
		}
	}
}

func (b *Bot) handle(update *tgbotapi.Update) []error {
	log.Printf("new update: %+v", update)
	var res []tgbotapi.Chattable
	var err error
	errs := make([]error, 0)
	switch {
	case update.Message != nil:
		res, err = b.handleMessage(update.Message)
	case update.CallbackQuery != nil:
		res, err = b.handleQuery(update.CallbackQuery)
	default:
		errs = append(errs, errors.New("unable to handle such update"))
	}
	if err != nil {
		errs = append(errs, err)
	}
	for _, resp := range res {
		if err := b.send(resp); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (b *Bot) send(c tgbotapi.Chattable) error {
	msg, err := b.api.Send(c)
	log.Printf("send msg: %+v", msg)
	if err != nil {
		return err
	}
	return nil
}
