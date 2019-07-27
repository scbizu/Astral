package tl

type Sender interface {
	Send(string) error
	ResolveMessage([]string) string
	SendAndReturnID(string) (string, error)
}

type Editor interface {
	Edit(msgID string, content string) error
}

type IRC interface {
	Sender
	Editor
}
