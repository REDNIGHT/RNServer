package RNServer

//import 	"github.com/json-iterator/go"
//var json = jsoniter.ConfigCompatibleWithStandardLibrary

import (
	"RNCore"
	"encoding/json"
	//"time"
)

type SocketBuffer2Router struct {
	RNCore.MinNode

	In chan *SocketBuffer

	out func(*Router)
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

//
func NewSocketBuffer2Router(name string) *SocketBuffer2Router {
	return &SocketBuffer2Router{RNCore.NewMinNode(name), make(chan *SocketBuffer, RNCore.InChanMinCount), nil}
}

func (this *SocketBuffer2Router) Out(out func(*Router)) {
	this.out = out
}

func (this *SocketBuffer2Router) Run() {
	go func() {
		for {
			select {
			case socketBuffer := <-this.In:
				josnData := &JosnData{}
				err := json.Unmarshal(socketBuffer.Buffer, josnData)
				if err != nil {
					this.Error("err != nil  err=" + err.Error())
					continue
				}

				this.out(&Router{socketBuffer.SocketID, josnData})
			}
		}
	}()
}
