package RNCore

//"time"

type Broadcast struct {
	In func() interface{}

	outs []func(interface{})
}

func NewBroadcast() *Broadcast {
	return &Broadcast{nil, make([]func(interface{}), InChanMinLen)}
}

func (this *Broadcast) OutAdd(outs ...func(interface{})) {
	this.outs = append(this.outs, outs...)
}

func (this *Broadcast) Run() {
	defer CatchPanic(func(v interface{}) bool {
		if RNCDebug {
			return false
		}
		go this.Run()
		return true
	}, this)

	for {
		v := this.In()
		for i := 0; i < len(this.outs); i++ {
			this.outs[i](v)
		}
	}
}
