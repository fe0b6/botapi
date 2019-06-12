package botapi

import (
	"log"
	"net/http"

	"github.com/fe0b6/botapi/tg"
	"github.com/fe0b6/botapi/types"
	"github.com/fe0b6/botapi/vk"
)

// ParseRequest - Разбираем запрос
func (bot *Bot) ParseRequest(r *http.Request, opt *MessageOptions) (ans *types.Message, err error) {
	switch opt.Source {
	case "vk":
		ans, err = vk.ParseRequest(r, bot.VK.Secret)
		if err != nil {
			if err.Error() != "skip" {
				log.Println("[error]", err)
			}
			return
		}
	case "tg":
		ans, err = tg.ParseRequest(r)
		if err != nil {
			if err.Error() != "skip" {
				log.Println("[error]", err)
			}
			return
		}
	}

	return
}

// SendMessage - Отправляем сообщение
func (bot *Bot) SendMessage(req *types.Message, ans *types.Message) (err error) {
	switch req.Source {
	case "vk":
		_, err = vk.SendMessage(req, ans, &vk.MessageOptions{Token: bot.VK.Token})
	case "tg":
		_, err = tg.SendMessage(req, ans, &vk.MessageOptions{Token: bot.TG.Token})
	default:
		log.Println("[error]", "unknown source", req.Source)
		return
	}

	if err != nil {
		if err.Error() != "skip" {
			log.Println("[error]", err)
		}
		return
	}

	return
}
