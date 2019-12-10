package vk

import (
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

	if b.MaxIdleConnsPerHost == 0 {
		b.MaxIdleConnsPerHost = DefaultMaxIdleConnsPerHost
	}

	if b.IdleConnTimeout == 0 {
		b.IdleConnTimeout = DefaultIdleConnTimeout
	}

	bot := &Bot{
		Host:          b.Host,
		Token:         b.Token,
		SendTimeout:   b.SendTimeout,
		UsePrometheus: b.UsePrometheus, // TODO
		Transport: &http.Transport{
			MaxIdleConnsPerHost: b.MaxIdleConnsPerHost,
			IdleConnTimeout:     b.IdleConnTimeout,
		},
	}

	return bot
}
