package plugin

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

// Handler defines plugin must impl interface
type Handler interface {
	Register(msg *tgbotapi.Message) tgbotapi.MessageConfig
}

// TGPlugin defines the common telegram plugin
type TGPlugin struct {
	enable    bool
	configure tgbotapi.MessageConfig
	name      string
}

// NewTGPlugin init the tg plugin
func NewTGPlugin(name string, msg *tgbotapi.Message, handler Handler) *TGPlugin {
	return &TGPlugin{
		enable:    true,
		configure: handler.Register(msg),
		name:      name,
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
