package tg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fe0b6/botapi/types"
)

// ParseRequest - Разбираем запрос, проверяем его корректность и приводим к стандартному виду
func ParseRequest(r *http.Request) (ans *types.Message, err error) {
	// Читаем данные запроса
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	log.Println(string(b))

	var upd Update
	err = json.Unmarshal(b, &upd)
	if err != nil {
		log.Println("[error]", err)
		log.Println(string(b))
		return
	}

	// преобразуем в стандартный вид
	ans, err = toStandard(&upd)
	if err != nil {
		if err.Error() != "skip" {
			log.Println("[error]", err)
		}
		return
	}

	return
}

// Приводим к стандартному виду
func toStandard(upd *Update) (ans *types.Message, err error) {
	ans = &types.Message{
		ID:     types.ID{ID: upd.Message.MessageID},
		FromID: types.ID{ID: upd.Message.From.ID},
		ChatID: types.ID{ID: upd.Message.Chat.ID},
		Time:   time.Unix(upd.Message.Date, 0),
		Text:   upd.Message.Text,
	}

	ans.Source = "tg"

	return
}
