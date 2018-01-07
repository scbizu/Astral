package telegram

import (
	"errors"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	tokenKey    = "BOTKEY"
	tgAPIDomain = "https://scnace.cc:443/"
)

var (
	errEmptyArgs = errors.New("plz input args")
	token        string
)

func init() {
	token = os.Getenv(tokenKey)
}

func connectTG() (bot *tgbotapi.BotAPI, err error) {
	return tgbotapi.NewBotAPI(token)
}
