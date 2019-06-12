package botapi

import (
	"log"
	"net/http"

	"github.com/fe0b6/botapi/types"
	"github.com/fe0b6/botapi/vk"
)

// ParseRequest - Разбираем запрос
func ParseRequest(r *http.Request, opt *MessageOptions) (ans *types.Message, err error) {
	switch opt.Source {
	case "vk":
		ans, err = vk.ParseRequest(r, opt.VKSecret)
		if err != nil {
			log.Println("[error]", err)
			return
		}
	}

	return
}
