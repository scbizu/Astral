package talker

import (
	"fmt"

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
	// ChannelName defines build/deploy channel name
	ChannelName = "@AstralServerNotifaction"
)

// NewNotifaction init the Notifaction instance
func NewNotifaction(repo string, stage map[string]string, commit string) *Notifaction {
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
		channelName:  "",
		projectName:  repo,
		testStatus:   testStatus,
		buildStatus:  buildStatus,
		depolyStatus: deployStatus,
		commitNotes:  commit,
	}
}

// Notify sends the msg to the tg channel
func (n *Notifaction) Notify() tgbotapi.MessageConfig {
	text := fmt.Sprintf("**Commit Note**: `%s` ", n.commitNotes)
	text = fmt.Sprintf("%s\n **Test Status**: `%v`", text, n.testStatus)
	text = fmt.Sprintf("%s\n **Build Status**: `%v`", text, n.buildStatus)
	text = fmt.Sprintf("%s\n **Deploy Status**: `%v`", text, n.depolyStatus)
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: ChannelName,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
