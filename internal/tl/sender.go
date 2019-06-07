package tl

type Sender interface {
	Send(string) error
	ResolveMessage([]string) string
}
