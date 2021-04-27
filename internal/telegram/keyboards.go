package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	menuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(NewConfigCmd.button()),
		tgbotapi.NewInlineKeyboardRow(DonateCmd.button()),
		tgbotapi.NewInlineKeyboardRow(HelpCmd.button()),
	)

	backToMenuButton = tgbotapi.NewInlineKeyboardButtonData("<< back to menu", MenuCmd.Command)

	newConfigKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(ConfigForNewKeysCmd.button()),
		tgbotapi.NewInlineKeyboardRow(ConfigForPublicKeyCmd.button()),
		tgbotapi.NewInlineKeyboardRow(backToMenuButton),
	)

	configForPublicKeyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"How to generate wireguard key pair?",
				"https://www.wireguard.com/quickstart/#key-generation",
			),
		),
		tgbotapi.NewInlineKeyboardRow(backToMenuButton),
	)

	donateKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backToMenuButton),
	)

	helpKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(backToMenuButton),
	)
)

func init() {
	MenuCmd.keyboard = &menuKeyboard
	NewConfigCmd.keyboard = &newConfigKeyboard
	ConfigForPublicKeyCmd.keyboard = &configForPublicKeyKeyboard
	DonateCmd.keyboard = &donateKeyboard
	HelpCmd.keyboard = &helpKeyboard
}
