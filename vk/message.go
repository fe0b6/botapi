package vk

import (
	"time"

	"github.com/fe0b6/botapi"
)

// Answer - Объект ответа
type Answer struct {
	Err error
}

// Error - возвращаем ошибку
func (a *Answer) Error() error {
	return a.Err
}

// OutMessage - объект сообщения
type OutMessage struct {
	ChatID      botapi.ID
	PeerID      botapi.ID
	UserID      botapi.ID
	RandomID    botapi.ID
	Text        string
	SendTimeout time.Duration
	Buttons     []botapi.Button
	Attach      []botapi.Attach

	ReplyTo botapi.ID
	ForwardMessages []botapi.ID

	DisableWebPagePreview bool
}

// NewMessage - Создаем новое сообщение
func NewMessage() *OutMessage {
	return &OutMessage{}
}

// WithChatID - добавляем получателя сообщения
func (m *OutMessage) WithChatID(id botapi.ID) *OutMessage {
	m.ChatID = id
	return m
}

// WithPeerID - добавляем получателя сообщения
func (m *OutMessage) WithPeerID(id botapi.ID) *OutMessage {
	m.PeerID = id
	return m
}

// WithUserID - добавляем получателя сообщения
func (m *OutMessage) WithUserID(id botapi.ID) *OutMessage {
	m.UserID = id
	return m
}

// WithIdempotent  - добавляем ключ идемпотентности
func (m *OutMessage) WithIdempotent(key botapi.ID) *OutMessage {
	m.RandomID = key
	return m
}

// WithText - добавляем текст в сообщение
func (m *OutMessage) WithText(text string) *OutMessage {
	m.Text = text
	return m
}

// WithSendTimeout - добавляем таймаут на передачу сообщения
func (m *OutMessage) WithSendTimeout(timeout time.Duration) *OutMessage {
	m.SendTimeout = timeout
	return m
}

// WithButtons - добавляем кнопки
func (m *OutMessage) WithButtons(buttons ...botapi.Button) *OutMessage {
	m.Buttons = buttons
	return m
}

// WithAttach - добавляем аттач
func (m *OutMessage) WithAttach(attachments ...botapi.Attach) *OutMessage {
	m.Attach = attachments
	return m
}

// WithReplyTo - добавляем ид сообщения на которое отвечаем
func (m *OutMessage) WithReplyTo(id botapi.ID) *OutMessage {
	m.ReplyTo = id
	return m
}

// WithForwardMessages - добавляем ид сообщения которые надо переслать
func (m *OutMessage) WithForwardMessages(ids ...botapi.ID) *OutMessage {
	m.ForwardMessages = ids
	return m
}

// WithDisableWebPagePreview - отключаем превью ссылок
func (m *OutMessage) WithDisableWebPagePreview() *OutMessage {
	m.DisableWebPagePreview = true
	return m
}

// Send - отправляем сообщение
func (m *OutMessage) Send() *Answer {

	return nil
}

// Button - объект кнопки
type Button struct {
}

// Attach - объект аттача
type Attach struct {
}
