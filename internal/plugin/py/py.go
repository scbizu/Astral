package py

import (
	"bytes"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	plugin "github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
)

var _ plugin.IPlugin = (*Handler)(nil)

// Handler impl the PluginHandler
type Handler struct{}

// Register regists py plugin
func (h *Handler) Process(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	pyHandler := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
		pyStr := format(command.GetAllCommands())
		return tgbotapi.NewMessage(msg.Chat.ID, pyStr)
	}

	pyCommand := command.NewCommand(
		command.CommandShowAllCommand,
		"make py with @botfather",
		pyHandler,
	)

	return pyCommand.Do(msg)
}

func (h *Handler) Name() command.CommanderName {
	return command.CommandShowAllCommand
}

func (h *Handler) Enable() bool {
	return true
}

func format(commands []*command.Commander) (formatedStr string) {
	res := bytes.NewBufferString("")
	for idx, c := range commands {
		// the plainCommandStr satisfied botfather's suggested format
		plainCommandStr := fmt.Sprintf("%s - %s", c.Name, c.Usage)
		res.WriteString(plainCommandStr)
		if idx != len(commands)-1 {
			res.WriteString("\n")
		}
	}
	return res.String()
}
