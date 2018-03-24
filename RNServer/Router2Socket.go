package RNServer

import (
	"RNCore"
	//"time"
)

type Router2Socket struct {
	RNCore.MinNode

	In chan *Router

	Out func(*SocketBufferByName)
}

type RouterData struct {
	SocketID uintptr
	Json     []byte
}

func NewRouter2Socket(name string) *Router2Socket {
	return &Router2Socket{RNCore.NewMinNode(name), make(chan *Router, RNCore.InChanMinLen), nil}
}

func (this *Router2Socket) Run() {
	go func() {
		for {

			select {
			case router := <-this.In:
				buffer, err := json.Marshal(&RouterData{router.SocketID, router.JosnData.NJ})
				if err != nil {
					this.Error("err != nil  err=" + err.Error())
					continue
				}

				this.Out(&SocketBufferByName{router.JosnData.S, buffer})
			}
		}
	}()
}
