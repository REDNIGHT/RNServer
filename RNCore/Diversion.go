package RNCore

import (
	"time"
)

type Diversion struct {
	Node

	In chan interface{}

	outs []chan<- interface{}
}

func NewDiversion(name string) *Diversion {
	return &Diversion{Node: NewNode(name), In: make(chan interface{}, InChanCount)}
}

func (this *Diversion) SetOut(outs []chan<- interface{}, outNodeInfos ...string) {
	this.outs = outs

	//
	this.SetOutNodeInfos("out", outNodeInfos...)
}

func (this *Diversion) Run() {

	for i := 0; i < len(this.outs); i++ {
		go func() {
			this.outs[i] <- this.In
		}()
	}

	//
	for {
		//
		select {
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
type _DiversionStateInfo struct {
	StateInfo
	In int
}

func (this *Diversion) OnState() IStateInfo {
	return &_DiversionStateInfo{StateInfo{this}, len(this.In)}
}
