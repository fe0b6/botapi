package vk

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	// DefaultMaxIdleConnsPerHost - Сколько максимум конектов в ожидании может быть на 1 хост
	DefaultMaxIdleConnsPerHost = 20
	// DefaultIdleConnTimeout - Время бездействия конекта перед отключением
	DefaultIdleConnTimeout = 10 * time.Minute
)

// Builder - сборщик бота
type Builder struct {
	Token               string
	Host                string
	SendTimeout         time.Duration
	UsePrometheus       bool
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	Domain              string
	RetrySettings       map[int]int
	HTTPRetrySettings   map[int]int
}

// NewBot - Создаем нового бота
func NewBot() *Builder {
	return &Builder{}
}

// WithToken - Добавляем токен
func (b *Builder) WithToken(token string) *Builder {
	b.Token = token
	return b
}

// WithHost - Добавляем хост
func (b *Builder) WithHost(host string) *Builder {
	b.Host = host
	return b
}

// WithSendTimeout - Добавляем таймаут отправки сообщения
func (b *Builder) WithSendTimeout(timeout time.Duration) *Builder {
	b.SendTimeout = timeout
	return b
}

// WithPrometheus - Включаем передачу статистики в прометей
func (b *Builder) WithPrometheus(prom bool) *Builder {
	b.UsePrometheus = prom
	return b
}

// WithDomain - заменяем стандартный домен
func (b *Builder) WithDomain(domain string) *Builder {
	b.Domain = domain
	return b
}

// WithRetrySettings - заменяем стандартный домен
func (b *Builder) WithRetrySettings(settings map[int]int) *Builder {
	b.RetrySettings = settings
	return b
}

// WithHTTPRetrySettings - заменяем стандартный домен
func (b *Builder) WithHTTPRetrySettings(settings map[int]int) *Builder {
	b.HTTPRetrySettings = settings
	return b
}

// WithMaxIdleConnsPerHost - Указываем сколько конектов можно к одному хосту
func (b *Builder) WithMaxIdleConnsPerHost(count int) *Builder {
	b.MaxIdleConnsPerHost = count
	return b
}

// WithIdleConnTimeout - Указываем сколько времени можно висеть в ожидании
func (b *Builder) WithIdleConnTimeout(timeout time.Duration) *Builder {
	b.IdleConnTimeout = timeout
	return b
}

// Build - собираем бота
func (b *Builder) Build() *Bot {

	if b.Domain == "" {
		b.Domain = DefaultDomain
	}

	if len(b.RetrySettings) == 0 {
		b.RetrySettings = defaultRetrySettings
	}

	if len(b.HTTPRetrySettings) == 0 {
		b.HTTPRetrySettings = defaultHTTPRetrySettings
	}

	bot := &Bot{
		Host:              b.Host,
		Token:             b.Token,
		SendTimeout:       b.SendTimeout,
		UsePrometheus:     b.UsePrometheus, // TODO
		Rand:              rand.New(rand.NewSource(time.Now().UnixNano())),
		RetrySettings:     b.RetrySettings,
		HTTPRetrySettings: b.HTTPRetrySettings,
		Domain:            b.Domain,
	}

	if b.MaxIdleConnsPerHost > 0 && b.IdleConnTimeout > 0 {
		bot.Transport = &http.Transport{
			MaxIdleConnsPerHost: b.MaxIdleConnsPerHost,
			IdleConnTimeout:     b.IdleConnTimeout,
		}
	} else {
		bot.Transport = defaultHTTPTransport
	}

	return bot
}
