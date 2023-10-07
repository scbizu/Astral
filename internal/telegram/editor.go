package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (ts *Telegram) Edit(msgID string, content string) error {
	iMsgID, err := strconv.ParseInt(msgID, 10, 64)
	if err != nil {
		return err
	}
	msgConfig := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:          iMsgID,
			ChannelUsername: CNSC2EventChannelName,
		},
		Text:      content,
		ParseMode: tgbotapi.ModeMarkdown,
	}
	if _, err := ts.bot.Send(msgConfig); err != nil {
		return err
	}
	return nil
}
