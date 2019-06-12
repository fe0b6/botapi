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

	ans, err = api.MessagesSend(map[string]string{
		"user_id":          req.FromID.String(),
		"peer_id":          req.ChatID.String(),
		"chat_id":          req.ChatID.String(),
		"random_id":        getRandomID(opt.RandomID),
		"message":          botans.Text,
		"dont_parse_links": "1",
		"disable_mentions": "1",
	})
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

// Формируем случайный id
func getRandomID(id int64) string {
	if id != 0 {
		return strconv.FormatInt(id, 10)
	}

	rand.Seed(time.Now().UnixNano())
	id = time.Now().Unix() + (1500000000 + rand.Int63n(1000000000))

	return strconv.FormatInt(id, 10)
}
