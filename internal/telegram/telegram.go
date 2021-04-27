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
	// TODO: unify commands/callbacks handling
	if update.Message != nil {
		log.Printf("new message: %+v", update.Message)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "run /menu, silly")

		if update.Message.IsCommand() {
			cmd, ok := commands[update.Message.Command()]
			if ok {
				msg.Text = cmd.text
				msg.ReplyMarkup = *cmd.keyboard
				// TODO: run some wireguard logic if needed
			} else {
				log.Printf("message received with unknown command: %s", update.Message.Command())
			}
		}

		if err := b.send(msg); err != nil {
			return err
		}
	} else if update.CallbackQuery != nil {
		query := update.CallbackQuery
		log.Printf("new callback query: %+v", query)
		callback := tgbotapi.NewCallback(query.ID, "")
		if _, err := b.api.Request(callback); err != nil {
			return err
		}

		if query.Message == nil {
			return errors.New("callback query received without message | it is possible only for inline mode")
		}
		log.Printf("message from callback: %+v", query.Message)

		msg := tgbotapi.NewEditMessageText(query.Message.Chat.ID, query.Message.MessageID, "something went wrong, try again later")
		cmd, ok := commands[query.Data]
		if ok {
			msg.Text = cmd.text
			msg.ReplyMarkup = cmd.keyboard
			// TODO: run some wireguard logic if needed
		} else {
			log.Printf("callback query received with unknown data field: %s", query.Data)
		}
		if err := b.send(msg); err != nil {
			return err
		}
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
