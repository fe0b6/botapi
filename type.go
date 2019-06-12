package botapi

// Bot - Объект бота
type Bot struct {
	VK BotVK
}

// BotVK - Настройки ВК
type BotVK struct {
	Secret string
	Token  string
}

// MessageOptions - Настройки сообщения
type MessageOptions struct {
	Source string
}
