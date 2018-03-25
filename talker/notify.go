package talker

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifaction new the notifaction object
type Notifaction struct {
	channelName     string
	projectName     string
	isTestSucceed   bool
	isBuildSucceed  bool
	isDepolySucceed bool
	commitNotes     string
}

const (
	testStage   = "test"
	buildStage  = "build"
	deployStage = "deploy"
)

const (
	// ChannelChatID defines the bot name
	// Hacked this on the web client
	// Regex format: *c(.*)_*
	ChannelChatID = 1378084890
)

// NewNotifaction init the Notifaction instance
func NewNotifaction(repo string, stage map[string]bool, commit string) *Notifaction {
	if _, ok := stage[testStage]; !ok {
		return &Notifaction{}
	}
	if _, ok := stage[buildStage]; !ok {
		return &Notifaction{}
	}
	return &Notifaction{
		channelName:     "",
		projectName:     repo,
		isTestSucceed:   stage[testStage],
		isBuildSucceed:  stage[buildStage],
		isDepolySucceed: stage[deployStage],
		commitNotes:     commit,
	}
}

// Notify sends the msg to the tg channel
func (n *Notifaction) Notify() tgbotapi.MessageConfig {
	text := fmt.Sprintf("%s:%v", n.commitNotes, n.isBuildSucceed)
	return tgbotapi.NewMessage(ChannelChatID, text)
}
