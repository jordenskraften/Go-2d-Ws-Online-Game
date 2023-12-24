package lobby

type Chat struct {
	Name     string
	Messages []*ChatMessage
}

type ChatMessage struct {
	Username string
	Text     string
	Date     string
}

func NewChat(name string) *Chat {
	return &Chat{
		Name:     name,
		Messages: make([]*ChatMessage, 0),
	}
}
