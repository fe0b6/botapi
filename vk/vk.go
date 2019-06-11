package vk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// APIVersion - используемая версия API
	APIVersion = "5.95"
	// APIMethodURL - URL запросов к API
	APIMethodURL = "https://api.vk.com/method/"
)

var (
	httpTransport *http.Transport
)

func init() {
	httpTransport = &http.Transport{
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     10 * time.Minute,
	}
}

/*
	Запрос к ВК
*/

// Обертка для запроса к ВК
func (vk *API) request(method string, params map[string]string) (ans Response, err error) {
	if vk.AccessToken == "" {
		err = errors.New("no access token")
		log.Println("[error]", err)
		return
	}

	for {
		ans, err = vk.fullRequest(method, params)
		if err != nil {
			if vk.httpErrorWait(method) {
				continue
			}
			return
		}

		// Проверяем ответ
		if ans.Error.ErrorCode != 0 {
			if ans.Error.ErrorMsg == "Too many requests per second" {
				// Ждем между запросами
				if vk.floodWait(method) {
					continue
				}
			} else if ans.Error.ErrorMsg == "Runtime error occurred during code invocation: Comparing values of different or unsupported types" {
				log.Println("[error]", params["code"])
			}

			err = errors.New(ans.Error.ErrorMsg)
			return
		}

		break
	}

	return
}

// Запрос к ВК
func (vk *API) fullRequest(method string, params map[string]string) (ans Response, err error) {

	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}
	if params["v"] == "" {
		q.Add("v", APIVersion)
	}
	q.Add("access_token", vk.AccessToken)

	// Формируем запрос
	req, err := http.NewRequest("POST", APIMethodURL+method, strings.NewReader(q.Encode()))
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Отправляем запрос
	client := &http.Client{Transport: httpTransport}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if !strings.Contains(err.Error(), "connection reset by peer") && !strings.Contains(err.Error(), "context canceled") {
			log.Println("[error]", err)
		}
		return
	}

	// Если проблема с ответом
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.Println("[error]", resp.Status, resp.StatusCode)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Парсим ответ
	err = json.Unmarshal(content, &ans)
	if err != nil {
		log.Println("[error]", method, err, string(content))
		return
	}

	return
}

// Ждем между запросами если вк ответил что запросы слишком частые
func (vk *API) floodWait(method string) (ok bool) {

	// Определяем сколько времени будет ждать
	var sleepTime int
	if vk.retryCount < 5 {
		sleepTime = 1
	} else if vk.retryCount < 10 {
		sleepTime = 2
	} else if vk.retryCount < 20 {
		sleepTime = 3
	} else if vk.retryCount < 25 {
		sleepTime = 5
	} else {
		// Сбрасываем счетчик ожидания
		vk.Lock()
		vk.retryCount = 0
		vk.Unlock()
		return
	}

	// Увеличиваем счетчик
	vk.Lock()
	vk.retryCount++
	vk.Unlock()

	// Ждем
	time.Sleep(time.Duration(sleepTime) * time.Second)

	ok = true
	return
}

// Попытка повтора запроса при ошибки http
func (vk *API) httpErrorWait(method string) (ok bool) {

	if vk.httpRetryCount >= 3 {
		vk.Lock()
		vk.httpRetryCount = 0
		vk.Unlock()
		return
	}

	vk.Lock()
	vk.httpRetryCount++
	vk.Unlock()

	// Ждем
	time.Sleep(1 * time.Second)

	ok = true
	return
}
