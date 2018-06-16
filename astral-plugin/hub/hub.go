package hub

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	plugin "github.com/scbizu/Astral/astral-plugin"
	"github.com/scbizu/Astral/astral-plugin/py"
	"github.com/scbizu/Astral/astral-plugin/sayhi"
	"github.com/scbizu/Astral/astral-plugin/today-anime"
	"github.com/sirupsen/logrus"
)

// TGPluginHub defines the telegram plugin hub
type TGPluginHub struct {
	plugins []*plugin.TGPlugin
}

// NewTGPluginHub init empty plugin hub
func NewTGPluginHub(msg *tgbotapi.Message) *TGPluginHub {
	hub := &TGPluginHub{
		plugins: []*plugin.TGPlugin{},
	}
	hub.Init(msg)
	return hub
}

// Init regist all command
func (ph *TGPluginHub) Init(msg *tgbotapi.Message) {
	ph.AddPlugin(sayhi.Sayhi(msg))
	ph.AddPlugin(anime.Anime(msg))
	// py plugin must be put in the end
	defer ph.AddPlugin(py.PY(msg))
}

// GetEnabledTelegramPlugins get the all enable plugins
func (ph *TGPluginHub) GetEnabledTelegramPlugins() (activePlugins []*plugin.TGPlugin) {
	for _, p := range ph.plugins {
		if p.IsPluginEnable() {
			activePlugins = append(activePlugins, p)
		}
	}
	return
}

// AddPlugin adds the plugin
func (ph *TGPluginHub) AddPlugin(p *plugin.TGPlugin) {
	ph.plugins = append(ph.plugins, p)
}

// RegistTGEnabledPlugins regists telegram plugin
func (ph *TGPluginHub) RegistTGEnabledPlugins(rawmsg *tgbotapi.Message) (msg tgbotapi.MessageConfig) {

	for _, p := range ph.GetEnabledTelegramPlugins() {
		msg, _ = p.Run(rawmsg)
		logrus.Infof("[chatID:%d,msg:%s]", msg.ChatID, msg.Text)
		if p.Validate(msg) {
			return
		}
	}

	return
}
