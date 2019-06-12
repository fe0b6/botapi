package tg

import (
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/fe0b6/botapi/types"
)

// SendMessage - Отправляем сообщение
func SendMessage(req *types.Message, botans *types.Message, opt *MessageOptions) (err error) {
	api := API{AccessToken: opt.Token}

	h := SendMessageData{
		ChatID:                req.ChatID.ID,
		Text:                  botans.Text,
		DisableWebPagePreview: true,
	}

	// Если есть клавиатура
	if len(botans.Keyboard.Buttons) > 0 || botans.Keyboard.NeedHide {
		h.ReplyMarkup = formatKeyboard(&botans.Keyboard)
	}

	ans := api.SendMessageBig(h)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	for _, a := range ans {
		if !a.Ok {
			err = errors.New(a.Description)
			log.Println("[error]", err, a.ErrorCode)
			return
		}
	}

	return
}

func formatKeyboard(kb *types.Keyboard) interface{} {
	if kb.NeedHide {
		return ReplyKeyboardMarkup{Keyboard: [][]string{}}
	}

	keyboard := ReplyKeyboardMarkup{
		Keyboard:        [][]string{},
		OneTimeKeyboard: kb.OneTime,
		ResizeKeyboard:  true,
	}

	for _, ba := range kb.Buttons {
		butns := []string{}
		for _, b := range ba {
			butns = append(butns, b.Text)
		}

		keyboard.Keyboard = append(keyboard.Keyboard, butns)
	}

	return keyboard
}

// GetMe - Получаем инфу о боте
func (tg *API) GetMe() (ans APIResponse) {
	return tg.sendJSONData("getMe", nil)
}

// SetWebhook - Установка Webhook
func (tg *API) SetWebhook(url string) (ans APIResponse) {
	m := map[string]string{"url": url}

	return tg.sendJSONData("setWebhook", m)
}

// SendMessage - Отправка сообщения
func (tg *API) SendMessage(msg SendMessageData) (ans APIResponse) {

	// Если клавиатура не указана - делаем пустую
	if msg.HideReplyMarkup && msg.ReplyMarkup == nil {
		msg.ReplyMarkup = ReplyKeyboardMarkup{Keyboard: [][]string{}}
	}

	return tg.sendJSONData("sendMessage", msg)
}

// SendMessageBig - Отправка сообщения с проверкой на длину
func (tg *API) SendMessageBig(msg SendMessageData) (ans []APIResponse) {
	ans = []APIResponse{}

	// Если длина текста влезет в одно сообщение - просто отправляем
	if utf8.RuneCountInString(msg.Text) < TextMaxSize {
		ans = append(ans, tg.SendMessage(msg))
		return
	}

	// Разбиваем текст на блоки нужной длины
	texts := []string{}
	var tmp string
	for _, v := range strings.Split(msg.Text, " ") {
		// Если длина куска будет больше чем максимум - сохраняем предыдущий и начинаем новый кусок
		if utf8.RuneCountInString(tmp+v) >= TextMaxSize {
			texts = append(texts, tmp)
			tmp = ""
		}

		tmp += v + " "
	}
	// Не забываем добавить остаток текста
	if len(tmp) > 0 {
		texts = append(texts, tmp)
	}

	// Отправляем куски
	for _, text := range texts {
		msg.Text = text
		ans = append(ans, tg.SendMessage(msg))
	}

	return
}
