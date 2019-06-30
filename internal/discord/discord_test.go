package discord

import "testing"

const (
	testChannelID = "586817877366276106"
)

func TestSendMessage(t *testing.T) {
	s, err := NewBot()
	if err != nil {
		t.Fatalf("discord bot: %s", err.Error())
	}

	if err := s.SendToChannel(testChannelID, "Astral: [Hello Discord](https://discord.gg)"); err != nil {
		t.Fatalf("discord bot: %q", err.Error())
	}
}
