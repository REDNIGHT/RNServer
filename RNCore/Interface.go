package RNCore

type IName interface {
	Name() string
	Type_Name() string
}

type IPanic interface {
	OnCatchPanic(v interface{}, node IPanic, vs ...interface{}) bool
	OnPanicExit()
}

type IRun interface {
	Run()
}

type IClose interface {
	Close()
}

type IMessage interface {
	InCall() chan func(IMessage)
	SendMessage(func(IMessage))
}
