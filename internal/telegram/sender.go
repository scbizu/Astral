package telegram

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/scbizu/Astral/internal/tl"
)

const (
	// CNSC2EventChannelName  defines cn sc2 channel
	CNSC2EventChannelName = "@CNSC2EventChannel"
)

func NewTelegram(bot *Bot) *Telegram {
	return &Telegram{bot}
}

type Telegram struct {
	*Bot
}

func (ts *Telegram) Send(msg string, fs ...tl.Filter) error {
	for _, f := range fs {
		msg = f.F(msg)
	}
	if msg == "" {
		return nil
	}
	msgConfig := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: CNSC2EventChannelName,
		},
		Text:      msg,
		ParseMode: tgbotapi.ModeMarkdown,
	}
	if _, err := ts.bot.Send(msgConfig); err != nil {
		return err
	}
	return nil
}

func (ts *Telegram) SendAndReturnID(msg string, fs ...tl.Filter) (string, error) {
	for _, f := range fs {
		msg = f.F(msg)
	}
	if msg == "" {
		return "", nil
	}
	msgConfig := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: CNSC2EventChannelName,
		},
		Text:      msg,
		ParseMode: tgbotapi.ModeMarkdown,
	}
	resp, err := ts.bot.Send(msgConfig)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(resp.MessageID), nil
}

func (ts *Telegram) ResolveMessage(msgs []string) string {
	return strings.Join(msgs, "\n\n")
}
