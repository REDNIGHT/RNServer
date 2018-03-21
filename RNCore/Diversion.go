package RNCore

import (
//"time"
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
	var inCount uint = 0
	for {
		inCount++
		//
		select {
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
type _DiversionStateInfo struct {
	StateInfo
	InCount uint
}

func (this *Diversion) OnStateInfo(counts ...*uint) IStateInfo {
	return &_DiversionStateInfo{StateInfo{this}, *counts[0]}
}
