package RNServer

import (
	"RNCore"
	//"time"
)

type GateSockerBuffer2RPC struct {
	In chan *SocketBuffer

	Out func(*Gate2RPCContent)
}

type Gate2RPCContent struct {
	GateSocketID uintptr
	SocketID     uintptr
	Json         []byte
}

func NewGate2RPC() *GateSockerBuffer2RPC {
	return &GateSockerBuffer2RPC{make(chan *SocketBuffer, RNCore.InChanMinLen), nil}
}

func (this *GateSockerBuffer2RPC) Go() {
	go func() {
		for {
			//
			select {
			case socketBuffer := <-this.In:
				routerData := &RouterData{}
				json.Unmarshal(socketBuffer.Buffer, routerData)
				this.Out(&Gate2RPCContent{routerData.SocketID, socketBuffer.SocketID, routerData.Json})
				continue
			}
		}
	}()
}
