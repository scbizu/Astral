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

// Register defines the plugin common register
type Register func(*tgbotapi.Message) tgbotapi.MessageConfig

// Handler defines plugin must impl interface
type Handler interface {
	Register(msg *tgbotapi.Message) func(*tgbotapi.Message) tgbotapi.MessageConfig
}

// TGPlugin defines the common telegram plugin
type TGPlugin struct {
	enable   bool
	register Register
	name     string
}

// TGPluginHub defines the telegram plugin hub
type TGPluginHub struct {
	plugins []*TGPlugin
}

// NewTGPlugin init the tg plugin
func NewTGPlugin(name string, msg *tgbotapi.Message, handler Handler) *TGPlugin {
	return &TGPlugin{
		enable:   true,
		register: handler.Register(msg),
		name:     name,
	}
}

// Run runs the enabled plugins
func (p *TGPlugin) Run(msg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	if p.enable {
		msgConf := p.register(msg)
		return msgConf, nil
	}
	return tgbotapi.MessageConfig{}, fmt.Errorf("plugin %s is not enabled", p.name)
}

// NewTGPluginHub init empty plugin hub
func NewTGPluginHub(msg *tgbotapi.Message) *TGPluginHub {
	hub := &TGPluginHub{
		plugins: []*TGPlugin{},
	}
	hub.InitRegister(msg)
	return hub
}

// InitRegister regist all command
func (ph *TGPluginHub) InitRegister(msg *tgbotapi.Message) {
	ph.AddPlugin(NewTGPlugin(command.CommandSayhi.String(), msg, &sayhi.Handler{}))
	ph.AddPlugin(NewTGPlugin(command.CommandTodayAnime.String(), msg, &anime.Handler{}))
	// py plugin must be put in the end
	defer ph.AddPlugin(NewTGPlugin(command.CommandShowAllCommand.String(), msg, &py.Handler{}))
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

	for _, p := range ph.GetEnabledTelegramPlugins() {
		msg, _ = p.Run(rawmsg)
		logrus.Infof("[chatID:%d,msg:%s]", msg.ChatID, msg.Text)
		if p.validate(msg) {
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
