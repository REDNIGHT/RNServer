package RNCore

import (
//"time"
)

type Cast struct {
	In chan interface{}
}

func NewCast() *Cast {
	return &Cast{In: make(chan interface{})}
}

//
func (this *Cast) To(out func(in interface{})) {
	go func() {
		for {
			in := <-this.In
			out(in)
		}
	}()
}
