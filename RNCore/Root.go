package RNCore

import (
	"os"
	"os/signal"
)

type root struct {
	MNode
	ns []IName
}

//
var _root *root

func NewRoot(serverName string) *root {
	this := &root{NewMNode(serverName), make([]IName, 0)}
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
func (this *root) Add(ns ...IName) {
	for _, n := range ns {
		if this.Get(n.Name()) != nil {
			this.Panic(nil, "Get(n.Name()) != nil  n.Name()="+n.Name())
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
func (this *root) Broadcast(f func(IMessage)) {
	this.InCall() <- f

	for i := 0; i < len(this.ns); i++ {
		im, b := this.ns[i].(IMessage)
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
	defer this.CatchPanic()

	//
	for i := 0; i < len(this.ns); i++ {

		n, b := this.ns[i].(IRun)
		if b == true {
			go n.Run()
			this.Log("%v.Run()", n)
		}
	}

	//
	this.MNode.Run()

	go func() {
		// close
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		sig := <-c
		this.Log("closing down (signal: %v)", sig)

		this.Close()
	}()
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
