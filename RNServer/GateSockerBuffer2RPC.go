package RNServer

import (
	"../RNCore"
	//"time"
)

type GateSockerBuffer2RPC struct {
	RNCore.CallNode

	Out func(*Gate2RPCContent)
}

type Gate2RPCContent struct {
	GateSocketId uintptr
	SocketId     uintptr
	Json         []byte
}

func NewGateSockerBuffer2RPC() *GateSockerBuffer2RPC {
	return &GateSockerBuffer2RPC{RNCore.NewCallNode(), nil}
}

func (this *GateSockerBuffer2RPC) In(socketId uintptr, buffer []byte) {
	this.InCall() <- func() {
		routerData := &RouterData{}
		json.Unmarshal(buffer, routerData)
		this.Out(&Gate2RPCContent{routerData.SocketID, socketId, routerData.Json})
	}
}
