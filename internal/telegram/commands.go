package telegram

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler func(b *Bot, chatID int64) (tgbotapi.Chattable, error)

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
	NewConfigCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "newconfing",
			Description: "create new config file for public server",
		},
		text: "do you want new config for new generated key pair or for your public key?",
	}
	ConfigForNewKeysCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "generatekeys",
			Description: "create new config file for new generated key pair",
		},
		text: "this is your new config for public wireguard vpn server, keep it in secret!",
	}
	ConfigForPublicKeyCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "configforkey",
			Description: "create new config file for given public key",
		},
		text: "send me your wireguard public key, please",
	}
	DonateCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "donate",
			Description: "buy me a beer and get a private wg vpn [WIP]",
		},
		text: "sorry, donations aren't supported yet",
	}
	HelpCmd = command{
		BotCommand: tgbotapi.BotCommand{
			Command:     "help",
			Description: "show bot's functionality with description",
		},
		text: "hi, i'm wireguard bot\n\n" +
			"i can create new wg vpn configuration files for you\n" +
			// TODO: write proper help message
			"TODO: write proper help message",
	}
)

var commands = map[string]*command{
	MenuCmd.Command:               &MenuCmd,
	NewConfigCmd.Command:          &NewConfigCmd,
	ConfigForNewKeysCmd.Command:   &ConfigForNewKeysCmd,
	ConfigForPublicKeyCmd.Command: &ConfigForPublicKeyCmd,
	ConfigForPublicKeyCmd.Command: &ConfigForPublicKeyCmd,
	DonateCmd.Command:             &DonateCmd,
	HelpCmd.Command:               &HelpCmd,
}

// setMyCommands is adapted method from unreleased v5.0.1
// https://github.com/go-telegram-bot-api/telegram-bot-api/commit/4a2c8c4547a868841c1ec088302b23b59443de2b
func setMyCommands(api *tgbotapi.BotAPI) error {
	params := make(tgbotapi.Params)
	data, err := json.Marshal([]command{MenuCmd, NewConfigCmd, DonateCmd, HelpCmd})
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
