package telegram

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	plugin "github.com/scbizu/Astral/astral-plugin"
	"github.com/scbizu/Astral/getcert"
)

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
		var msg tgbotapi.MessageConfig
		msg = plugin.RegistTGEnabledPlugins(update.Message)

		if msg.Text == "" || msg.ChatID == 0 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		}
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}

	return
}
