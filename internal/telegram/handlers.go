package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
)

type responses []tgbotapi.Chattable

func (b *Bot) handleMessage(msg *tgbotapi.Message) (responses, error) {
	log.Printf("new message: %+v", msg)
	res := tgbotapi.NewMessage(msg.Chat.ID, "run /menu, silly")

	if msg.IsCommand() {
		cmd, ok := commands[msg.Command()]
		if !ok {
			return responses{res}, fmt.Errorf("message received with unknown command: %s", msg.Command())
		}
		res.Text = cmd.text
		res.ReplyMarkup = *cmd.keyboard
		// TODO: run some wireguard logic if needed
	}
	return responses{res}, nil
}

func (b Bot) handleQuery(query *tgbotapi.CallbackQuery) (responses, error) {
	log.Printf("new callback query: %+v", query)

	msg := query.Message
	if msg == nil {
		return nil, errors.New("callback query received without message | it is possible only for inline mode")
	}
	res := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "something went wrong, try again later")
	log.Printf("message from callback: %+v", msg)

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		return responses{res}, err
	}

	cmd, ok := commands[query.Data]
	if !ok {
		return responses{res}, fmt.Errorf("callback query received with unknown data field: %s", query.Data)
	}
	res.Text = cmd.text
	res.ReplyMarkup = cmd.keyboard
	// TODO: run some wireguard logic if needed
	if cmd.handler == nil {
		return responses{res}, nil
	}
	cfg, ok := cmd.handler(struct{}{}).(io.Reader)
	if !ok {
		return responses{res}, errors.New("uuupsiiiii")
	}
	file := tgbotapi.FileReader{
		Name:   "wg-tg-test.conf",
		Reader: cfg,
	}

	document := tgbotapi.NewDocument(msg.Chat.ID, file)
	return responses{res, document}, nil
}
