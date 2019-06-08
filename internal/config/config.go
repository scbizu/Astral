package config

import (
	"fmt"
	"os"
)

const (
	DiscordCNSC2ChannelID = "586225314078654484"
	discordBotClientIDKey = "ASTRAL_DISCORD_CLIENT_ID"
)

func GetDiscordClientID() string {
	return fmt.Sprintf("Bot %s", os.Getenv(discordBotClientIDKey))
}
