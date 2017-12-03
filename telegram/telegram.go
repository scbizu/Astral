package telegram

import (
	"errors"
	"log"
	"sync"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/storage"
)

const (
	tokenKey = "tg_token"
)

var (
	errEmptyArgs = errors.New("plz input args")
	redisOnce    sync.Once
	token        string
)

func init() {
	var redisClient *redis.Client
	redisOnce.Do(func() {
		redisClient = storage.NewRedisClient()
	})
	var err error
	token, err = redisClient.Get(tokenKey).Result()
	if err != nil {
		log.Fatalln(err)
	}
}

func connectTG() (bot *tgbotapi.BotAPI, err error) {
	return tgbotapi.NewBotAPI(token)
}

//PullAndReply pull the msg (long polling) and response
func PullAndReply() (err error) {
	bot, err := connectTG()
	if err != nil {
		return
	}
	config := tgbotapi.NewUpdate(0)
	config.Timeout = 60

	updates, err := bot.GetUpdatesChan(config)
	if err != nil {
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == CommandHello {
			if update.Message.CommandArguments() == "" {
				return errEmptyArgs
			}
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.CommandArguments())
			reply.ReplyToMessageID = update.Message.MessageID
			bot.Send(reply)
		}
	}
	return
}
