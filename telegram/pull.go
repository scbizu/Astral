package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		reply.ReplyToMessageID = update.Message.MessageID
		bot.Send(reply)
	}

	return
}
