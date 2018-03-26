package RNCore

import (
	"fmt"
	"reflect"
)

//
type MessageNode struct {
	inCall    chan func(IMessage)
	inMessage chan func(IMessage)

	inTotal uint
}

func NewMessageNode() MessageNode {
	return MessageNode{make(chan func(IMessage)), make(chan func(IMessage), InChanLen), 0}
}

func (this *MessageNode) InCall() chan func(IMessage) {
	return this.inCall
}
func (this *MessageNode) InMessage() chan func(IMessage) {
	return this.inMessage
}

//
func (this *MessageNode) SendCall() chan<- func(IMessage) {
	return this.inCall
}

func (this *MessageNode) SendMessage(f func(IMessage)) {
	this.inMessage <- f
	<-this.inMessage
}
func (this *MessageNode) OnMessage(f func(IMessage)) (close bool) {
	if f != nil {
		f(this)
		this.inMessage <- nil

		return false

	} else {

		this.inMessage <- nil

		//CloseSig
		return true
	}
}
func (this *MessageNode) Run() {
	for {
		this.inTotal++

		//
		select {
		case f := <-this.inCall:
			f(this)

		case f := <-this.inMessage:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

//
func (this *Node) Close() {
	this.SendMessage(nil)
	close(this.inMessage)
}

//
func (this *MessageNode) Name() string { return reflect.TypeOf(this).String() }

func (this *MessageNode) Type_Name() string {
	return this.Name()
}

func (this *MessageNode) SetOutNodeInfos(node_chan_Names ...string) {
	nodeInfos := &NodeInfo{this.Type_Name(), node_chan_Names}

	/*if len(nodeInfos.OutNames) == 1 {
		nodeInfos.OutNames[0] = fmt.Sprintf(nodeName+".%v", nodeInfos.OutNames[0])
	} else {
		for i := 0; i < len(nodeInfos.OutNames); i++ {
			nodeInfos.OutNames[i] = fmt.Sprintf(nodeName+".%v", nodeInfos.OutNames[i])
		}
	}*/

	AddNodeInfo(nodeInfos)
}

//IState
func (this *Node) GetStateInfo() *StateInfo {
	inTotal := this.inTotal
	this.inTotal = 0
	return NewStateInfo(this, inTotal)
}

func (this *Node) GetStateWarning(stateWarning func(name, warning string)) {
	this.TestChanOverload(stateWarning, "inMessage", len(this.inMessage))
}
func (this *Node) TestChanOverload(stateWarning func(name, warning string), chanName string, chanLen int) {
	if chanLen > ChanOverloadLen {
		stateWarning(this.Type_Name()+"."+chanName+".ChanOverload", fmt.Sprintf("%v", chanLen))
	}
}
