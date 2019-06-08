package push

import (
	"context"
	"fmt"

	"github.com/scbizu/Astral/internal/discord"
	"github.com/scbizu/Astral/internal/telegram"
	"github.com/scbizu/Astral/internal/tl"
)

func NewPushService() *PushService {
	return &PushService{ctx: context.TODO()}
}

type PushService struct {
	ctx context.Context
}

func (ps *PushService) ServePushSC2Event() error {
	tgBot, err := telegram.NewBot(false)
	if err != nil {
		return err
	}
	discordSender, err := discord.NewBot()
	if err != nil {
		return err
	}
	tgSender := telegram.NewTGSender(tgBot)
	f := tl.NewFetcher(tgSender, discordSender)
	if err := f.Do(); err != nil {
		return fmt.Errorf("tl: %s", err.Error())
	}
	return nil
}
