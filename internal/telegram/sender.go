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

func NewTGSender(bot *Bot) *TGSender {
	return &TGSender{bot}
}

type TGSender struct {
	*Bot
}

func (ts *TGSender) Send(msg string) error {
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

func (ts *TGSender) SendAndReturnID(msg string) (string, error) {
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

func (ts *TGSender) ResolveMessage(msgs []string) string {
	return strings.Join(msgs, "\n\n")
}
