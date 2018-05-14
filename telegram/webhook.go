package telegram

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	plugin "github.com/scbizu/Astral/astral-plugin"
	"github.com/scbizu/Astral/getcert"
	"github.com/scbizu/Astral/talker"
	"github.com/scbizu/Astral/talker/dce"
	"github.com/sirupsen/logrus"
)

//ListenWebHook is the tg api webhook mode
func ListenWebHook(debug bool) (err error) {
	bot, err := ConnectTG()
	if err != nil {
		return
	}
	if debug {
		bot.Debug = true
		logrus.Infof("bot auth passed as %s", bot.Self.UserName)
	}
	bot.RemoveWebhook()
	cert := getcert.NewDomainCert(tgAPIDomain)
	domainWithToken := fmt.Sprintf("%s%s", cert.GetDomain(), token)
	if _, err = bot.SetWebhook(tgbotapi.NewWebhook(domainWithToken)); err != nil {
		logrus.Infof("notify webhook failed:%s", err.Error())
		return
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		return err
	}

	logrus.Info(info.LastErrorMessage, info.LastErrorDate)

	pattern := fmt.Sprintf("/%s", token)
	updatesMsgChannel := bot.ListenForWebhook(pattern)
	logrus.Infof("msg in channel:%d", len(updatesMsgChannel))

	registDCEServer(bot)

	for update := range updatesMsgChannel {
		logrus.Infof("[raw msg]:%#v\n", update)

		if update.Message == nil {
			continue
		}
		pluginHub := plugin.NewEmptyTGPluginHub()
		var msg tgbotapi.MessageConfig
		if isMsgBadRequest(pluginHub.RegistTGEnabledPlugins(update.Message)) {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Astral服务酱表示不想理你")
		}

		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}

	return
}

func isMsgBadRequest(msg tgbotapi.MessageConfig) bool {
	if msg.Text == "" || msg.ChatID == 0 {
		return true
	}
	return false
}

func registDCEServer(bot *tgbotapi.BotAPI) {
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

	go http.ListenAndServe(port, nil)
}
