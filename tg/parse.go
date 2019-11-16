package tg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fe0b6/botapi/types"
)

// ParseRequest - Разбираем запрос, проверяем его корректность и приводим к стандартному виду
func ParseRequest(r *http.Request, opt *ParseOptions) (ans *types.Message, err error) {
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
	ans, err = toStandard(&upd, opt)
	if err != nil {
		if err.Error() != "skip" {
			log.Println("[error]", err)
		}
		return
	}

	return
}

// Приводим к стандартному виду
func toStandard(upd *Update, opt *ParseOptions) (ans *types.Message, err error) {
	ans = &types.Message{
		ID:     types.ID{ID: upd.Message.MessageID},
		FromID: types.ID{ID: upd.Message.From.ID},
		ChatID: types.ID{ID: upd.Message.Chat.ID},
		Time:   time.Unix(upd.Message.Date, 0),
		Text:   upd.Message.Text,
		Contact: types.Contact{
			PhoneNumber: upd.Message.Contact.PhoneNumber,
			FirstName:   upd.Message.Contact.FirstName,
			LastName:    upd.Message.Contact.LastName,
			UserID:      types.ID{ID: upd.Message.Contact.UserID},
		},
		Photos: make([]types.Photo, 0, len(upd.Message.Photo)),
	}

	for _, ph := range upd.Message.Photo {
		ans.Photos = append(ans.Photos, types.Photo{
			URL:      fmt.Sprintf(FileEndpoint, opt.Token, ph.FileID),
			Width:    ph.Width,
			Height:   ph.Height,
			FileSize: ph.FileSize,
		})
	}

	ans.Source = "tg"

	return
}
