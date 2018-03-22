package RNCore

import (
	"fmt"
	"reflect"
)

type IMinNode interface {
	Name() string
	Run()

	Log(format string, a ...interface{})
}

//
type MinNode struct {
	name string
}

func NewMinNode(name string) MinNode {
	return MinNode{name}
}

func (this *MinNode) Name() string { return this.name }

func (this *MinNode) Type_Name() string {
	return reflect.TypeOf(this).String() + "." + this.Name()
}

//
func (this *MinNode) SetOutNodeInfo(outNames ...string) {
}

func (this *Node) SetOutNodeInfos(node_chan_Names ...string) {
	nodeInfos := &NodeInfo{this.Type_Name(), node_chan_Names}

	/*if len(nodeInfos.OutNames) == 1 {
		nodeInfos.OutNames[0] = fmt.Sprintf(nodeName+".%v", nodeInfos.OutNames[0])
	} else {
		for i := 0; i < len(nodeInfos.OutNames); i++ {
			nodeInfos.OutNames[i] = fmt.Sprintf(nodeName+".%v", nodeInfos.OutNames[i])
		}
	}*/

	InNodeInfo() <- nodeInfos
}

//
func (this *MinNode) Log(format string, a ...interface{}) {
	Print(this, printLogLevel, format, a)
}
func (this *MinNode) Warn(format string, a ...interface{}) {
	Warn(this, printWarnLevel, format, a)
}
func (this *MinNode) Error(format string, a ...interface{}) {
	Error(this, printErrorLevel, format, a)
}
func (this *MinNode) Debug(format string, a ...interface{}) {
	Debug(this, printDebugLevel, format, a)
}
func (this *MinNode) Panic(format string, a ...interface{}) {
	panic(this.Type_Name() + fmt.Sprintf(format, a...))
}
