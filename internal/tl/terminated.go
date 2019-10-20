package tl

import "github.com/sirupsen/logrus"

type Stash chan Match

var (
	stashMatchChan Stash
)

func GetStashChan() Stash {
	if stashMatchChan == nil {
		stashMatchChan = make(chan Match)
	}
	return stashMatchChan
}

func (s Stash) Put(m Match) {
	s <- m
}

func (s Stash) Run(ircs ...IRC) {
	for match := range s {
		if match.vs.P1 == "" || match.vs.P2 == "" {
			logrus.Warnf("match series parser rules not set")
			continue
		}
		if match.detailURL == nil {
			continue
		}
		vs, err := GetFinalMatchRes(match.detailURL, match.vs.P1, match.vs.P2)
		if err != nil {
			logrus.Warnf("update message: %q", err)
			continue
		}
		match.vs = vs
		match.timeCountingDown = "已结束"
		match.isOnGoing = false
		msg := match.GetMDMatchInfo()
		for _, irc := range ircs {
			content := irc.ResolveMessage([]string{msg})
			if err := irc.Send(content); err != nil {
				logrus.Errorf("send FIN message: %q", err)
				continue
			}
		}
	}
}
