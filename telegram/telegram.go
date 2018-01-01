package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/getcert"
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

//PullAndReply pull the msg (long polling) and response
func PullAndReply() (err error) {
	bot, err := connectTG()
	if err != nil {
		return
	}
	bot.Debug = true
	log.Printf("bot auth passed as %s", bot.Self.UserName)

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

//ListenWebHook is the tg api webhook mode
func ListenWebHook(debug bool) (err error) {
	bot, err := connectTG()
	if err != nil {
		return
	}
	if debug {
		bot.Debug = true
		log.Printf("bot auth passed as %s", bot.Self.UserName)
	}
	bot.RemoveWebhook()
	cert := getcert.NewDomainCert(tgAPIDomain)
	domainWithToken := fmt.Sprintf("%s%s", cert.GetDomain(), token)
	if _, err = bot.SetWebhook(tgbotapi.NewWebhook(domainWithToken)); err != nil {
		log.Printf("notify webhook failed:%s", err.Error())
		return
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Panicln(err)
	}
	log.Println(info.LastErrorMessage, info.LastErrorDate)

	pattern := fmt.Sprintf("/%s", token)
	updatesMsgChannel := bot.ListenForWebhook(pattern)
	log.Printf("msg in channel:%d", len(updatesMsgChannel))

	port := fmt.Sprintf(":%s", os.Getenv("LISTENPORT"))

	go http.ListenAndServe(port, nil)

	for update := range updatesMsgChannel {
		log.Printf("[raw msg]:%+v\n", update)

		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}

	return
}
