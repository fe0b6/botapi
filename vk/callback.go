package vk

import (
	"encoding/json"
	"log"
)

// CallBackObj - Объект запроса через callback
type CallBackObj struct {
	Type    string          `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	Secret  string          `json:"secret"`

	Message      MessagesGetAns       `json:"-"`
	MessageAllow CallbackMessageAllow `json:"-"`

	Parsed bool `json:"-"`
}

// Parse - Парсим объект
func (cbo *CallBackObj) Parse() (err error) {
	if cbo.Parsed {
		return
	}

	switch cbo.Type {
	case "message_new", "message_reply", "message_edit":
		err = json.Unmarshal(cbo.Object, &cbo.Message)
	case "message_allow", "message_deny":
		err = json.Unmarshal(cbo.Object, &cbo.MessageAllow)
	}

	if err != nil {
		log.Println("[error]", err)
		return
	}

	cbo.Parsed = true
	return
}

// CallbackMessageAllow - объект подписки на сообщения
type CallbackMessageAllow struct {
	UserID int    `json:"user_id"`
	Key    string `json:"key"`
}
