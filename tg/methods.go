package tg

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/fe0b6/botapi/types"
)

// GetFile - Отправляем сообщение
func GetFile(fileID string, opt *MessageOptions) (file File, err error) {
	api := API{AccessToken: opt.Token}

	fdata := SendGetFile{FileID: fileID}

	ans := api.GetFile(fdata)
	if !ans.Ok {
		err = errors.New(ans.Description)
		return
	}

	err = json.Unmarshal(ans.Result, &file)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	return
}

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

	buttons := [][]ReplyKeyboardButton{}
	keyboard := ReplyKeyboardMarkup{
		OneTimeKeyboard: kb.OneTime,
		ResizeKeyboard:  true,
	}

	for _, ba := range kb.Buttons {
		butns := []ReplyKeyboardButton{}
		for _, b := range ba {
			bn := ReplyKeyboardButton{Text: b.Text}

			if b.RequestContact {
				bn.RequestContact = true
			}

			if b.RequestLocation {
				bn.RequestLocation = true
			}

			butns = append(butns, bn)
		}

		buttons = append(buttons, butns)
	}

	keyboard.Keyboard = buttons

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

// GetFile - Отправка сообщения
func (tg *API) GetFile(sgf SendGetFile) (ans APIResponse) {

	return tg.sendJSONData("getFile", sgf)
}
