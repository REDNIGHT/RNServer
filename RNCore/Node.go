package RNCore

import (
	"reflect"
)

//
type Node struct {
	name string
}

func NewNode(name string) Node {
	return Node{name}
}

func (this *Node) Name() string      { return this.name }
func (this *Node) Type_Name() string { return reflect.TypeOf(this).String() + "." + this.Name() }

//
func (this *Node) Log(format string, a ...interface{}) {
	Print(this, format, a)
}
func (this *Node) Warn(format string, a ...interface{}) {
	Warn(this, format, a)
}
func (this *Node) Error(format string, a ...interface{}) {
	Error(this, format, a)
}
func (this *Node) Debug(format string, a ...interface{}) {
	Debug(this, format, a)
}
func (this *Node) Panic(format string, a ...interface{}) {
	Panic(this, format, a)
}

//IState
func (this *Node) SetOutNodeInfos(node_chan_Names ...string) {
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
