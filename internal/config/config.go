package config

import (
	"fmt"
	"os"
)

var (
	OpenAIAPIKEy      = os.Getenv("OPENAI_API_KEY")
	OpenAIAPIEndpoint = os.Getenv("OPENAI_COMPLETION_ENDPOINT")
)

const (
	DiscordCNSC2ChannelID = "586225314078654484"
	discordBotClientIDKey = "ASTRAL_DISCORD_CLIENT_ID"
)

func GetDiscordClientID() string {
	return fmt.Sprintf("Bot %s", os.Getenv(discordBotClientIDKey))
}
