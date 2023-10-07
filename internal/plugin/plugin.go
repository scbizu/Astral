package plugin

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/scbizu/Astral/internal/telegram/command"
	"github.com/sirupsen/logrus"
)

type Handler func(msg *tgbotapi.Message) tgbotapi.MessageConfig

// TGPlugin defines the common telegram plugin
type TGPlugin struct {
	enable    bool
	configure tgbotapi.MessageConfig
	name      string
}

type IPlugin interface {
	Name() command.CommanderName
	Enable() bool
	Process(*tgbotapi.Message) tgbotapi.MessageConfig
}

// NewTGPlugin init the tg plugin
func NewTGPlugin(name command.CommanderName,
	conf tgbotapi.MessageConfig,
) *TGPlugin {
	return &TGPlugin{
		enable:    true,
		configure: conf,
		name:      name.String(),
	}
}

// IsPluginEnable returns if current plugin was enabled or not
func (p *TGPlugin) IsPluginEnable() bool {
	return p.enable
}

// Run runs the enabled plugins
func (p *TGPlugin) Run(msg *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	if p.enable {
		return p.configure, nil
	}
	return tgbotapi.MessageConfig{}, fmt.Errorf("plugin %s is not enabled", p.name)
}

// Validate validates message
func (p *TGPlugin) Validate(conf tgbotapi.MessageConfig) bool {
	logrus.Infof("[chatid]: %d,[msg Text]:%s", conf.ChatID, conf.Text)
	if conf.ChatID != 0 && conf.Text != "" {
		return true
	}
	return false
}
