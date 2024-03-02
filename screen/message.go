package screen

import (
	"sync"
)

var maxMessages = 20

type SafeMessageQueue struct {
	Mu       sync.Mutex
	Messages []string
}

func NewSafeMessageQueue() *SafeMessageQueue {
	return &SafeMessageQueue{
		Messages: make([]string, 0, maxMessages),
	}
}

func (smq *SafeMessageQueue) AddMessage(message string) {
	smq.Mu.Lock()
	defer smq.Mu.Unlock()

	if len(smq.Messages) == maxMessages {
		smq.Messages = smq.Messages[1:]
	}

	smq.Messages = append(smq.Messages, message)
}
