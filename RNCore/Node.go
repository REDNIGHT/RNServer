package RNCore

import (
//"time"
)

type IMinNode interface {
	Name() string
	Run()

	Log(format string, a ...interface{})
}

type INode interface {
	Name() string
	Init()
	Register()
	Run()
	Close()
	Destroy()

	Log(format string, a ...interface{})

	State()
	DebugChanState()
}

//
type Node struct {
	MinNode

	//
	StateSig chan bool

	//
	CloseSig chan bool
}

func NewNode(name string) Node {
	return Node{NewMinNode(name), make(chan bool), make(chan bool)}
}

func (this *Node) Init()     {}
func (this *Node) Register() {}
func (this *Node) Run() {
	this.Panic("//todo... Run")

	var inCount uint = 0
	for {
		inCount++

		//
		select {

		case <-this.StateSig:
			this.OnState(&inCount)
			this.StateSig <- true
			continue

		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

//
func (this *Node) State() {
	this.StateSig <- true
	<-this.StateSig
}
func (this *Node) OnState(counts ...*uint) {
	InStateInfo() <- IState(this).OnStateInfo(counts...)

	//this.beginStateTime = time.Now()
	for i := 0; i < len(counts); i++ {
		*counts[i] = 0
	}
}
func (this *Node) OnStateInfo(counts ...*uint) *StateInfo {
	this.Panic("//todo...  OnState")
	return NewStateInfo(this, 0)
}

func (this *Node) DebugChanState() {
	this.Panic("//todo...  DebugChanState")
	this.OnDebugChanState("StateSig", len(this.StateSig))
	this.OnDebugChanState("CloseSig", len(this.CloseSig))
}
func (this *Node) OnDebugChanState(chanName string, chanLen int) {
	if chanLen > ChanOverloadLen {
		InChanOverload() <- &ChanOverload{this.Type_Name() + "." + chanName, chanLen}
	}
}

func (this *Node) Close() {
	this.CloseSig <- true
	<-this.CloseSig
	close(this.CloseSig)
}
func (this *Node) Destroy() {
}
