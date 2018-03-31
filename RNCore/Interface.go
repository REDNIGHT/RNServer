package RNCore

import "reflect"

type IName interface {
	Name() string
}

type IRun interface {
	Run()
}

type IClose interface {
	Close()
}

type ICall interface {
	InCall() chan<- func(ICall)
}

type IMessage interface {
	SendMessage(func(IMessage))
}

func Type_Name(node interface{}) string {
	if node == nil {
		return "nil"
	}
	tn := reflect.TypeOf(node).String()
	if in, b := node.(IName); b {
		tn += "." + in.Name()
	}
	return tn
}
