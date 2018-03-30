package RNServer

//import 	"github.com/json-iterator/go"
//var json = jsoniter.ConfigCompatibleWithStandardLibrary

import (
	"../RNCore"
	//"time"
)

type SocketBuffer2Router struct {
	RNCore.CallNode

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
func NewSocketBuffer2Router() *SocketBuffer2Router {
	return &SocketBuffer2Router{RNCore.NewCallNode(), nil}
}

func (this *SocketBuffer2Router) In(socketBuffer *SocketBuffer) {
	this.InCall() <- func() {
		josnData := &JosnData{}
		err := json.Unmarshal(socketBuffer.Buffer, josnData)
		if err != nil {
			this.Error("err != nil  err=" + err.Error())
			return
		}

		this.out(&Router{socketBuffer.SocketId, josnData})
	}
}

func (this *SocketBuffer2Router) Out(out func(*Router)) {
	this.out = out
}
