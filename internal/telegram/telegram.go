package telegram

import (
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	tokenKey    = "BOTKEY"
	tgAPIDomain = "https://astral-on-telegram.fly.dev"
)

var (
	errEmptyArgs = errors.New("plz input args")
	token        string
)

func init() {
	token = os.Getenv(tokenKey)
}

// ConnectTG returns the bot instance
func ConnectTG() (bot *tgbotapi.BotAPI, err error) {
	return tgbotapi.NewBotAPI(token)
}
