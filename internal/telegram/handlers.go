package telegram

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/yeqown/go-qrcode"
)

type responses []tgbotapi.Chattable

func (b *Bot) handleMessage(msg *tgbotapi.Message) (responses, error) {
	log.Printf("new message: %+v", msg)
	res0 := tgbotapi.NewMessage(msg.Chat.ID, "run /menu, silly")

	if !msg.IsCommand() {
		return responses{res0}, nil
	}
	cmd, ok := commands[msg.Command()]
	if !ok {
		return responses{res0}, errors.Errorf("message received with unknown command: %s", msg.Command())
	}
	res0.Text = cmd.text
	res0.ReplyMarkup = cmd.keyboard
	if cmd.handler == nil {
		return responses{res0}, nil
	}

	res1, err := cmd.handler(b, msg.Chat.ID, msg.CommandArguments())
	if err != nil {
		return responses{errorMessage(msg.Chat.ID, msg.MessageID, false)}, err
	}
	return append(responses{res0}, res1...), nil
}

func (b *Bot) handleQuery(query *tgbotapi.CallbackQuery) (responses, error) {
	log.Printf("new callback query: %+v", query)

	if query.Message == nil {
		return nil, errors.New("callback query received without message | it is possible only for inline mode")
	}
	log.Printf("message from callback: %+v", query.Message)
	chatID, msgID := query.Message.Chat.ID, query.Message.MessageID
	sorryMsg := errorMessage(chatID, msgID, true)

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		return responses{sorryMsg}, errors.Wrap(err, "failed to process callback query")
	}

	cmd, ok := commands[query.Data]
	if !ok {
		return responses{sorryMsg}, errors.Errorf("callback query received with unknown data field: %s", query.Data)
	}
	res0 := tgbotapi.NewEditMessageText(chatID, msgID, cmd.text)
	res0.ReplyMarkup = cmd.keyboard

	if cmd.handler == nil {
		return responses{res0}, nil
	}

	res1, err := cmd.handler(b, chatID, query.Message.CommandArguments())
	if err != nil {
		return responses{sorryMsg}, errors.Wrap(err, "unable to create new config")
	}
	return append(responses{res0}, res1...), nil
}

func (b *Bot) handleConfigForNewKeys(chadID int64, _ string) (responses, error) {
	cfg, err := b.wireguard.CreateConfigForNewKeys()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new config")
	}
	// create qrcode
	content, err := ioutil.ReadAll(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read new config")
	}
	options := []qrcode.ImageOption{
		qrcode.WithLogoImageFilePNG("assets/logo-min.png"),
		qrcode.WithQRWidth(7),
		qrcode.WithBuiltinImageEncoder(qrcode.PNG_FORMAT),
	}
	qrc, err := qrcode.New(string(content), options...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create qr code from config")
	}
	qr := bytes.Buffer{}
	if err := qrc.SaveTo(&qr); err != nil {
		return nil, errors.Wrap(err, "failed to read new qr code")
	}
	timestamp := time.Now().Unix()
	name := strconv.FormatInt(timestamp, 10)
	file0 := tgbotapi.NewPhoto(chadID, tgbotapi.FileReader{
		Name:   name + ".png",
		Reader: &qr,
	})

	file1 := tgbotapi.NewDocument(chadID, tgbotapi.FileBytes{
		Name:  name + ".conf",
		Bytes: content,
	})
	thumb, _ := os.Open("assets/logo-min.png")
	file1.Thumb = tgbotapi.FileReader{
		Name:   "thumb",
		Reader: thumb,
	}
	return responses{file0, file1}, nil
}

func (b *Bot) handleConfigForPublicKey(chadID int64, arg string) (responses, error) {
	if arg == "" {
		return nil, nil
	}
	cfg, err := b.wireguard.CreateConfigForPublicKey(arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new config")
	}
	timestamp := time.Now().Unix()
	file := tgbotapi.FileReader{
		Name:   fmt.Sprintf("wg-tg-%d.conf", timestamp),
		Reader: cfg,
	}
	return responses{tgbotapi.NewDocument(chadID, file)}, nil
}

func init() {
	ConfigForNewKeysCmd.handler = (*Bot).handleConfigForNewKeys
	ConfigForPublicKeyCmd.handler = (*Bot).handleConfigForPublicKey
}

const sorry = "something went wrong, sorry\n" +
	"or not\n" +
	"👉🏻👈🏻"

func errorMessage(chatID int64, msgID int, edit bool) (res tgbotapi.Chattable) {
	if edit {
		res = tgbotapi.NewEditMessageTextAndMarkup(
			chatID, msgID, sorry,
			tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(goToMenuButton),
			),
		)
	} else {
		res = tgbotapi.NewMessage(chatID, sorry)
	}
	return
}
