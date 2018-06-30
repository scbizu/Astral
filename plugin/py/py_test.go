package py

import (
	"testing"

	"github.com/scbizu/Astral/telegram/command"
)

func TestFormatPYCommands(t *testing.T) {
	command.NewCommand(command.CommanderName("testCommand"), "testing", nil)
	cmds := command.GetAllCommands()
	expected := "testCommand - testing"
	if !eq(format(cmds), expected) {
		t.Errorf("not equal,get %s,expected %s", format(cmds), expected)
	}
}

func eq(get string, expected string) (isEq bool) {
	isEq = (get == expected)
	return
}
