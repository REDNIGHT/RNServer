package RNCore

import (
//"time"
)

//
type Node struct {
	MinNode
	MessageNode

	InTotal uint
}

func NewNode(name string) Node { return Node{NewMinNode(name), NewMessageNode(), 0} }

func (this *Node) Name() string      { return this.MinNode.Name() }
func (this *Node) Type_Name() string { return this.MinNode.Type_Name() }

//INode
func (this *Node) Run() {
	this.Panic("//todo... Run")

	for {
		this.InTotal++

		//
		select {

		case f := <-this.inMessage:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

//
func (this *Node) Close() {
	this.inMessage <- nil
	<-this.inMessage
	close(this.inMessage)
}

//IState
func (this *Node) GetStateInfo() *StateInfo {
	inTotal := this.InTotal
	this.InTotal = 0
	return NewStateInfo(this, inTotal)
}

func (this *Node) DebugChanState(chanOverload chan *ChanOverload) {
	this.TestChanOverload(chanOverload, "inMessage", len(this.inMessage))
}
func (this *Node) TestChanOverload(chanOverload chan *ChanOverload, chanName string, chanLen int) {
	if chanLen > ChanOverloadLen {
		chanOverload <- &ChanOverload{this.Type_Name() + "." + chanName, chanLen}
	}
}
