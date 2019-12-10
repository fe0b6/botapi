package botapi

// Bot - Объект бота
type Bot interface {
	SendMessage(OutMessage) Answer
	NewMessage() OutMessage
}
