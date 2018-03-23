package RNCore

type IMessage interface {
	MessageChan() chan func(IMessage)
	SendMessage(func(IMessage))
}
