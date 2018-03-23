package RNCore

import (
	"os"
	"os/signal"
)

type root struct {
	Node
	minNodes []IMinNode
	nodes    []INode
}

//
var _root *root

func NewRoot(serverName string) *root {
	if _root != nil {
		panic("_root != nil")
	}
	_root = &root{NewNode(serverName), make([]IMinNode, 0), make([]INode, 0)}
	return _root
}
func Root() *root {
	return _root
}

//
func (this *root) Add(nodes ...IMinNode) {
	for _, node := range nodes {
		if this.Get(node.Name()) != nil {
			this.Panic("Get(node.Name()) != nil  node.Name()=" + node.Name())
		}

		this.minNodes = append(this.minNodes, node)

		if node.(INode) != nil {
			this.nodes = append(this.nodes, node.(INode))
		}
	}
}
func (this *root) Get(name string) IMinNode {
	for i := 0; i < len(this.minNodes); i++ {
		if this.minNodes[i].Name() == name {
			return this.minNodes[i]
		}
	}

	return nil
}
func (this *root) GetCount() int {
	return len(this.minNodes)
}
func (this *root) GetByIndex(index int) IMinNode {
	return this.minNodes[index]
}
func (this *root) ForEach(f func(IMinNode)) {
	f(this)

	for i := 0; i < len(this.minNodes); i++ {
		f(this.minNodes[i])
	}
}
func (this *root) BroadcastMessage(f func(IMessage)) {
	this.ForEach(func(node IMinNode) {
		if im, b := node.(IMessage); b == true {
			im.SendMessage(f)
		}
	})
}

//
func (this *root) Init() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Init()

		n.Log("Init()")
	}
}

func (this *root) Register() {
	for i := 0; i < len(this.nodes); i++ {

		n := this.nodes[i]
		n.Register()

		n.Log("Register()")
	}
}

func (this *root) Run() {
	for i := 0; i < len(this.minNodes); i++ {

		n := this.minNodes[i]
		go n.Run()

		n.Log("Run()")
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

		n.Log("Close()")
	}
}

func (this *root) Destroy() {
	for i := len(this.nodes) - 1; i >= 0; i-- {

		n := this.nodes[i]
		n.Destroy()

		n.Log("Destroy()")
	}
}
