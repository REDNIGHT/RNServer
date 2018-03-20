package RNServer

import (
	"RNCore"
	"encoding/json"
	"time"
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
	for {
		//
		select {

		case socketBuffer := <-this.InSocketsBuffer:
			this.onSocketBuffer(socketBuffer)
			continue

		case <-time.After(time.Second * RNCore.StateTime):
			this.State()
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
type _GateBuffer2RPCStateInfo struct {
	RNCore.StateInfo
	InSocketsBuffer int
}

func (this *GateBuffer2RPC) OnState() RNCore.IStateInfo {
	return &_GateBuffer2RPCStateInfo{RNCore.StateInfo{this}, len(this.InSocketsBuffer)}
}
