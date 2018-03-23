package RNCore

import (
	"os"
	"os/signal"
)

type root struct {
	Node
	minNodes     []IMinNode
	messageNodes []IMessage
	nodes        []INode
	ns           []IName
}

//
var _root *root

func NewRoot(serverName string) *root {
	if _root != nil {
		panic("_root != nil")
	}
	_root = &root{NewNode(serverName), make([]IMinNode, 0), make([]IMessage, 0), make([]INode, 0), make([]IName, 0)}
	return _root
}
func Root() *root {
	return _root
}

//
func (this *root) Add(ns ...IName) {
	for _, n := range ns {
		if this.Get(n.Name()) != nil {
			this.Panic("Get(n.Name()) != nil  n.Name()=" + n.Name())
		}

		if imn, b := n.(IMinNode); b == true {
			this.minNodes = append(this.minNodes, imn)
		}

		if im, b := n.(IMessage); b == true {
			this.messageNodes = append(this.messageNodes, im)
		}

		if n, b := n.(INode); b == true {
			this.nodes = append(this.nodes, n)
		}

		this.ns = append(this.ns, n)
	}
}
func (this *root) Get(name string) IName {
	for i := 0; i < len(this.ns); i++ {
		if this.ns[i].Name() == name {
			return this.ns[i]
		}
	}

	return nil
}

func (this *root) GetCount() int {
	return len(this.ns)
}
func (this *root) GetByIndex(index int) IName {
	return this.ns[index]
}
func (this *root) ForEach(f func(IName)) {
	f(this)

	for i := 0; i < len(this.ns); i++ {
		f(this.ns[i])
	}
}
func (this *root) BroadcastMessage(f func(IMessage)) {
	this.SendMessage(f)

	for i := 0; i < len(this.messageNodes); i++ {
		this.messageNodes[i].SendMessage(f)
	}
}

//
func (this *root) Init() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Init()

		this.Log("%v.Init()", n)
	}
}

func (this *root) Register() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Register()

		this.Log("%v.Register()", n)
	}
}

func (this *root) Run() {
	for i := 0; i < len(this.ns); i++ {

		n, b := this.ns[i].(IRun)
		if b {
			go n.Run()
			this.Log("%v.Run()", n)
		}
	}

	//
	go func() {
		// close
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		sig := <-c
		this.Log("closing down (signal: %v)", sig)

		this.Node.Close() //这行代码可以退出下面的for循环
	}()

	//
	for {
		select {
		case f := <-this.messageChan:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

func (this *root) Close() {
	for i := len(this.nodes) - 1; i >= 0; i-- {

		n := this.nodes[i]
		n.Close()

		this.Log("%v.Close()", n)
	}
}

func (this *root) Destroy() {
	for i := len(this.nodes) - 1; i >= 0; i-- {

		n := this.nodes[i]
		n.Destroy()

		this.Log("%v.Destroy()", n)
	}
}
