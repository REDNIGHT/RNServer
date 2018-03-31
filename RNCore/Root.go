package RNCore

import (
	"os"
	"os/signal"
)

type root struct {
	MNode
	ns []interface{}
}

//
var _root *root

func NewRoot(serverName string) *root {
	this := &root{NewMNode(serverName), make([]interface{}, 0)}
	if _root != nil {
		this.Panic(nil, "_root != nil")
	}
	_root = this
	return _root
}
func Root() *root {
	return _root
}

//
func getName(n interface{}) string {
	if in, b := n.(IName); b {
		return in.Name()
	} else {
		return Type_Name(n)
	}
}
func (this *root) Add(ns ...interface{}) {
	for _, n := range ns {
		name := getName(n)
		if this.Get(name) != nil {
			this.Panic(nil, "Get(name) != nil  name="+name)
		}

		this.ns = append(this.ns, n)
	}
}
func (this *root) Get(name string) interface{} {
	for _, n := range this.ns {
		node_name := getName(n)
		if node_name == name {
			return n
		}
	}

	return nil
}

func (this *root) GetCount() int {
	return len(this.ns)
}
func (this *root) GetByIndex(index int) interface{} {
	return this.ns[index]
}
func (this *root) ForEach(f func(interface{})) {
	f(this)

	for i := 0; i < len(this.ns); i++ {
		f(this.ns[i])
	}
}
func (this *root) Broadcast(f func(ICall)) {
	this.InCall() <- f

	for i := 0; i < len(this.ns); i++ {
		im, b := this.ns[i].(ICall)
		if b == true {
			im.InCall() <- f
		}
	}
}

func (this *root) BroadcastMessage(f func(IMessage)) {
	this.SendMessage(f)

	for i := 0; i < len(this.ns); i++ {
		im, b := this.ns[i].(IMessage)
		if b == true {
			im.SendMessage(f)
		}
	}
}

func (this *root) Run() {
	//defer this.CatchPanic()

	//
	for i := 0; i < len(this.ns); i++ {

		n, b := this.ns[i].(IRun)
		if b == true {
			go n.Run()
			this.Log("%v.Run()", n)
		}
	}

	go func() {
		// close
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		sig := <-c
		this.Log("closing down (signal: %v)", sig)

		this.Close()
	}()

	//
	this.MNode.Run()
}

func (this *root) Close() {
	for i := len(this.ns) - 1; i >= 0; i-- {

		ic, b := this.ns[i].(IClose)
		if b == true {
			ic.Close()
		}
		this.Log("%v.Close()", ic)
	}

	this.MNode.Close()
}
