package RNServer

import (
	"RNCore"
	//"time"
)

type Router2Socket struct {
	RNCore.CallNode

	Out func(s string, buffer []byte)
}

type RouterData struct {
	SocketID uintptr
	Json     []byte
}

func NewRouter2Socket() *Router2Socket {
	return &Router2Socket{RNCore.NewCallNode(), nil}
}

func (this *Router2Socket) In(socketID uintptr, josnData *JosnData) {
	this.InCall() <- func() {
		buffer, err := json.Marshal(&RouterData{socketID, josnData.NJ})
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			return
		}

		this.Out(josnData.S, buffer)
	}
}
