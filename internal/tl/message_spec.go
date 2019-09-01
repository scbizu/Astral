package tl

type Filter interface {
	F(string) string
}

type Sender interface {
	Send(string, ...Filter) error
	SendAndReturnID(string, ...Filter) (string, error)
	ResolveMessage([]string) string
}

type Editor interface {
	Edit(msgID string, content string) error
}

type IRC interface {
	Sender
	Editor
}
