package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/getcert"
	"github.com/scbizu/Astral/plugin/hub"
	"github.com/scbizu/Astral/talker"
	"github.com/scbizu/Astral/talker/dce"
	"github.com/scbizu/Astral/tl"
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

func (b *Bot) ServeBotUpdateMessage() error {
	b.bot.RemoveWebhook()
	cert := getcert.NewDomainCert(tgAPIDomain)
	domainWithToken := fmt.Sprintf("%s%s", cert.GetDomain(), token)
	if _, err := b.bot.SetWebhook(tgbotapi.NewWebhook(domainWithToken)); err != nil {
		logrus.Errorf("notify webhook failed:%s", err.Error())
		return err
	}
	if b.isDebugMode {
		info, err := b.bot.GetWebhookInfo()
		if err != nil {
			return err
		}

		logrus.Debug(info.LastErrorMessage, info.LastErrorDate)
	}
	pattern := fmt.Sprintf("/%s", token)
	if b.isDebugMode {
		logrus.Debugf("token: %s", token)
	}
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

func (b *Bot) ServePushSC2Event() {
	f := tl.NewFetcher(b.bot)
	if err := f.Do(); err != nil {
		logrus.Errorf("tl: %s", err.Error())
		return
	}
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
				logrus.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			defer r.Body.Close()
			dceObj, err := dce.NewDCEObj(string(body))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			noti := talker.NewNotifaction(dceObj.GetRepoName(),
				dceObj.GetStageMap(), dceObj.GetCommitMsg(),
				dceObj.GetBuildDuration())
			respMsg, err := bot.Send(noti.Notify())
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(respMsg.Text))
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Astral denied your request"))
		}
	}

	http.HandleFunc("/dce", dceListenHandler)

	port := fmt.Sprintf(":%s", os.Getenv("LISTENPORT"))

	http.ListenAndServe(port, nil)
}
