package RNCore

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
