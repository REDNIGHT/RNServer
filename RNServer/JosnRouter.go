package RNServer

//import 	"github.com/json-iterator/go"
//var json = jsoniter.ConfigCompatibleWithStandardLibrary

import (
	"RNCore"
	"encoding/json"
	//"time"
)

//import "unsafe"

type JosnRouter struct {
	RNCore.Node

	InSocketsBuffer chan *SocketBuffer

	outRouter chan<- *Router
}

type Router struct {
	SocketID uintptr
	JosnData *JosnData
}

//
type JosnData struct {
	S  string //server
	NJ []byte //next json
}

//IFilterName
func (this *JosnData) FilterName() string {
	return this.S
}

//
func NewJosnRouter(name string) *JosnRouter {
	return &JosnRouter{Node: RNCore.NewNode(name), InSocketsBuffer: make(chan *SocketBuffer, RNCore.InChanCount)}
}

func (this *JosnRouter) SetOut(outRouter chan<- *Router, node_chan_name string) {
	this.outRouter = outRouter

	//
	this.SetOutNodeInfos("outRouter", node_chan_name)
}

func (this *JosnRouter) Run() {

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case socketBuffer := <-this.InSocketsBuffer:
			josnData := &JosnData{}
			err := json.Unmarshal(socketBuffer.Buffer, josnData)
			if err != nil {
				this.Error("err != nil  err=" + err.Error())
				continue
			}

			this.rpc(josnData, socketBuffer.SocketID)

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

func (this *JosnRouter) rpc(josnData *JosnData, socketID uintptr) {
	//var id = uintptr(unsafe.Pointer(socket))
	this.outRouter <- &Router{socketID, josnData}
}

//
type _JosnRouterStateInfo struct {
	RNCore.StateInfo
	InCount uint
}

func (this *JosnRouter) OnStateInfo(counts ...*uint) RNCore.IStateInfo {
	return &_JosnRouterStateInfo{RNCore.StateInfo{this}, *counts[0]}
}
