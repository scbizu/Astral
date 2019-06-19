package talker

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifaction new the notifaction object
type Notifaction struct {
	channelName  string
	serverStatus string
}

const (
	// ChannelName defines build/deploy channel name
	ChannelName = "@AstralServerNotification"
)

// NewServerStatusNotifaction new server status
func NewServerStatusNotifaction(serverStatus string) *Notifaction {
	return &Notifaction{
		channelName:  ChannelName,
		serverStatus: serverStatus,
	}
}

// ServerStatusMsg builds the server status message
func (n *Notifaction) ServerStatusMsg() tgbotapi.MessageConfig {
	t := `
	[Astral Server Status]\n
	Astral has some problems,Please Login And Check.\n
	Astral Server Event: **%s**,\n
	Astral Server Event Timestamp(Server Time): %d\n
	`

	text := fmt.Sprintf(t, n.serverStatus, time.Now().Unix())

	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: ChannelName,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
