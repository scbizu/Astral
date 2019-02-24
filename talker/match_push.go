package talker

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	// CNSC2EventChannelName  defines cn sc2 channel
	CNSC2EventChannelName = "@CNSC2EventChannel"
)

type MatchPush struct {
	matches     []string
	channelName string
}

func NewMatchPush(matches []string) MatchPush {
	return MatchPush{
		matches:     matches,
		channelName: CNSC2EventChannelName,
	}
}

func (mp MatchPush) GetPushMessage() tgbotapi.MessageConfig {
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: mp.channelName,
		},
		Text:      strings.Join(mp.matches, "\n\n"),
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
