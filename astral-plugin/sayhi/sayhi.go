//Package sayhi is the telegram plugin
package sayhi

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/telegram/command"
)

const (
	masterName = "scnace"
)

//Register regists sayhi plugin
func Register(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	sayHandler := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
		user := msg.From.UserName
		if user == masterName {
			user = "master"
		}
		return tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("hi,%s", user))
	}
	say := command.NewCommand(command.CommandSayhi, "say hi to every one", sayHandler)
	return say.Do(msg)
}
