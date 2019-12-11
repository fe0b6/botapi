package botapi

import (
	"time"
)

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

// Button - объект кнопоки
type Button interface {
}

// Attach - объект аттача
type Attach interface {
}
