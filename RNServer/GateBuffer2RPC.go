package RNServer

import (
	"RNCore"
	"encoding/json"
	//"time"
)

type GateBuffer2RPC struct {
	RNCore.Node

	InSocketsBuffer chan *SocketBuffer

	outGate2RPCContent chan<- *Gate2RPCContent
}

func NewGate2RPC(name string) *GateBuffer2RPC {
	return &GateBuffer2RPC{RNCore.NewNode(name), make(chan *SocketBuffer, RNCore.InChanCount), nil}
}

func (this *GateBuffer2RPC) SetOut(outGate2RPCContent chan<- *Gate2RPCContent, node_chan_name string) {
	this.outGate2RPCContent = outGate2RPCContent

	//
	this.SetOutNodeInfos("outGate2RPCContent", node_chan_name)
}

func (this *GateBuffer2RPC) Run() {

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {

		case socketBuffer := <-this.InSocketsBuffer:
			this.onSocketBuffer(socketBuffer)
			continue

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

type Gate2RPCContent struct {
	GateSocketID uintptr
	SocketID     uintptr
	Json         []byte
}

func (this *GateBuffer2RPC) onSocketBuffer(socketBuffer *SocketBuffer) {
	routerData := &RouterData{}
	json.Unmarshal(socketBuffer.Buffer, routerData)
	this.outGate2RPCContent <- &Gate2RPCContent{routerData.SocketID, socketBuffer.SocketID, routerData.Json}
}

//

func (this *GateBuffer2RPC) OnStateInfo(counts ...*uint) *RNCore.StateInfo {
	return RNCore.NewStateInfo(this, *counts[0])
}
