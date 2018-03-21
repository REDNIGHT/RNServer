package RNCore

import (
	"fmt"
	//"time"
)

type INode interface {
	Name() string
	Init()
	Register()
	Run()
	State()
	Close()
	Destroy()

	Log(format string, a ...interface{})

	GetOutNodeInfos() []string
}

//
type Node struct {
	name string

	//
	StateSig chan bool
	//beginStateTime time.Time
	OutNodeInfos []string

	//
	CloseSig chan bool
}

func NewNode(name string) Node {
	return Node{name, make(chan bool), nil, make(chan bool)}
}

func (this *Node) Name() string { return this.name }

func (this *Node) Init()     {}
func (this *Node) Register() {}
func (this *Node) Run() {

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

func (this *Node) Close() {
	this.CloseSig <- true
	<-this.CloseSig
	close(this.CloseSig)
}
func (this *Node) Destroy() {
}

//
func (this *Node) State() {
	this.StateSig <- true
	<-this.StateSig
}
func (this *Node) OnState(counts ...*uint) {
	InState <- IState(this).OnStateInfo(counts...)

	//this.beginStateTime = time.Now()
	for i := 0; i < len(counts); i++ {
		*counts[i] = 0
	}
}
func (this *Node) OnStateInfo(counts ...*uint) *StateInfo {
	panic("//todo...  OnState")
	return NewStateInfo(this, 0)
}
func (this *Node) GetOutNodeInfos() []string {
	return this.OutNodeInfos
}
func (this *Node) SetOutNodeInfos(outName string, outNodeInfos ...string) {
	this.OutNodeInfos = outNodeInfos

	if len(this.OutNodeInfos) == 1 {
		this.OutNodeInfos[0] = fmt.Sprintf(outName+".%v", this.OutNodeInfos[0])
	} else {
		for i := 0; i < len(this.OutNodeInfos); i++ {
			this.OutNodeInfos[i] = fmt.Sprintf(outName+"%v.%v", i, this.OutNodeInfos[i])
		}
	}
}

//
func (this *Node) Log(format string, a ...interface{}) {
	Print(this, printLogLevel, format, a)
}
func (this *Node) Warn(format string, a ...interface{}) {
	Warn(this, printWarnLevel, format, a)
}
func (this *Node) Error(format string, a ...interface{}) {
	Error(this, printErrorLevel, format, a)
}
func (this *Node) Debug(format string, a ...interface{}) {
	Debug(this, printDebugLevel, format, a)
}
