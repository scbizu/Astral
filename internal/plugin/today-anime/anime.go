package anime

import (
	"bytes"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	plugin "github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
	"github.com/sirupsen/logrus"
)

// Handler impl the PluginHandler
type Handler struct{}

func getCurrentDay() time.Weekday {
	return time.Now().Weekday()
}

//Register register anime plugin
func (h *Handler) Register(msg *tgbotapi.Message) tgbotapi.MessageConfig {

	handler := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
		animeInfo, err := GetAllAnimes()
		if err != nil {
			asError := fmt.Errorf("astral has an error:[%s]", err.Error())
			return tgbotapi.NewMessage(msg.Chat.ID, asError.Error())
		}
		allmsg := bytes.NewBufferString("")
		for _, info := range animeInfo {
			allmsg.WriteString(info.FormatLinkInMarkdownPreview())
			allmsg.WriteByte('\n')
		}
		if len(animeInfo) == 0 {
			allmsg.WriteString("nothing update today.")
		}
		logrus.Infof("all anime msg:%s", allmsg.String())
		return tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           msg.Chat.ID,
				ReplyToMessageID: 0,
			},
			Text:      allmsg.String(),
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}

	commandDetails := "fetch all anime today!TIPS:D站为了保持一致,默认全部已更新"
	todayAnime := command.NewCommand(command.CommandTodayAnime, commandDetails, handler)
	return todayAnime.Do(msg)

}

// Anime returns anime plugin
func Anime(msg *tgbotapi.Message) *plugin.TGPlugin {
	return plugin.NewTGPlugin(command.CommandTodayAnime.String(), msg, plugin.Handler(&Handler{}))
}
