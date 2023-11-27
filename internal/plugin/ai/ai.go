package ai

import (
	"bytes"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rakyll/openai-go/chat"
	"github.com/scbizu/Astral/internal/openai"
	"github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
)

var _ plugin.IPlugin = (*AICommands)(nil)

type AICommands struct{}

func (ai *AICommands) Name() command.CommanderName {
	return command.CommandAIChat
}

func (ai *AICommands) Enable() bool {
	return true
}

func (ai *AICommands) Process(msg *tgbotapi.Message) tgbotapi.Chattable {
	cmd := command.NewCommand(
		command.CommandAIChat,
		"chat with openAI",
		func(msg *tgbotapi.Message) tgbotapi.Chattable {
			resp, err := openai.GetOpenAIClient().CreateCompletion(context.TODO(), &chat.CreateCompletionParams{
				Messages: []*chat.Message{
					{Role: "user", Content: msg.Text},
				},
			})
			if err != nil {
				return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			}
			msgBuffer := bytes.NewBuffer(nil)
			for _, choice := range resp.Choices {
				fmt.Fprintf(msgBuffer, "%s\n", choice.Message.Content)
			}
			return tgbotapi.NewMessage(msg.Chat.ID, msgBuffer.String())
		},
	)
	return cmd.Do(msg)
}
