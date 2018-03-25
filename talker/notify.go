package talker

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Notifaction new the notifaction object
type Notifaction struct {
	channelName     string
	projectName     string
	isTestSucceed   string
	isBuildSucceed  string
	isDepolySucceed string
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
	// ChannelName defines build/deploy channel name
	ChannelName = "@AstralServerNotifaction"

	successIcon = "✅"

	failedIcon = "❌"
)

func convert2Icon(isSuccess bool) string {
	if isSuccess {
		return successIcon
	}
	return failedIcon
}

// NewNotifaction init the Notifaction instance
func NewNotifaction(repo string, stage map[string]bool, commit string) *Notifaction {
	if _, ok := stage[testStage]; !ok {
		return &Notifaction{}
	}
	if _, ok := stage[buildStage]; !ok {
		return &Notifaction{}
	}
	testStatus := convert2Icon(stage[testStage])
	buildStatus := convert2Icon(stage[buildStage])
	deployStatus := convert2Icon(stage[deployStage])

	return &Notifaction{
		channelName:     "",
		projectName:     repo,
		isTestSucceed:   testStatus,
		isBuildSucceed:  buildStatus,
		isDepolySucceed: deployStatus,
		commitNotes:     commit,
	}
}

// Notify sends the msg to the tg channel
func (n *Notifaction) Notify() tgbotapi.MessageConfig {
	text := fmt.Sprintf("**Commit Note**: `%s` ", n.commitNotes)
	text = fmt.Sprintf("%s\n **Test Status**: %v", text, n.isTestSucceed)
	text = fmt.Sprintf("%s\n **Build Status**: %v", text, n.isBuildSucceed)
	text = fmt.Sprintf("%s\n **Deploy Status**: %v", text, n.isDepolySucceed)
	return tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChannelUsername: ChannelName,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}
}
