package RNServer

import (
	"RNCore"
	"encoding/json"
	//"time"
)

type Router2Socket struct {
	RNCore.Node

	InRouter chan *Router

	outSendBufferByName chan<- *SocketBufferByName
}

type RouterData struct {
	SocketID uintptr
	Json     []byte
}

func NewRouter2Socket(name string) *Router2Socket {
	return &Router2Socket{Node: RNCore.NewNode(name), InRouter: make(chan *Router, RNCore.InChanCount)}
}

func (this *Router2Socket) SetOut(outSendBufferByName chan<- *SocketBufferByName, node_chan_name string) {
	this.outSendBufferByName = outSendBufferByName

	//
	this.SetOutNodeInfos("outSendBufferByName", node_chan_name)
}

func (this *Router2Socket) Run() {

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case router := <-this.InRouter:
			buffer, err := json.Marshal(&RouterData{router.SocketID, router.JosnData.NJ})
			if err != nil {
				this.Error("err != nil  err=" + err.Error())
				continue
			}

			this.outSendBufferByName <- &SocketBufferByName{router.JosnData.S, buffer}

			//
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
type _Router2SocketStateInfo struct {
	RNCore.StateInfo
	InCount uint
}

func (this *Router2Socket) OnStateInfo(counts ...*uint) RNCore.IStateInfo {
	return &_Router2SocketStateInfo{RNCore.StateInfo{this}, *counts[0]}
}