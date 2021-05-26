# wireguard-telegram-bot

<img alt="It was supposed to be a cool logo here" src="https://github.com/skoret/wireguard-telegram-bot/raw/dev/assets/logo.png" width="256" height="256" align="right">

Simple-Dimple Telegram Bot for Wireguard VPN config generation

## Functionality

- `/menu` — list available commands
- `/newkeys` — create a new config file and qr code for new generated key pair
- `/pubkey` — create a new config file template for the public key you provided
- `/help` — print this message

## Public Wireguard Telegram Bot

Just drop a message to [@wrgrdtgbot](https://t.me/wrgrdtgbot) and ask him for some new config for you and your friends   
[Install](https://www.wireguard.com/install/) Wireguard client for your device and import generated file or scan qr code

<p align="center">
  <img alt="bot screen 1" src="https://github.com/skoret/wireguard-telegram-bot/raw/dev/assets/bot/bot_1.png" width="300" />
  <img alt="bot screen 2" src="https://github.com/skoret/wireguard-telegram-bot/raw/dev/assets/bot/bot_2.png" width="300" /> 
</p>
<p align="center">
  <img alt="bot screen 3" src="https://github.com/skoret/wireguard-telegram-bot/raw/dev/assets/bot/bot_3.png" width="300" />
  <img alt="bot screen 4" src="https://github.com/skoret/wireguard-telegram-bot/raw/dev/assets/bot/bot_4.png" width="300" />
</p>

> **Disclaimer:** stability, availability and security **are not** guaranteed! Sorry not sorry 👉🏻👈🏻

## Setup your own Wireguard Telegram Bot

- Go to [@BotFather](https://t.me/BotFather), send him `/newbot`, choose a bot's name and username, and receive Telegram Bot API Token
- Go to AWS, GCP, whatever ☁️ and setup your remote server in desired region
  - You need to open corresponding port (e.g. `udp:51820`)
- Install `go`, `wireguard` and `wireguard-tools` on your server
  - Someday, we hope there will be a handy Dockerfile for it 🐳
- Generate Wireguard key pair for your server, create appropriate config file (e.g. `wg0.conf`) and run Wireguard
  - You're all big boys, you'll handle it
- `git clone git@github.com:skoret/wireguard-telegram-bot.git`
- `cd wireguard-telegram-bot && cp .env.example .env`
- Set env variables in `.env` file:

  | Variable              | Content | Notes |
    | :-------------------- | ------- | ----- |
  | `TELEGRAM_APITOKEN`   | your Telegram Bot API token from [@BotFather](https://t.me/BotFather) | keep it in _secret_! |
  | `ADMIN_USERNAMES`     | list of Telegram usernames, separated by commas, who are allowed to access this bot | leave variable _empty_ for public access |
  | `DNS_IPS`             | list of DNS ip addresses, separated by commas | e.g. `8.8.8.8,8.8.4.4` |
  | `SERVER_ENDPOINT`     | `<your_machine's_external_ip:open_port>` | |
  | `WIREGUARD_INTERFACE` | new Wireguard interface name | e.g. `wg0` |
  | `TEMPLATES_FOLDER`    | path to configuration template files | probably, you don't wanna change it |
  | `DEV_MODE`            | `false` for common uses<br />`true` for mocked internal wireguard client | dev mode suitable for manual bot ui tests |
- `sudo go run cmd/bot/main.go`
- 🎉 🍻 🥳

---
We hope the bot will be helpful. The code is not of the best quality. Contributions are welcome!

---
### Acknowledgements
- Thanks to @randallmunroe and all [ipython/xkcd-font](https://github.com/ipython/xkcd-font) contributors for such an awesome stuff
- Thanks to authors from [Noun Project](https://thenounproject.com/)
  - [Arrow by Andre](https://thenounproject.com/icon/1771844/)
  - [dragon by P Thanga Vignesh](https://thenounproject.com/icon/2863783/)
  - [Telegram by Danil Polshin](https://thenounproject.com/icon/1634539/)
