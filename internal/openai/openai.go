package openai

import (
	"sync"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
	"github.com/scbizu/Astral/internal/config"
)

var (
	openAIClientOnce sync.Once
	openAIChatClient *chat.Client
)

func GetOpenAIClient() *chat.Client {
	openAIClientOnce.Do(func() {
		openAISession := openai.NewSession(config.OpenAIAPIKEy)
		c := chat.NewClient(openAISession, "")
		c.CreateCompletionEndpoint = config.OpenAIAPIEndpoint
		openAIChatClient = c
	})
	return openAIChatClient
}
