package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Telegram constants
const (

	// Максимальный размер текста для сообщения
	TextMaxSize = 4000
	// Максимальный размер описания
	CaptionMaxSize = 200
)

var (
	// APIEndpoint is the endpoint for all API methods, with formatting for Sprintf
	APIEndpoint = "https://api.telegram.org/bot%s/%s"
	// FileEndpoint is the endpoint for downloading a file from Telegram
	FileEndpoint = "https://api.telegram.org/file/bot%s/%s"

	httpTransport *http.Transport
)

func init() {
	httpTransport = &http.Transport{
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     10 * time.Minute,
	}
}

// SetEndpoint - Устанавливаем какой endpoint использовать
func SetEndpoint(domain string) {
	APIEndpoint = "https://" + domain + "/bot%s/%s"
	FileEndpoint = "https://" + domain + "/file/bot%s/%s"
}

// Отправляем json
func (tg *API) sendJSONData(method string, data interface{}) (ans APIResponse) {
	for {
		ans = tg.sendJSONDataFull(method, data)

		// Если переборщили с кол-вом сообщенией - подождем и попробуем заново
		if !ans.Ok && ans.ErrorCode == 429 {
			// Если повтора ждать не надо
			if tg.RetryDontWait {
				break
			}

			if tg.floodWait(ans) {
				continue
			}
		}

		break
	}

	return
}

func (tg *API) sendJSONDataFull(method string, data interface{}) (ans APIResponse) {
	// Формируем json данные
	b, err := json.Marshal(&data)
	if err != nil {
		log.Println("[error]", method, err)
		return
	}

	// Если нет клавиатуры - удаляем отметку о ней
	b = bytes.Replace(b, []byte(`,"reply_markup":null`), []byte(""), -1)

	// Формируем запрос
	req, err := http.NewRequest("POST", tg.getRequestURL(method), bytes.NewBuffer(b))
	if err != nil {
		log.Println("[error]", method, err)
		return
	}
	// Добавляем заголовое о том что это json
	req.Header.Set("Content-Type", "application/json")

	// Делаем запрос
	client := &http.Client{Transport: httpTransport}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if resp != nil {
			ans.Description = resp.Status
			ans.ErrorCode = resp.StatusCode
		} else {
			ans.ErrorCode = 500
		}

		log.Println("[error]", method, err)
		return
	}

	// Проверяем ответ
	ans = tg.checkAnswer(method, resp)

	return
}

// Проверяем ответ телеграма
func (tg *API) checkAnswer(method string, resp *http.Response) (ans APIResponse) {
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", method, err)
		return
	}

	// Разбираем ответ
	err = json.Unmarshal(body, &ans)
	if err != nil {
		log.Println("[error]", method, err)
		return
	}

	// Если с ответом не все ок
	if !ans.Ok {
		log.Println("[error]", method, string(body))
		return
	}

	return
}

// Формируем url для запроса
func (tg *API) getRequestURL(method string) string {
	return fmt.Sprintf(APIEndpoint, tg.AccessToken, method)
}

// Ждем между запросами если телеграм ответил что запросы слишком частые
func (tg *API) floodWait(ans APIResponse) (ok bool) {
	// Если уже много ждали - выходим
	if tg.retryCount >= 5 {
		// Сбрасываем счетчик
		tg.Lock()
		tg.retryCount = 0
		tg.Unlock()
		return
	}

	// Увеличиваем счетчик
	tg.Lock()
	tg.retryCount++
	tg.Unlock()

	// Ждем
	time.Sleep(time.Duration(ans.Parameters.RetryAfter) * time.Second)

	ok = true
	return
}
