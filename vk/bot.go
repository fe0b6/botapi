package vk

import (
	"net/http"
	"time"

	"github.com/fe0b6/botapi"
)

// Bot - Структура бота
type Bot struct {
	Token         string
	Host          string
	Transport     *http.Transport
	SendTimeout   time.Duration
	UsePrometheus bool
}

// SendMessage - шлем сообщение
func (b *Bot) SendMessage(msg botapi.OutMessage) *Answer {
	// TODO обвязка статы и прочего
	ans := msg.Send().(*Answer)

	return ans
}

// NewMessage - создаем новое сообщение
func (b *Bot) NewMessage() *OutMessage {
	return NewMessage()
}
