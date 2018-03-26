package RNCore

type IName interface {
	Name() string
	Type_Name() string
}

type IRun interface {
	Run()
}

type IClose interface {
	Close()
}

type IMessage interface {
	SendCall() chan<- func(IMessage)
	SendMessage(func(IMessage))
}
