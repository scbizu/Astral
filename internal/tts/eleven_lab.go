package tts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	domain  string
	voiceID string
	apiKey  string
	model   string
}

type BadResponse struct {
	Detail struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"detail"`
}

func NewElevenLabClient() *Client {
	return &Client{
		domain:  "https://api.elevenlabs.io/v1/text-to-speech/",
		voiceID: "21m00Tcm4TlvDq8ikWAM",
		apiKey:  os.Getenv("ELEVEN_LAB_API_KEY"),
		model:   "eleven_monolingual_v1",
	}
}

func (c *Client) ToSpeech(ctx context.Context, text string) ([]byte, error) {
	u, err := url.JoinPath(c.domain, c.voiceID)
	if err != nil {
		return nil, err
	}
	body := make(map[string]any)
	body["text"] = text
	body["model_id"] = c.model
	body["voice_setting"] = map[string]any{
		"stability":        0.5,
		"similarity_boost": 0.5,
	}
	bodyBytes := bytes.NewBuffer(nil)
	if err := json.NewEncoder(bodyBytes).Encode(body); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, u, bodyBytes)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		details := &BadResponse{}
		if err := json.NewDecoder(bytes.NewBuffer(respBody)).Decode(&details); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[11lab] bad response: %s", details.Detail.Message)
	}

	return respBody, nil
}
