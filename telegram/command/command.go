package command

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Handler defines handle func
type Handler func(msg *tgbotapi.Message) tgbotapi.MessageConfig

//CommanderName defines command's literal name
type CommanderName string

//Commander defines command obj
type Commander struct {
	Name     CommanderName
	Usage    string
	behavior Handler
}

const (
	//CommandSayhi == /sayhi
	CommandSayhi CommanderName = "sayhi"
	//CommandTodayAnime == /today_anime,combine all animes from all srcs.
	CommandTodayAnime CommanderName = "todayanime"
)

//DefaultBehavior defines the  default behavior of commander
var DefaultBehavior = func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	defaultText := fmt.Sprintf("%s 命中了,但是作者什么都没有实现哦😞", msg.Command())
	return tgbotapi.NewMessage(msg.Chat.ID, defaultText)
}

//NewCommand init a command
func NewCommand(name CommanderName, u string, handler Handler) *Commander {
	if handler == nil {

	}
	return &Commander{
		Name:     name,
		Usage:    u,
		behavior: handler,
	}
}

//Do will run the command behavior
func (c *Commander) Do(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	if msg.Command() != string(c.Name) {
		log.Printf("%s skiped command %s", msg.Command(), string(c.Name))
		return tgbotapi.MessageConfig{}
	}
	if c.behavior == nil {
		c.behavior = DefaultBehavior
	}
	return c.behavior(msg)
}
