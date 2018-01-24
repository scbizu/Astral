package anime

import (
	"bytes"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/telegram/command"
)

func getCurrentDay() time.Weekday {
	return time.Now().Weekday()
}

//Register regists anime plugin
func Register(msg *tgbotapi.Message) tgbotapi.MessageConfig {
	handler := func(msg *tgbotapi.Message) tgbotapi.MessageConfig {
		log.Printf("fetching anime")
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
		log.Printf("all anime msg:%s", allmsg.String())
		return tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           msg.Chat.ID,
				ReplyToMessageID: 0,
			},
			Text:      allmsg.String(),
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}
	commandDetails := "fetch all anime today!"
	todayAnime := command.NewCommand(command.CommandTodayAnime, commandDetails, handler)
	return todayAnime.Do(msg)
}
