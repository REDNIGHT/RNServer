package RNCore

import (
//"time"
)

//
func CastTo(in func() interface{}, out func(in interface{})) {
	go func() {
		for {
			out(in())
		}
	}()
}

//
func CastTo2(inOut func()) {
	go func() {
		for {
			inOut()
		}
	}()
}

//
type Cast struct {
	In  func() interface{}
	Out func(interface{})
}

func (this *Cast) Go() {
	go func() {
		this.Out(this.In())
	}()
}

type AllInOne Cast
