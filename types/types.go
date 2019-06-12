package types

import (
	"log"
	"strconv"
	"time"
)

// Message - Структура сообщения
type Message struct {
	ID      ID        `json:"id"`
	Text    string    `json:"Text"`
	Time    time.Time `json:"time"`
	FromID  ID        `json:"from_id"`
	ChatID  ID        `json:"chat_id"`
	Source  string    `json:"source"`
	IsAllow bool      `json:"is_allow"`
	IsDeny  bool      `json:"is_deny"`
	//	Attachments
}

// ID - объект айди
type ID struct {
	ID interface{}
}

func (id *ID) String() (str string) {
	switch t := id.ID.(type) {
	case string:
		str = id.ID.(string)
	case float64:
		str = strconv.FormatFloat(id.ID.(float64), 'f', -1, 64)
	case int64:
		str = strconv.FormatInt(id.ID.(int64), 10)
	default:
		log.Println("[error]", "unknown type", t)
	}

	return
}
