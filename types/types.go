package types

import (
	"log"
	"strconv"
	"time"
)

// Message - Структура сообщения
type Message struct {
	ID             ID        `json:"id"`
	Text           string    `json:"Text"`
	Command        string    `json:"command"`
	Time           time.Time `json:"time"`
	FromID         ID        `json:"from_id"`
	ChatID         ID        `json:"chat_id"`
	Source         string    `json:"source"`
	IsAllow        bool      `json:"is_allow"`
	IsDeny         bool      `json:"is_deny"`
	IsConfirmation bool      `json:"is_confirmation"`
	Keyboard       Keyboard  `json:"keyboard"`
	Contact        Contact   `json:"contact"`
	//	Attachments
}

// Contact - объект контакта
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      ID     `json:"user_id"`
}

// Keyboard - объект клавиатуры
type Keyboard struct {
	OneTime  bool       `json:"one_time"`
	NeedHide bool       `json:"need_hide"`
	Buttons  [][]Button `json:"buttons"`
}

// Button - объект кнопки
type Button struct {
	Color           string `json:"color"`
	Text            string `json:"text"`
	Command         string `json:"command"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
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
	case int:
		str = strconv.Itoa(id.ID.(int))
	default:
		log.Printf("[error] unknown type %T!\n", t)
	}

	return
}
