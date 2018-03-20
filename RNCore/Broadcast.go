package RNCore

import (
	"time"
)

type Broadcast struct {
	Node

	In chan interface{}

	outs []chan<- interface{}
}

func NewBroadcast(name string) *Broadcast {
	return &Broadcast{Node: NewNode(name), In: make(chan interface{}, InChanCount)}
}

func (this *Broadcast) SetOut(outs []chan<- interface{}, outNodeInfos ...string) {
	this.outs = outs

	//
	this.SetOutNodeInfos("out", outNodeInfos...)
}

func (this *Broadcast) Run() {

	//
	for {
		//
		select {
		case in := <-this.In:
			for i := 0; i < len(this.outs); i++ {
				this.outs[i] <- in
			}

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
type _BroadcastStateInfo struct {
	StateInfo
	In int
}

func (this *Broadcast) OnState() IStateInfo {
	return &_BroadcastStateInfo{StateInfo{this}, len(this.In)}
}
