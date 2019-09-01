package tl

import "github.com/sirupsen/logrus"

type Stash chan PMatch

var (
	stashMatchChan Stash
)

func GetStashChan() Stash {
	if stashMatchChan == nil {
		stashMatchChan = make(chan PMatch)
	}
	return stashMatchChan
}

func (s Stash) Put(m PMatch) {
	s <- m
}

func (s Stash) Run(ircs ...IRC) {
	for {
		select {
		case msg := <-s:
			match := msg.rawMatches[msg.matchIndex]
			vs, err := GetFinalMatchRes(match.detailURL, match.vs.P1, match.vs.P2)
			if err != nil {
				logrus.Warnf("update message: %q", err)
				vs = Versus{
					P1:      match.vs.P1,
					P2:      match.vs.P2,
					P1Score: match.vs.P1Score,
					P2Score: match.vs.P2Score,
				}
			}
			msg.rawMatches[msg.matchIndex] = Match{
				isOnGoing:        false,
				vs:               vs,
				timeCountingDown: "已结束",
				series:           match.series,
			}
			var strs []string
			for _, m := range msg.rawMatches {
				strs = append(strs, m.GetMDMatchInfo())
			}
			for _, irc := range ircs {
				content := irc.ResolveMessage(strs)
				if err := irc.Send(content); err != nil {
					logrus.Errorf("send FIN message: %q", err)
					continue
				}
			}
		}
	}
}
