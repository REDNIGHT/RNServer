package RNCore

import (
	"fmt"
	"time"
)

type INode interface {
	Name() string
	Init()
	Register()
	Run()
	Close()
	Destroy()

	Log(format string, a ...interface{})

	GetOutNodeInfos() []string
}

//
type Node struct {
	name string

	//
	OutNodeInfos []string

	//
	CloseSig chan bool
}

func NewNode(name string) Node {
	return Node{name, nil, make(chan bool, 1)}
}

func (this *Node) Name() string { return this.name }

func (this *Node) Init()     {}
func (this *Node) Register() {}
func (this *Node) Run() {

	for {
		//
		select {

		case <-time.After(time.Second * StateTime):
			this.State()
			continue

		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

func (this *Node) Close() {
	this.CloseSig <- true
}
func (this *Node) Destroy() {
	<-this.CloseSig
	close(this.CloseSig)
}

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

//
func (this *Node) State() {
	InState <- IState(this).OnState()
}
func (this *Node) OnState() IStateInfo {
	panic("//todo...  OnState")
	return &StateInfo{this}
}
func (this *Node) GetOutNodeInfos() []string {
	return this.OutNodeInfos
}

func (this *Node) SetOutNodeInfos(outName string, outNodeInfos ...string) {
	this.OutNodeInfos = outNodeInfos

	if len(this.OutNodeInfos) == 1 {
		this.OutNodeInfos[0] = fmt.Sprintf("out.%v", this.OutNodeInfos[0])
	} else {
		for i := 0; i < len(this.OutNodeInfos); i++ {
			this.OutNodeInfos[i] = fmt.Sprintf("out%v.%v", i, this.OutNodeInfos[i])
		}
	}
}
