package telegram

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler func(b *Bot, chatID int64, arg string) (responses, error)

type command struct {
	tgbotapi.BotCommand
	text     string
	keyboard *tgbotapi.InlineKeyboardMarkup
	handler  handler
}

var (
	MenuCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "menu",
			Description: "bot menu",
		},
		text: "so, what do you want?",
	}
	ConfigForNewKeysCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "newkeys",
			Description: "new config for new key pair",
		},
		text: "",
	}
	ConfigForPublicKeyCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "pubkey",
			Description: "new config for your public key",
		},
		text: "send me your wireguard public key, like that:\n" +
			"`/pubkey <your key in base64>`",
	}
	HelpCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "help",
			Description: "show bot's functionality with description",
		},
		text: "hi, i'm wireguard telegram bot\n" +
			"i can create new wireguard vpn configuration files for you\n\n" +
			"/menu — available commands\n" +
			"/newkeys — new config for new key pair\n" +
			"/pubkey — new config for your public key\n" +
			"/help — this message",
	}
)

var commands = map[string]*command{
	MenuCmd.Command:               &MenuCmd,
	ConfigForNewKeysCmd.Command:   &ConfigForNewKeysCmd,
	ConfigForPublicKeyCmd.Command: &ConfigForPublicKeyCmd,
	HelpCmd.Command:               &HelpCmd,
}

// setMyCommands is adapted method from unreleased v5.0.1
// https://github.com/go-telegram-bot-api/telegram-bot-api/commit/4a2c8c4547a868841c1ec088302b23b59443de2b
func setMyCommands(api *tgbotapi.BotAPI) error {
	params := make(tgbotapi.Params)
	data, err := json.Marshal([]command{
		MenuCmd,
		ConfigForNewKeysCmd,
		ConfigForPublicKeyCmd,
		HelpCmd,
	})
	if err != nil {
		return err
	}
	params.AddNonEmpty("commands", string(data))
	_, err = api.MakeRequest("setMyCommands", params)
	if err != nil {
		return err
	}
	return nil
}
