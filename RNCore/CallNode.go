package RNCore

type CallNode struct {
	inCall chan func(ICall)
}

func NewCallNode() CallNode {
	return CallNode{make(chan func(ICall), InChanMinLen)}
}

func (this *CallNode) InCall() chan<- func(ICall) {
	return this.inCall
}

func (this *CallNode) Run() {
	defer this.CatchPanic(func(v interface{}) bool {
		if RNCDebug {
			return false
		}
		go this.Run()
		return true
	})

	for {
		f := <-this.inCall
		f(this)
	}
}

func (this *CallNode) CatchPanic(onCatchPanic func(interface{}) bool) {
	CatchPanic(onCatchPanic, this)
}

//
//func (this *CallNode) Name() string { return reflect.TypeOf(this).String() }
//func (this *CallNode) Type_Name() string { return this.Name() }

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
