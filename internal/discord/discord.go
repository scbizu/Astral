package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/scbizu/Astral/internal/config"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot() (*Bot, error) {
	bot, err := discordgo.New(config.GetDiscordClientID())
	if err != nil {
		return nil, err
	}
	return &Bot{session: bot}, nil
}

func (b *Bot) ResolveMessage(msgs []string) string {
	return strings.Join(msgs, "\n\n")
}

func (b *Bot) Send(msg string) error {
	if _, err := b.session.ChannelMessageSend(
		config.DiscordCNSC2ChannelName,
		msg,
	); err != nil {
		return err
	}

	return nil
}
