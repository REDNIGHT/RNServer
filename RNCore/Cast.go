package RNCore

import (
	"time"
)

type Cast struct {
	Node

	In chan interface{}

	out func(interface{})
}

func NewCast(name string) *Cast {
	return &Cast{Node: NewNode(name), In: make(chan interface{}, InChanCount)}
}

func (this *Cast) SetOut(out func(interface{}), node_chan_name string) {
	this.out = out

	//
	this.SetOutNodeInfos("out", node_chan_name)
}

func (this *Cast) Run() {

	//
	for {
		//
		select {
		case in := <-this.In:
			this.out(in)

		case <-time.After(time.Second * StateTime):
			this.State()
			continue
		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

//
type _CastStateInfo struct {
	StateInfo
	In int
}

func (this *Cast) OnState() IStateInfo {
	return &_CastStateInfo{StateInfo{this}, len(this.In)}
}

//
func CastTo(In chan interface{}, Out func(in interface{})) {
	go func() {
		for {
			in := <-In
			Out(in)
		}
	}()
}
