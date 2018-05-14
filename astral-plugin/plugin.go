package plugin

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scbizu/Astral/astral-plugin/py"
	"github.com/scbizu/Astral/astral-plugin/sayhi"
	"github.com/scbizu/Astral/astral-plugin/today-anime"
	"github.com/scbizu/Astral/telegram/command"
	"github.com/sirupsen/logrus"
)

// Handler defines the plugin handler
type Handler func(*tgbotapi.Message) tgbotapi.MessageConfig

// TGPlugin defines the common telegram plugin
type TGPlugin struct {
	enable bool
	handle Handler
	name   string
}

// TGPluginHub defines the telegram plugin hub
type TGPluginHub struct {
	plugins []*TGPlugin
}

// NewTGPlugin init the tg plugin
func NewTGPlugin(name string, handler Handler) *TGPlugin {
	return &TGPlugin{
		enable: true,
		handle: handler,
		name:   name,
	}
}

// Run runs the enabled plugins
func (p *TGPlugin) Run(msg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	if p.enable {
		msgConf := p.handle(msg)
		return msgConf, nil
	}
	return tgbotapi.MessageConfig{}, fmt.Errorf("plugin %s is not enabled", p.name)
}

// NewEmptyTGPluginHub init empty plugin hub
func NewEmptyTGPluginHub() *TGPluginHub {
	return &TGPluginHub{
		plugins: []*TGPlugin{},
	}
}

// InitRegister regist all command
func (ph *TGPluginHub) InitRegister(msg *tgbotapi.Message) {
	ph.AddPlugin(NewTGPlugin(command.CommandSayhi.String(), Handler(sayhi.Register(msg))))
	ph.AddPlugin(NewTGPlugin(command.CommandTodayAnime.String(), Handler(anime.Register(msg))))
	// py plugin must be put in the end
	ph.AddPlugin(NewTGPlugin(command.CommandShowAllCommand.String(), Handler(py.Register(msg))))
}

// GetEnabledTelegramPlugins get the all enable plugins
func (ph *TGPluginHub) GetEnabledTelegramPlugins() (activePlugins []*TGPlugin) {
	for _, p := range ph.plugins {
		if p.enable {
			activePlugins = append(activePlugins, p)
		}
	}
	return
}

// AddPlugin adds the plugin
func (ph *TGPluginHub) AddPlugin(p *TGPlugin) {
	ph.plugins = append(ph.plugins, p)
}

// RegistTGEnabledPlugins regists telegram plugin
func (ph *TGPluginHub) RegistTGEnabledPlugins(rawmsg *tgbotapi.Message) (msg tgbotapi.MessageConfig) {

	ph.InitRegister(rawmsg)

	for _, p := range ph.GetEnabledTelegramPlugins() {
		msg, _ = p.Run(rawmsg)
		logrus.Infof("[chatID:%d,msg:%s]", msg.ChatID, msg.Text)
		if !p.validate(msg) {
			return
		}
	}

	return
}

func (p *TGPlugin) validate(conf tgbotapi.MessageConfig) bool {
	logrus.Infof("[chatid]: %d,[msg Text]:%s", conf.ChatID, conf.Text)
	if conf.ChatID != 0 && conf.Text != "" {
		return true
	}
	return false
}
