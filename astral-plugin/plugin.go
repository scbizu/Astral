package plugin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/astral-plugin/lunch"
	"github.com/scbizu/Astral/astral-plugin/sayhi"
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
	return
}
