package telegram

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type responses []tgbotapi.Chattable

func (b *Bot) handleMessage(msg *tgbotapi.Message) (responses, error) {
	log.Printf("new message: %+v", msg)
	res := tgbotapi.NewMessage(msg.Chat.ID, "run /menu, silly")

	if msg.IsCommand() {
		cmd, ok := commands[msg.Command()]
		if !ok {
			return responses{res}, errors.Errorf("message received with unknown command: %s", msg.Command())
		}
		res.Text = cmd.text
		res.ReplyMarkup = *cmd.keyboard
		// TODO: run some wireguard logic if needed
	}
	return responses{res}, nil
}

func (b *Bot) handleQuery(query *tgbotapi.CallbackQuery) (responses, error) {
	log.Printf("new callback query: %+v", query)

	if query.Message == nil {
		return nil, errors.New("callback query received without message | it is possible only for inline mode")
	}
	log.Printf("message from callback: %+v", query.Message)
	chatID, msgID := query.Message.Chat.ID, query.Message.MessageID
	sorry := errorMessage(chatID, msgID)

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		return responses{sorry}, errors.Wrap(err, "failed to process callback query")
	}

	cmd, ok := commands[query.Data]
	if !ok {
		return responses{sorry}, errors.Errorf("callback query received with unknown data field: %s", query.Data)
	}
	res := tgbotapi.NewEditMessageText(chatID, msgID, cmd.text)
	res.ReplyMarkup = cmd.keyboard

	if cmd.handler == nil {
		return responses{res}, nil
	}
	document, err := cmd.handler(b, chatID)
	if err != nil {
		return responses{sorry}, errors.Wrap(err, "unable to create new config")
	}
	return responses{res, document}, nil
}

func (b *Bot) handleConfigForNewKeys(chadID int64) (tgbotapi.Chattable, error) {
	cfg, err := b.wireguard.CreateNewConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new config")
	}
	timestamp := time.Now().Unix()
	file := tgbotapi.FileReader{
		Name:   fmt.Sprintf("wg-tg-%d.conf", timestamp),
		Reader: cfg,
	}
	return tgbotapi.NewDocument(chadID, file), nil
}

func init() {
	ConfigForNewKeysCmd.handler = (*Bot).handleConfigForNewKeys
}

func errorMessage(chatID int64, msgID int) tgbotapi.Chattable {
	return tgbotapi.NewEditMessageTextAndMarkup(
		chatID, msgID, "something went wrong, sorry 👉🏻👈🏻",
		tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(goToMenuButton),
		),
	)
}
