package RNCore

//"time"

//分流
type Diversion struct {
	In func() interface{}

	outs []func(interface{})
}

func NewDiversion() *Diversion {
	return &Diversion{nil, make([]func(interface{}), 0)}
}

func (this *Diversion) OutAdd(outs ...func(interface{})) {
	this.outs = append(this.outs, outs...)
}

func (this *Diversion) Run() {
	defer CatchPanic(func(v interface{}) bool {
		if RNCDebug {
			return false
		}
		go this.Run()
		return true
	}, this)

	for i := 0; i < len(this.outs); i++ {
		go func() {
			this.outs[i](this.In())
		}()
	}
}
