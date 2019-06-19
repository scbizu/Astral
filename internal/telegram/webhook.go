package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/internal/plugin/hub"
	"github.com/scbizu/Astral/pkg/getcert"
	"github.com/scbizu/Astral/pkg/talker"
	"github.com/scbizu/Astral/pkg/talker/dce"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	isDebugMode bool
}

func NewBot(isDebugMode bool) (*Bot, error) {
	bot := new(Bot)
	tgConn, err := ConnectTG()
	if err != nil {
		return nil, err
	}
	bot.bot = tgConn
	if isDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
		bot.isDebugMode = true
		tgConn.Debug = true
		logrus.Infof("bot auth passed as %s", tgConn.Self.UserName)
	}
	return bot, nil
}

func listenWebhook() {
	port := fmt.Sprintf(":%s", os.Getenv("LISTENPORT"))
	http.ListenAndServe(port, nil)
}

func (b *Bot) ServeBotUpdateMessage() error {
	go listenWebhook()
	b.bot.RemoveWebhook()
	cert := getcert.NewDomainCert(tgAPIDomain)
	domainWithToken := fmt.Sprintf("%s/%s", cert.GetDomain(), token)
	if _, err := b.bot.SetWebhook(tgbotapi.NewWebhook(domainWithToken)); err != nil {
		logrus.Errorf("notify webhook failed:%s", err.Error())
		return err
	}
	if b.isDebugMode {
		logrus.SetLevel(logrus.DebugLevel)
		info, err := b.bot.GetWebhookInfo()
		if err != nil {
			return err
		}
		logrus.Debug(info.LastErrorMessage, info.LastErrorDate)
	}
	pattern := fmt.Sprintf("/tg/%s", token)
	updatesMsgChannel := b.bot.ListenForWebhook(pattern)

	logrus.Debugf("msg in channel:%d", len(updatesMsgChannel))
	for update := range updatesMsgChannel {
		logrus.Debugf("[raw msg]:%#v\n", update)

		if update.Message == nil {
			continue
		}
		pluginHub := hub.NewTGPluginHub(update.Message)
		msg := pluginHub.RegistTGEnabledPlugins(update.Message)
		if isMsgBadRequest(msg) {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Astral服务酱表示不想理你")
		}

		mbNames, ok := isMsgNewMember(update)
		if ok {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("撒花欢迎新基佬(%v)入群", mbNames))
		}

		_, ok = isMsgLeftMember(update)
		if ok {
			// drop the member left message
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		b.bot.Send(msg)
	}
	return nil
}

func (b *Bot) ServePushAstralServerMessage() {
	go healthCheck(b.bot)
	registerDCEServer(b.bot)
}

func isMsgNewMember(update tgbotapi.Update) ([]string, bool) {
	members := update.Message.NewChatMembers
	if members == nil {
		return nil, false
	}
	if len(*members) == 0 {
		return nil, false
	}
	var mbNames []string
	for _, m := range *members {
		mbNames = append(mbNames, m.String())
	}
	return mbNames, true
}

func isMsgLeftMember(update tgbotapi.Update) (string, bool) {
	if update.Message.LeftChatMember == nil {
		return "", false
	}
	mbName := update.Message.LeftChatMember.String()
	return mbName, true
}

func isMsgBadRequest(msg tgbotapi.MessageConfig) bool {
	if msg.Text == "" || msg.ChatID == 0 {
		return true
	}
	return false
}

func registerDCEServer(bot *tgbotapi.BotAPI) {
	dceListenHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.Errorf("dce: read webhook msg failed: %q", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			defer r.Body.Close()

			logrus.Debugf("dce webhook: %s", string(body))

			dceObj, err := dce.NewDCEObj(string(body))
			if err != nil {
				logrus.Errorf("dce: read webhook msg failed: %q", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			msgConfig := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChannelUsername: talker.ChannelName,
				},
				Text:      dceObj.Fmt(),
				ParseMode: tgbotapi.ModeMarkdown,
			}
			respMsg, err := bot.Send(msgConfig)
			if err != nil {
				logrus.Errorf("telegram bot: send server info failed: %q", err.Error())
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(respMsg.Text))
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Astral has denied your request"))
		}
	}

	http.HandleFunc("/alert/dce", dceListenHandler)

	port := fmt.Sprintf(":%s", os.Getenv("LISTENPORT"))

	http.ListenAndServe(port, nil)
}
