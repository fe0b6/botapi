package vk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fe0b6/botapi/types"
)

// ParseRequest - Разбираем запрос, проверяем его корректность и приводим к стандартному виду
func ParseRequest(r *http.Request, secret string) (ans *types.Message, err error) {
	// Читаем данные запроса
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	log.Println(string(b))

	// Парсим объект
	var cbo CallBackObj
	err = json.Unmarshal(b, &cbo)
	if err != nil {
		log.Println("[error]", err)
		log.Println(string(b))
		return
	}

	// Проверяем ключ
	if secret != cbo.Secret {
		err = errors.New("bad secret")
		log.Println("[error]", err, cbo.Secret)
		return
	}

	// преобразуем в стандартный вид
	ans, err = toStandard(&cbo)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

// Приводим к стандартному виду
func toStandard(cbo *CallBackObj) (ans *types.Message, err error) {
	err = cbo.Parse()
	if err != nil {
		log.Println("[error]", err)
		return
	}

	switch cbo.Type {
	case "confirmation":
		ans = &types.Message{IsConfirmation: true}
	case "message_allow":
		ans = &types.Message{
			IsAllow: true,
			FromID:  types.ID{ID: cbo.MessageAllow.UserID},
			Text:    cbo.MessageAllow.Key,
		}
	case "message_deny":
		ans = &types.Message{
			IsAllow: false,
			FromID:  types.ID{ID: cbo.MessageAllow.UserID},
		}
	case "message_new":
		ans = &types.Message{
			FromID: types.ID{ID: cbo.Message.FromID},
			ChatID: types.ID{ID: cbo.Message.PeerID},
			Time:   time.Unix(cbo.Message.Date, 0),
			Text:   cbo.Message.Text,
		}

		if cbo.Message.Payload != "" {
			ans.Command = parsePayload(cbo.Message.Payload)
		}
	}

	ans.Source = "vk"

	return
}

func parsePayload(payload string) string {
	var h map[string]string
	err := json.Unmarshal([]byte(payload), &h)
	if err != nil {
		log.Println("[error]", err)
		return ""
	}

	return h["command"]
}
