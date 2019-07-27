package discord

import (
	"fmt"
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
	return fmt.Sprintf("%s\n\n", strings.Join(msgs, "\n\n"))
}

func (b *Bot) Send(msg string) error {
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.session.Close()
	if _, err := b.session.ChannelMessageSendEmbed(
		config.DiscordCNSC2ChannelID,
		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Event Update",
					Value: msg,
				},
			}},
	); err != nil {
		return err
	}
	return nil
}

func (b *Bot) SendAndReturnID(msg string) (string, error) {
	if err := b.session.Open(); err != nil {
		return "", err
	}
	defer b.session.Close()
	resp, err := b.session.ChannelMessageSendEmbed(
		config.DiscordCNSC2ChannelID,
		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Event Update",
					Value: msg,
				},
			}},
	)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (b *Bot) SendToChannel(channelID string, msg string) error {
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.session.Close()
	if _, err := b.session.ChannelMessageSendEmbed(
		channelID,
		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Event Update",
					Value: msg,
				},
			}},
	); err != nil {
		return err
	}

	return nil
}
