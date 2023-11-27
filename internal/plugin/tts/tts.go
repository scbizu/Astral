package tts

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
	itts "github.com/scbizu/Astral/internal/tts"
)

var _ plugin.IPlugin = (*TTSCommand)(nil)

type TTSCommand struct{}

func (t *TTSCommand) Name() command.CommanderName {
	return command.CommandTTS
}

func (t *TTSCommand) Enable() bool {
	return true
}

func (t *TTSCommand) Process(msg *tgbotapi.Message) tgbotapi.Chattable {
	cmd := command.NewCommand(
		command.CommandTTS,
		"tts",
		func(msg *tgbotapi.Message) tgbotapi.Chattable {
			c := itts.NewElevenLabClient()
			bs, err := c.ToSpeech(context.Background(), msg.Text)
			if err != nil {
				return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			}
			return tgbotapi.NewVoice(msg.Chat.ID, tgbotapi.FileBytes{
				Name:  fmt.Sprintf("%s.mp3", msg.Text),
				Bytes: bs,
			})
		},
	)
	return cmd.Do(msg)
}
