package telegram

import (
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/pkg/talker"
	"github.com/sirupsen/logrus"
)

const (
	healthCheckDomian = "https://api.scnace.me/"
)

func healthCheck(bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(10 * time.Minute)

	for t := range ticker.C {
		logrus.Infof("health check at %s", t.String())
		resp, err := http.Get(healthCheckDomian)
		if err != nil {
			logrus.Errorf("health check failed:%s", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusBadGateway {
			t := talker.NewServerStatusNotifaction("Nginx is currently DOWN")
			if _, err := bot.Send(t.ServerStatusMsg()); err != nil {
				logrus.Errorf("health check failed:%s", err.Error())
			}
		}
	}
}
