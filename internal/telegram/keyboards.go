package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (cmd command) button() tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(cmd.Description, cmd.Command)
}

var (
	menuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(NewConfigCmd.button()),
		tgbotapi.NewInlineKeyboardRow(DonateCmd.button()),
		tgbotapi.NewInlineKeyboardRow(HelpCmd.button()),
	)

	goToMenuButton = tgbotapi.NewInlineKeyboardButtonData("go to menu", MenuCmd.Command)

	newConfigKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(ConfigForNewKeysCmd.button()),
		tgbotapi.NewInlineKeyboardRow(ConfigForPublicKeyCmd.button()),
		tgbotapi.NewInlineKeyboardRow(goToMenuButton),
	)

	configForPublicKeyKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"How to generate wireguard key pair?",
				"https://www.wireguard.com/quickstart/#key-generation",
			),
		),
		tgbotapi.NewInlineKeyboardRow(goToMenuButton),
	)

	donateKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(goToMenuButton),
	)

	helpKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(goToMenuButton),
	)
)

func init() {
	MenuCmd.keyboard = &menuKeyboard
	NewConfigCmd.keyboard = &newConfigKeyboard
	ConfigForPublicKeyCmd.keyboard = &configForPublicKeyKeyboard
	DonateCmd.keyboard = &donateKeyboard
	HelpCmd.keyboard = &helpKeyboard
}
