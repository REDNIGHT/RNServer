package RNCore

import (
	"reflect"
)

type IMessage interface {
	MessageChan() chan func(IMessage)
	SendMessage(func(IMessage))
}

//
type MessageNode struct {
	messageChan chan func(node IMessage)
}

func NewMessageNode() MessageNode {
	return MessageNode{make(chan func(IMessage))}
}

//IMessage
func (this *MessageNode) MessageChan() chan func(node IMessage) { return this.messageChan }
func (this *MessageNode) SendMessage(f func(IMessage)) {
	mc := this.MessageChan()
	mc <- f
	<-mc
}
func (this *MessageNode) OnMessage(f func(IMessage)) (close bool) {
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
func (this *MessageNode) Run() {
	for {
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

	InNodeInfo() <- nodeInfos
}
