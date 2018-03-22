package RNCore

import (
//"time"
)

type Broadcast struct {
	In func() interface{}

	outs []func(interface{})
}

func NewBroadcast() *Broadcast {
	return &Broadcast{nil, make([]func(interface{}), InChanMinCount)}
}

func (this *Broadcast) OutAdd(outs ...func(interface{})) {
	this.outs = append(this.outs, outs...)
}

func (this *Broadcast) Go() {
	go func() {
		for {
			v := this.In()
			for i := 0; i < len(this.outs); i++ {
				this.outs[i](v)
			}
		}
	}()
}
