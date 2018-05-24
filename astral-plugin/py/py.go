package py

import (
	"bytes"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/telegram/command"
)

// Handler impl the PluginHandler
type Handler struct{}

// Register regists py plugin
func (h *Handler) Register(msg *tgbotapi.Message) func(*tgbotapi.Message) tgbotapi.MessageConfig {
	pyRegister := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
		pyHandler := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
			pyStr := format(command.GetAllCommands())
			return tgbotapi.NewMessage(msg.Chat.ID, pyStr)
		}
		pyCommand := command.NewCommand(command.CommandShowAllCommand, "make py with @botfather", pyHandler)
		return pyCommand.Do(msg)
	}
	return pyRegister
}

func format(commands []*command.Commander) (formatedStr string) {
	res := bytes.NewBufferString("")
	for idx, c := range commands {
		//the plainCommandStr satisfied botfather's suggested format
		plainCommandStr := fmt.Sprintf("%s - %s", c.Name, c.Usage)
		res.WriteString(plainCommandStr)
		if idx != len(commands)-1 {
			res.WriteString("\n")
		}
	}
	return res.String()
}
