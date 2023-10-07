// Package sayhi is the telegram plugin
package sayhi

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
)

const (
	masterName = "scnace"
)

var _ plugin.IPlugin = (*Handler)(nil)

// Handler impl the PluginHandler
type Handler struct{}

func (h *Handler) Enable() bool {
	return true
}

// Register regists sayhi plugin
func (h *Handler) Process(msg *tgbotapi.Message) tgbotapi.MessageConfig {
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

func (h *Handler) Name() command.CommanderName {
	return command.CommandSayhi
}
