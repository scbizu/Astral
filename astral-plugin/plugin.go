package plugin

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/astral-plugin/lunch"
	"github.com/scbizu/Astral/astral-plugin/py"
	"github.com/scbizu/Astral/astral-plugin/sayhi"
	"github.com/scbizu/Astral/astral-plugin/today-anime"
	"github.com/scbizu/wechat-go/wxweb"
)

//RegisterWechatEnabledPlugins regists wechat plugins
//to the main wx session.
func RegisterWechatEnabledPlugins(session *wxweb.Session) {
	// replier.Register(session, autoReply)
	lunch.Register(session, nil)
}

//RegistTGEnabledPlugins regists telegram plugin
func RegistTGEnabledPlugins(rawmsg *tgbotapi.Message) (msg tgbotapi.MessageConfig) {
	msg = sayhi.Register(rawmsg)
	if checkMarkedMsg(msg) {
		return
	}

	msg = anime.Register(rawmsg)
	if checkMarkedMsg(msg) {
		return
	}

	msg = py.Register(rawmsg)
	log.Println(msg.ChatID, msg.Text)
	return
}

func checkMarkedMsg(msg tgbotapi.MessageConfig) bool {
	log.Printf("[check chatid]:%d,[check msgText]:%s", msg.ChatID, msg.Text)
	if msg.ChatID != 0 && msg.Text != "" {
		return true
	}
	return false
}
