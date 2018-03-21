package RNCore

import (
//"time"
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

	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case in := <-this.In:
			for i := 0; i < len(this.outs); i++ {
				this.outs[i] <- in
			}

		case <-this.StateSig:
			this.OnState(&inCount)
			this.StateSig <- true
			continue
		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

//
func (this *Broadcast) OnStateInfo(counts ...*uint) *StateInfo {
	return NewStateInfo(this, *counts[0])
}
