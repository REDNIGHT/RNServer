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

	messageChan chan func(node IMessage)

	InTotal uint
}

func NewNode(name string) Node {
	return Node{NewMinNode(name), make(chan func(IMessage)), 0}
}

//IMessage
func (this *Node) MessageChan() chan func(node IMessage) { return this.messageChan }
func (this *Node) SendMessage(f func(IMessage)) {
	mc := this.MessageChan()
	mc <- f
	<-mc
}
func (this *Node) OnMessage(f func(IMessage)) (close bool) {
	if f != nil {
		f(this)
		this.messageChan <- nil

		return false

	} else {

		//CloseSig
		this.messageChan <- nil

		return true
	}
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
	it := this.InTotal
	this.InTotal = 0
	return NewStateInfo(this, it)
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
