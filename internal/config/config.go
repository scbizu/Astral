package config

import (
	"os"
)

const (
	DiscordCNSC2ChannelName = "CNSC2Event"
	discordBotClientIDKey   = "ASTRAL_DISCORD_CLIENT_ID"
)

func GetDiscordClientID() string {
	return os.Getenv(discordBotClientIDKey)
}
