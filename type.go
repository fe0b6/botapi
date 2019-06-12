package botapi

// Bot - Объект бота
type Bot struct {
	VK BotVK
	TG BotTG
}

// BotVK - Настройки ВК
type BotVK struct {
	Secret string
	Token  string
}

// BotTG - Настройки tg
type BotTG struct {
	Token string
}

// MessageOptions - Настройки сообщения
type MessageOptions struct {
	Source string
}
