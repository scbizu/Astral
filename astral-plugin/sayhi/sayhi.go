//Package sayhi is the telegram plugin
package sayhi

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	plugin "github.com/scbizu/Astral/astral-plugin"
	"github.com/scbizu/Astral/telegram/command"
)

const (
	masterName = "scnace"
)

// Handler impl the PluginHandler
type Handler struct{}

//Register regists sayhi plugin
func (h *Handler) Register(msg *tgbotapi.Message) func(*tgbotapi.Message) tgbotapi.MessageConfig {

	sayhiRegister := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
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

	return sayhiRegister

}

// Sayhi returns sayhi plugin
func Sayhi(msg *tgbotapi.Message) *plugin.TGPlugin {
	return plugin.NewTGPlugin(command.CommandSayhi.String(), msg, plugin.Handler(&Handler{}))
}
