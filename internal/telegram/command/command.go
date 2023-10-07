package command

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Handler defines handle func
type Handler func(msg *tgbotapi.Message) tgbotapi.MessageConfig

// CommanderName defines command's literal name
type CommanderName string

// Commander defines command obj
type Commander struct {
	Name     CommanderName
	Usage    string
	behavior Handler
}

const (
	// CommandSayhi == /sayhi
	CommandSayhi CommanderName = "sayhi"
	// CommandTodayAnime == /today_anime,combine all animes from all srcs.
	CommandTodayAnime CommanderName = "todayanime"
	// CommandShowAllCommand makes py with @botfather
	CommandShowAllCommand CommanderName = "show_commands"
	// CommandAIChat makes bot chat with you
	CommandAIChat CommanderName = "chat"
)

var allCommandsMapping *sync.Map

// DefaultBehavior defines the  default behavior of commander
var DefaultBehavior = func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	defaultText := fmt.Sprintf("%s 命中了,但是作者什么都没有实现哦😞", msg.Command())
	return tgbotapi.NewMessage(msg.Chat.ID, defaultText)
}

func init() {
	allCommandsMapping = new(sync.Map)
}

// NewCommand init a command
func NewCommand(name CommanderName, u string, handler Handler) *Commander {
	if handler == nil {
		handler = DefaultBehavior
	}
	c := &Commander{
		Name:     name,
		Usage:    u,
		behavior: handler,
	}
	setCommand(c)

	return c
}

// String return commanderName
func (c CommanderName) String() string {
	return string(c)
}

// Do will run the command behavior
func (c *Commander) Do(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	if msg.Command() != string(c.Name) {
		logrus.Infof("%s skiped command %s", msg.Command(), string(c.Name))
		return tgbotapi.MessageConfig{}
	}
	return c.behavior(msg)
}

func setCommand(command *Commander) {
	allCommandsMapping.Store(command.Name, command)
}

// GetAllCommands gets all commands name
func GetAllCommands() (commands []*Commander) {
	allCommandsMapping.Range(func(_ interface{}, value interface{}) bool {
		if _, ok := value.(*Commander); !ok {
			return false
		}
		commands = append(commands, value.(*Commander))
		return true
	})
	return
}
