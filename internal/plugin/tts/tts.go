package tts

import (
	"bytes"
	"context"
	"fmt"
	"strings"

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
			text := strings.TrimPrefix(msg.Text, "/tts")
			bs, err := c.ToSpeech(context.Background(), text)
			if err != nil {
				return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			}
			return tgbotapi.NewVoice(msg.Chat.ID, tgbotapi.FileReader{
				Name:   fmt.Sprintf("%s.ogg", msg.Text),
				Reader: bytes.NewBuffer(bs),
			})
		},
	)
	return cmd.Do(msg)
}
