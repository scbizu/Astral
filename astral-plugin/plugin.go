package plugin

import (
	"github.com/scbizu/Astral/astral-plugin/lunch"
	"github.com/scbizu/wechat-go/wxweb"
)

//RegisterAllEnabledPlugins scan the setting file(plugin.yaml) and register them
//to the main wx session.
func RegisterAllEnabledPlugins(session *wxweb.Session) {
	// replier.Register(session, autoReply)
	lunch.Register(session, nil)
}
