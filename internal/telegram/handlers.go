package telegram

import (
	"log"
	"net"
	"os"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	cfgs "github.com/skoret/wireguard-bot/internal/wireguard/configs"
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

	msg := query.Message
	if msg == nil {
		return nil, errors.New("callback query received without message | it is possible only for inline mode")
	}
	res := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, "something went wrong, try again later")
	log.Printf("message from callback: %+v", msg)

	callback := tgbotapi.NewCallback(query.ID, "")
	if _, err := b.api.Request(callback); err != nil {
		return responses{res}, errors.Wrap(err, "failed to process callback query")
	}

	cmd, ok := commands[query.Data]
	if !ok {
		return responses{res}, errors.Errorf("callback query received with unknown data field: %s", query.Data)
	}
	res.Text = cmd.text
	res.ReplyMarkup = cmd.keyboard
	if cmd.handler == nil {
		return responses{res}, nil
	}
	document, err := cmd.handler(b, msg.Chat.ID)
	if err != nil {
		return responses{res}, errors.Wrap(err, "unable to create new config")
	}
	return responses{res, document}, nil
}

func (b *Bot) handleConfigForNewKeys(chadID int64) (tgbotapi.Chattable, error) {
	pri, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}
	// TODO: get ip net dynamically
	address := "10.8.0.3/32"
	clientConfig := cfgs.ClientConfig{
		Address:    address,
		PrivateKey: pri.String(),
		DNS:        []string{"8.8.8.8", "8.8.4.4"},

		PublicKey:  os.Getenv("SERVER_PUB_KEY"),
		AllowedIPs: []string{"0.0.0.0/0"},
		Endpoint:   os.Getenv("SERVER_ENDPOINT"),
	}
	cfgFile, err := cfgs.ProcessClientConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// wg server conf update
	pub := pri.PublicKey()
	_, ipNet, err := net.ParseCIDR(address)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ip with mask")
	}

	cfg := wgtypes.Config{
		ReplacePeers: false,
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey:                   pub,
				Remove:                      false,
				UpdateOnly:                  false,
				PresharedKey:                nil,
				Endpoint:                    nil,
				PersistentKeepaliveInterval: nil,
				ReplaceAllowedIPs:           false,
				AllowedIPs:                  []net.IPNet{*ipNet},
			},
		},
	}

	// TODO: use wgctrl client from Bot.wireguard field
	c, err := wgctrl.New()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open wgctrl")
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Println(err)
		}
	}()
	if err := c.ConfigureDevice("wg0", cfg); err != nil {
		return nil, errors.Wrap(err, "failed to update server configuration")
	}

	cmd := exec.Command("wg-quick", "save", "wg0")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, errors.Wrap(err, "failed to dump server config to conf file")
	}

	file := tgbotapi.FileReader{
		Name:   "wg-tg-test.conf",
		Reader: cfgFile,
	}

	return tgbotapi.NewDocument(chadID, file), nil
}

func init() {
	ConfigForNewKeysCmd.handler = (*Bot).handleConfigForNewKeys
}
