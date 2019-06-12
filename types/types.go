package types

import "time"

// Message - Структура сообщения
type Message struct {
	ID      interface{} `json:"id"`
	Text    string      `json:"Text"`
	Time    time.Time   `json:"time"`
	FromID  interface{} `json:"from_id"`
	Source  string      `json:"source"`
	IsAllow bool        `json:"is_allow"`
	IsDeny  bool        `json:"is_deny"`
	//	Attachments
}
