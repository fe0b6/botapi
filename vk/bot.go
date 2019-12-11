package vk

import (
	"errors"
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
	// TODO добавить сборщик статы
	message, ok := msg.(*OutMessage)
	if !ok {
		return &Answer{Err: errors.New(botapi.ErrBadMessageType)}
	}

	// Собираем хэш для отправки
	h := map[string]string{
		"user_id":   message.UserID.String(),
		"peer_id":   message.PeerID.String(),
		"chat_id":   message.ChatID.String(),
		"random_id": message.RandomID.String(),
		"message":   message.Text,
	}

	if message.DisableWebPagePreview {
		h["dont_parse_links"] = "1"
	}

	// TODO обвязка статы и прочего
	//ans := msg.Send().(*Answer)

	return nil
}

// NewMessage - создаем новое сообщение
func (b *Bot) NewMessage() *OutMessage {
	return NewMessage()
}
