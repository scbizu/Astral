package telegram

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func (ts *Telegram) Send(msg string) error {
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

func (ts *Telegram) SendAndReturnID(msg string) (string, error) {
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
