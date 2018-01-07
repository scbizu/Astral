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
	tokenKey    = "tg_token"
	tgAPIDomain = "https://scnace.cc:443/"
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
