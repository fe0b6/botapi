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

	ReplyTo         botapi.ID
	ForwardMessages []botapi.ID

	DisableWebPagePreview bool
}

// NewMessage - Создаем новое сообщение
func NewMessage() *OutMessage {
	return &OutMessage{}
}

// WithChatID - добавляем получателя сообщения
func (m *OutMessage) WithChatID(id botapi.ID) botapi.OutMessage {
	m.ChatID = id
	return m
}

// WithPeerID - добавляем получателя сообщения
func (m *OutMessage) WithPeerID(id botapi.ID) botapi.OutMessage {
	m.PeerID = id
	return m
}

// WithUserID - добавляем получателя сообщения
func (m *OutMessage) WithUserID(id botapi.ID) botapi.OutMessage {
	m.UserID = id
	return m
}

// WithIdempotent  - добавляем ключ идемпотентности
func (m *OutMessage) WithIdempotent(key botapi.ID) botapi.OutMessage {
	m.RandomID = key
	return m
}

// WithReplyTo - добавляем ид сообщения на которое отвечаем
func (m *OutMessage) WithReplyTo(id botapi.ID) botapi.OutMessage {
	m.ReplyTo = id
	return m
}

// WithForwardMessages - добавляем ид сообщения которые надо переслать
func (m *OutMessage) WithForwardMessages(ids ...botapi.ID) botapi.OutMessage {
	m.ForwardMessages = ids
	return m
}

// WithText - добавляем текст в сообщение
func (m *OutMessage) WithText(text string) botapi.OutMessage {
	m.Text = text
	return m
}

// WithSendTimeout - добавляем таймаут на передачу сообщения
func (m *OutMessage) WithSendTimeout(timeout time.Duration) botapi.OutMessage {
	m.SendTimeout = timeout
	return m
}

// WithButtons - добавляем кнопки
func (m *OutMessage) WithButtons(buttons ...botapi.Button) botapi.OutMessage {
	m.Buttons = buttons
	return m
}

// WithAttach - добавляем аттач
func (m *OutMessage) WithAttach(attachments ...botapi.Attach) botapi.OutMessage {
	m.Attach = attachments
	return m
}

// WithDisableWebPagePreview - отключаем превью ссылок
func (m *OutMessage) WithDisableWebPagePreview() botapi.OutMessage {
	m.DisableWebPagePreview = true
	return m
}

// Button - объект кнопки
type Button struct {
}

// Attach - объект аттача
type Attach struct {
}
