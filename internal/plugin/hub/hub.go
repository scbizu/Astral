package hub

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	plugin "github.com/scbizu/Astral/internal/plugin"
	"github.com/scbizu/Astral/internal/plugin/ai"
	"github.com/scbizu/Astral/internal/plugin/py"
	"github.com/scbizu/Astral/internal/plugin/sayhi"
	anime "github.com/scbizu/Astral/internal/plugin/today-anime"
	"github.com/sirupsen/logrus"
)

// TGPluginHub defines the telegram plugin hub
type TGPluginHub struct {
	msg     *tgbotapi.Message
	plugins []*plugin.TGPlugin
}

var plugins = []plugin.IPlugin{
	&sayhi.Handler{},
	&anime.Handler{},
	&ai.AICommands{},
	// py plugin must be put in the end
	&py.Handler{},
}

// NewTGPluginHub init empty plugin hub
func NewTGPluginHub(msg *tgbotapi.Message) *TGPluginHub {
	hub := &TGPluginHub{
		plugins: []*plugin.TGPlugin{},
		msg:     msg,
	}
	hub.Init()
	return hub
}

// Init regist all command
func (ph *TGPluginHub) Init() {
	for _, p := range plugins {
		ph.AddTGPlugin(p)
	}
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

func (ph *TGPluginHub) AddTGPlugin(
	p plugin.IPlugin,
) {
	ph.plugins = append(ph.plugins,
		plugin.NewTGPlugin(p.Name(),
			p.Process(ph.msg),
		),
	)
}

// Do iters telegram plugin
func (ph *TGPluginHub) Do(rawmsg *tgbotapi.Message) (msg tgbotapi.Chattable) {
	for _, p := range ph.GetEnabledTelegramPlugins() {
		msg, _ = p.Run(rawmsg)
		if _, ok := msg.(tgbotapi.MessageConfig); ok {
			msgConf := msg.(tgbotapi.MessageConfig)
			logrus.Infof("[chatID:%d,msg:%s]", msgConf.ChatID, msgConf.Text)
			if p.Validate(msgConf) {
				return
			}
		}
	}

	return
}
