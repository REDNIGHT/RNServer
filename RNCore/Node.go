package RNCore

import (
//"time"
)

type INode interface {
	Name() string
	Init()
	Register()
	Run()
	Close()
	Destroy()

	Log(format string, a ...interface{})
}

//
type Node struct {
	MinNode
	MessageNode

	InTotal uint
}

func NewNode(name string) Node {
	return Node{NewMinNode(name), NewMessageNode(), 0}
}

//INode
func (this *Node) Init()     {}
func (this *Node) Register() {}
func (this *Node) Run() {
	this.Panic("//todo... Run")

	for {
		this.InTotal++

		//
		select {

		case f := <-this.messageChan:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

//
func (this *Node) Close() {
	this.messageChan <- nil
	<-this.messageChan
	close(this.messageChan)
}
func (this *Node) Destroy() {
}

//IState
func (this *Node) GetStateInfo() *StateInfo {
	inTotal := this.InTotal
	this.InTotal = 0
	return NewStateInfo(this, inTotal)
}

func (this *Node) DebugChanState(chanOverload chan *ChanOverload) {
	this.Panic("//todo...  DebugChanState")
	this.TestChanOverload(chanOverload, "messageChan", len(this.messageChan))
}
func (this *Node) TestChanOverload(chanOverload chan *ChanOverload, chanName string, chanLen int) {
	if chanLen > ChanOverloadLen {
		chanOverload <- &ChanOverload{this.Type_Name() + "." + chanName, chanLen}
	}
}
