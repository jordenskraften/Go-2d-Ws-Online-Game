package lobby

import (
	"log"
	"sync"
	"time"
)

type Chat struct {
	Name     string
	Messages []*ChatMessage
	mu       sync.RWMutex
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

func (ch *Chat) AddChatMessage(username string, text string) *ChatMessage {
	msg := ChatMessage{
		Username: username,
		Text:     text,
		Date:     time.Now().Format("2006-01-02 15:04:05"),
	}
	ch.mu.Lock()
	ch.Messages = append(ch.Messages, &msg)
	ch.mu.Unlock()
	//для дебага
	messages := ch.GetLastFiveMessages()
	log.Println("--------------------")
	log.Println("Пять последних сообщений чата")
	for _, msg := range messages {
		log.Printf("Username: %s, Text: %s, Date: %s\n", msg.Username, msg.Text, msg.Date)
	}
	log.Println("--------------------")
	return &msg
}

func (ch *Chat) GetAllMessages() []*ChatMessage {
	ch.mu.RLock()
	defer ch.mu.RUnlock()
	messages := make([]*ChatMessage, len(ch.Messages))
	copy(messages, ch.Messages)
	return messages
}

func (ch *Chat) GetLastFiveMessages() []*ChatMessage {
	ch.mu.RLock()
	defer ch.mu.RUnlock()

	totalMessages := len(ch.Messages)
	startIndex := 0

	if totalMessages > 5 {
		startIndex = totalMessages - 5
	}

	lastFiveMessages := make([]*ChatMessage, 0)
	for i := startIndex; i < totalMessages; i++ {
		msg := *ch.Messages[i]
		lastFiveMessages = append(lastFiveMessages, &msg)
	}

	return lastFiveMessages
}
