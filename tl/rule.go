// Package tl provides TeamLiquid API wrappers
package tl

// RuleConfig is the TL API TOS spec
type RuleConfig struct {
	commonActionDuration int64
	// specialActionDuration includes "action=parse","action=ask","action=askargs"
	specialActionDuration int64

	userAgent string
}

func NewRuleConfig() RuleConfig {
	return RuleConfig{
		commonActionDuration:  2,
		specialActionDuration: 30,
		userAgent:             "CNSC2EventsBot(https://t.me/AstralAwesomeBot;scbizu@gmail.com)",
	}
}

func (r RuleConfig) GetCommonActionDuration() int64 {
	return r.commonActionDuration
}

func (r RuleConfig) GetSpecialActionDuration() int64 {
	return r.specialActionDuration
}

func (r RuleConfig) GetUA() string {
	return r.userAgent
}
