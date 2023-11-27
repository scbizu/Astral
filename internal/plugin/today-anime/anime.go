package anime

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	plugin "github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/telegram/command"
)

var _ plugin.IPlugin = (*Handler)(nil)

// Handler impl the PluginHandler
type Handler struct{}

// Register register anime plugin
func (h *Handler) Process(msg *tgbotapi.Message) tgbotapi.Chattable {
	handler := func(msg *tgbotapi.Message) tgbotapi.Chattable {
		return tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           msg.Chat.ID,
				ReplyToMessageID: 0,
			},
			Text:      "施工中",
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}

	commandDetails := "fetch all anime today!TIPS:D站为了保持一致,默认全部已更新"
	todayAnime := command.NewCommand(command.CommandTodayAnime, commandDetails, handler)
	return todayAnime.Do(msg)
}

func (h *Handler) Name() command.CommanderName {
	return command.CommandTodayAnime
}

func (h *Handler) Enable() bool {
	return true
}
