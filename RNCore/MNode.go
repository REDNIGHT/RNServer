package RNCore

import (
	"fmt"
	//"reflect"
)

//
type MNode struct {
	Node

	inCall    chan func(IMessage)
	inMessage chan func(IMessage)

	inTotal uint
}

func NewMNode(name string) MNode {
	return MNode{NewNode(name), make(chan func(IMessage)), make(chan func(IMessage), InChanLen), 0}
}

//
func (this *MNode) InCall() chan func(IMessage) {
	return this.inCall
}
func (this *MNode) InMessage() chan func(IMessage) {
	return this.inMessage
}

func (this *MNode) SendMessage(f func(IMessage)) {
	this.inMessage <- f
	<-this.inMessage
}
func (this *MNode) Run() {
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
func (this *MNode) OnMessage(f func(IMessage)) (close bool) {
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

//
func (this *MNode) Close() {
	this.SendMessage(nil)
	close(this.inMessage)
}

//IState
func (this *MNode) GetStateInfo() *StateInfo {
	inTotal := this.inTotal
	this.inTotal = 0
	return NewStateInfo(this, inTotal)
}

func (this *MNode) GetStateWarning(stateWarning func(name, warning string)) {
	this.TestChanOverload(stateWarning, "inCall", len(this.inCall))
	this.TestChanOverload(stateWarning, "inMessage", len(this.inMessage))
}
func (this *MNode) TestChanOverload(stateWarning func(name, warning string), chanName string, chanLen int) {
	if chanLen > ChanOverloadLen {
		stateWarning(this.Type_Name()+"."+chanName+".ChanOverload", fmt.Sprintf("%v", chanLen))
	}
}
