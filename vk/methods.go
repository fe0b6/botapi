package vk

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/fe0b6/botapi/types"
)

// SendMessage - Отправляем сообщение
func SendMessage(req *types.Message, botans *types.Message, opt *MessageOptions) (ans int, err error) {
	api := API{AccessToken: opt.Token}

	h := map[string]string{
		"user_id":          req.FromID.String(),
		"peer_id":          req.ChatID.String(),
		"chat_id":          req.ChatID.String(),
		"random_id":        getRandomID(opt.RandomID),
		"message":          botans.Text,
		"dont_parse_links": "1",
		"disable_mentions": "1",
	}

	// Если есть клавиатура
	if len(botans.Keyboard.Buttons) > 0 || botans.Keyboard.NeedHide {
		h["keyboard"] = formatKeyboard(&botans.Keyboard)
	}

	ans, err = api.MessagesSend(h)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// MessagesSend - отправка сообщения
func (vk *API) MessagesSend(params map[string]string) (ans int, err error) {
	// Отправляем запрос
	r, err := vk.request("messages.send", params)
	if err != nil {
		return
	}

	// Парсим данные
	err = json.Unmarshal(r.Response, &ans)
	if err != nil {
		log.Println("[error]", err, string(r.Response))
		return
	}

	return
}

func formatKeyboard(kb *types.Keyboard) string {
	if kb.NeedHide {
		return `{"buttons":[],"one_time":true}`
	}

	keyboard := Keyboard{Buttons: [][]KeyboardButton{}, OneTime: kb.OneTime}

	for _, ba := range kb.Buttons {
		butns := []KeyboardButton{}
		for _, b := range ba {
			if b.Color == "" {
				b.Color = "primary"
			}
			butns = append(butns, KeyboardButton{
				Color: b.Color,
				Action: KeyboardButtonAction{
					Type:  "text",
					Label: b.Text,
				},
			})
		}

		keyboard.Buttons = append(keyboard.Buttons, butns)
	}

	b, err := json.Marshal(keyboard)
	if err != nil {
		log.Println("[error]", err)
		return ""
	}

	return string(b)
}

// Формируем случайный id
func getRandomID(id int64) string {
	if id != 0 {
		return strconv.FormatInt(id, 10)
	}

	rand.Seed(time.Now().UnixNano())
	id = time.Now().Unix() + (1500000000 + rand.Int63n(1000000000))

	return strconv.FormatInt(id, 10)
}
