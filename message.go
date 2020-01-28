package botapi

import (
	"time"
)

// Message - объект мессаджа
type Message struct {
	RecipientID           ID
	ChatID                ID
	SenderID              ID
	ReplyTo               ID
	ForwardMessages       []ID
	IdempotentKey         string
	Text                  string
	DisableWebPagePreview bool
	Keyboard              Keyboard
	Attach                []Attach
	SendTimeout           time.Duration
}

// Keyboard - объект клавиатуры
type Keyboard struct {
	Buttons [][]Button
	OneTime bool
}

// Button - объект кнопки
type Button struct {
	Color           string `json:"color"`
	Text            string `json:"text"`
	Command         string `json:"command"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

// Attach - объект аттача
type Attach struct {
}

// ID - интерфейс айдишника
type ID interface {
	String() string
}

// Answer - объект ответа на сообщение
type Answer interface {
	Error() error
}

// OutMessage - объект исходящего сообщения
type OutMessage interface {
	WithChatID(ID) OutMessage
	WithPeerID(ID) OutMessage
	WithUserID(ID) OutMessage
	WithIdempotent(ID) OutMessage
	WithReplyTo(ID) OutMessage
	WithForwardMessages(...ID) OutMessage
	WithText(string) OutMessage
	WithSendTimeout(time.Duration) OutMessage
	WithButtons(...Button) OutMessage
	WithAttach(...Attach) OutMessage
	WithDisableWebPagePreview() OutMessage
}
