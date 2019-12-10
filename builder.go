package botapi

import "time"

// Builder - сборщик
type Builder interface {
	WithToken(string) Builder
	WithHost(string) Builder
	WithSendTimeout(time.Duration) Builder
	WithMaxIdleConnsPerHost(time.Duration) Builder
	WithIdleConnTimeout(time.Duration) Builder
	WithPrometheus(bool) Builder
	Build() Bot
}
