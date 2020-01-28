package vk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fe0b6/botapi"
)

const (
	// APIVersion - используемая версия API
	APIVersion = "5.95"
	// APIMethodURL - URL запросов к API
	APIMethodURL = "https://%s/method/"
	// DefaultDomain - домен по умолчанию
	DefaultDomain = "api.vk.com"

	defaultSendTimeout = time.Minute * 15
)

var (
	defaultRetrySettings     = map[int]int{5: 1, 10: 5, 25: 10}
	defaultHTTPRetrySettings = map[int]int{5: 5, 10: 15, 25: 30}

	defaultHTTPTransport *http.Transport
)

func init() {
	defaultHTTPTransport = &http.Transport{
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     10 * time.Minute,
	}
}

// Bot - Структура бота
type Bot struct {
	Token             string
	Host              string
	Transport         *http.Transport
	SendTimeout       time.Duration
	UsePrometheus     bool
	Rand              *rand.Rand
	RetrySettings     map[int]int
	HTTPRetrySettings map[int]int
	Domain            string
}

// SendMessage - шлем сообщение
func (b *Bot) SendMessage(ctx context.Context, message *botapi.Message) *Answer {

	resp, err := b.request(ctx, "messages.send", b.FormatMessage(message), message.SendTimeout)
	if err != nil {
		return &Answer{Err: err}
	}

	if resp.Error.ErrorCode != 0 {
		return &Answer{Err: fmt.Errorf("%s, with code: %d", resp.Error.ErrorMsg, resp.Error.ErrorCode)}
	}

	return nil
}

// FormatMessage - Приводим сообщение к нужному виде для отправки
func (b *Bot) FormatMessage(message *botapi.Message) map[string]string {
	// Собираем хэш для отправки
	h := map[string]string{
		"user_id":   message.SenderID.String(),
		"peer_id":   message.RecipientID.String(),
		"chat_id":   message.ChatID.String(),
		"random_id": strconv.Itoa(b.Rand.Int()), // TODO check max value size
		"message":   message.Text,
	}

	if message.DisableWebPagePreview {
		h["dont_parse_links"] = "1"
	}

	// Собираем клавиатуру
	if len(message.Keyboard.Buttons) > 0 {
		buttons := make([][]map[string]interface{}, len(message.Keyboard.Buttons))

		for _, buttonArr := range message.Keyboard.Buttons {
			btns := make([]map[string]interface{}, len(buttonArr))

			for _, button := range buttonArr {
				color := "primary"
				if button.Color != "" {
					color = button.Color
				}

				btns = append(btns, map[string]interface{}{
					"color": color,
					"action": map[string]string{
						"type":    "text",
						"label":   button.Text,
						"payload": button.Command,
					},
				})
			}

			buttons = append(buttons, btns)
		}

		buttonJSON, err := json.Marshal(buttons)
		if err != nil {
			log.Fatalln(err)
		}

		h["keyboard"] = string(buttonJSON)
	}

	// собираем аттач
	if len(message.Attach) > 0 {

	}

	return h
}

func (b *Bot) request(
	ctx context.Context,
	method string,
	params map[string]string,
	sendTimeout time.Duration,
) (ans Response, err error) {
	if b.Token == "" {
		err = errors.New("no access token")
		return
	}

	if sendTimeout == 0 {
		if b.SendTimeout == 0 {
			sendTimeout = defaultSendTimeout
		} else {
			sendTimeout = b.SendTimeout
		}
	}

	var httpRetryCount, floodRetryCount int
	for {
		ans, err = b.fullRequest(ctx, method, params, sendTimeout)
		if err != nil {
			if b.httpErrorWait(httpRetryCount) {
				httpRetryCount++
				continue
			}
			return
		}

		// Проверяем ответ
		if ans.Error.ErrorCode != 0 {
			if ans.Error.ErrorMsg == "Too many requests per second" {
				// Ждем между запросами
				if b.floodWait(floodRetryCount) {
					floodRetryCount++
					continue
				}
			}

			err = errors.New(ans.Error.ErrorMsg)
			return
		}

		break
	}

	return
}

// Запрос к ВК
func (b *Bot) fullRequest(
	ctx context.Context,
	method string,
	params map[string]string,
	sendTimeout time.Duration,
) (ans Response, err error) {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}

	if params["v"] == "" {
		q.Add("v", APIVersion)
	}

	q.Add("access_token", b.Token)

	ctx, cancel := context.WithTimeout(ctx, sendTimeout)
	defer cancel()

	// Формируем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(APIMethodURL, b.Domain)+method, strings.NewReader(q.Encode()))
	if err != nil {
		return
	}

	// Отправляем запрос
	client := &http.Client{Transport: b.Transport}
	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return
	}

	// Если проблема с ответом
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Парсим ответ
	err = json.Unmarshal(content, &ans)
	if err != nil {
		return
	}

	return
}

// Ждем между запросами если вк ответил что запросы слишком частые
func (b *Bot) floodWait(retryCount int) bool {
	return b.retryWait(b.RetrySettings, retryCount)
}

// Попытка повтора запроса при ошибки http
func (b *Bot) httpErrorWait(retryCount int) bool {
	return b.retryWait(b.HTTPRetrySettings, retryCount)
}

// Ждем между запросами если вк ответил что запросы слишком частые
func (b *Bot) retryWait(settings map[int]int, count int) (ok bool) {

	counts := make([]int, 0, len(settings))
	for k := range settings {
		counts = append(counts, k)
	}

	sort.Ints(counts)

	var sleepTime int
	for _, c := range counts {
		if count < c {
			sleepTime = settings[c]
			break
		}
	}

	if sleepTime == 0 {
		return
	}

	// Ждем
	time.Sleep(time.Duration(sleepTime) * time.Second)

	ok = true
	return
}
