package RNCore

import (
	"reflect"
)

type CallNode struct {
	inCall chan func()
}

func NewCallNode() CallNode {
	return CallNode{make(chan func(), InChanMinLen)}
}

func (this *CallNode) InCall() chan<- func() {
	return this.inCall
}

func (this *CallNode) Run() {
	defer this.CatchPanic()

	for {
		f := <-this.inCall
		f()
	}
}

func (this *CallNode) CatchPanic(vs ...interface{}) {
	CatchPanic(this, vs)
}
func (this *CallNode) OnCatchPanic(v interface{}, node IPanic, vs ...interface{}) bool {
	return false
}
func (this *CallNode) OnPanicExit() {
}

//
func (this *CallNode) Name() string { return reflect.TypeOf(this).String() }
func (this *CallNode) Type_Name() string {
	return this.Name()
}

//
func (this *CallNode) Log(format string, a ...interface{}) {
	Print(this, format, a)
}
func (this *CallNode) Warn(format string, a ...interface{}) {
	Warn(this, format, a)
}
func (this *CallNode) Error(format string, a ...interface{}) {
	Error(this, format, a)
}
func (this *CallNode) Debug(format string, a ...interface{}) {
	Debug(this, format, a)
}
func (this *CallNode) Panic(v interface{}, format string, a ...interface{}) {
	Panic(this, v, format, a)
}
