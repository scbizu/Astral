package discord

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/scbizu/Astral/internal/config"
)

type Bot struct {
	session *discordgo.Session
	sync.Mutex
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

func (b *Bot) teardown() {
	b.session.Close()
	b.Unlock()
}

func (b *Bot) Send(msg string) error {
	b.Lock()
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.teardown()

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
	b.Lock()
	if err := b.session.Open(); err != nil {
		return "", err
	}
	defer b.teardown()
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
	b.Lock()
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.teardown()
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

func (b *Bot) Edit(msgID string, content string) error {
	return b.editByChannel(config.DiscordCNSC2ChannelID, msgID, content)
}

func (b *Bot) editByChannel(channelID string, msgID string, content string) error {
	b.Lock()
	if err := b.session.Open(); err != nil {
		return err
	}
	defer b.teardown()
	if _, err := b.session.ChannelMessageEditEmbed(channelID, msgID, &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Event Closed",
				Value: content,
			},
		}},
	); err != nil {
		return err
	}

	return nil
}
