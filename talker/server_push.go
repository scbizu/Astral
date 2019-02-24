package talker

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifaction new the notifaction object
type Notifaction struct {
	channelName  string
	projectName  string
	testStatus   string
	buildStatus  string
	depolyStatus string
	commitNotes  string
	duration     int64
	serverStatus string
}

const (
	testStage   = "test"
	buildStage  = "build"
	deployStage = "deploy"
)

const (
	// ChannelName defines build/deploy channel name
	ChannelName = "@AstralServerNotification"
)

// NewNotifaction init the Notifaction instance
func NewNotifaction(repo string, stage map[string]string, commit string, duration int64) *Notifaction {
	if _, ok := stage[testStage]; !ok {
		return &Notifaction{}
	}
	if _, ok := stage[buildStage]; !ok {
		return &Notifaction{}
	}
	testStatus := stage[testStage]
	buildStatus := stage[buildStage]
	deployStatus := stage[deployStage]

	return &Notifaction{
		channelName:  ChannelName,
		projectName:  repo,
		testStatus:   testStatus,
		buildStatus:  buildStatus,
		depolyStatus: deployStatus,
		commitNotes:  commit,
		duration:     duration,
	}
}

// Notify sends the msg to the tg channel
func (n *Notifaction) Notify() tgbotapi.MessageConfig {
	text := fmt.Sprintf("[Commit Note]: `%s` ", n.commitNotes)
	text = fmt.Sprintf("%s\n **Test Status**: `%v`", text, n.testStatus)
	text = fmt.Sprintf("%s\n **Build Status**: `%v`", text, n.buildStatus)
	text = fmt.Sprintf("%s\n **Deploy Status**: `%v`", text, n.depolyStatus)
	text = fmt.Sprintf("%s\n **Duration**: `%v s`", text, n.duration)
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: n.channelName,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}
}

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
			ChannelUsername: n.channelName,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
