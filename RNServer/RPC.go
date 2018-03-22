package RNServer

import (
	"RNCore"
	"encoding/json"
	"reflect"
	//"time"
)

type RPC struct {
	RNCore.MinNode

	In chan *Gate2RPCContent

	out reflect.Value
}

func NewRPC(name string) *RPC {
	return &RPC{RNCore.NewMinNode(name), make(chan *Gate2RPCContent, RNCore.InChanCount), reflect.Value{}}
}

func (this *RPC) Out(out interface{}) {
	this.out = reflect.ValueOf(out).Elem()
}

func (this *RPC) Go() {
	go func() {
		for {
			select {
			case content := <-this.In:
				this.onRPC(content)
				continue
			}
		}
	}()
}

type FuncContent struct {
	F  string
	NJ []byte
}

func (this *RPC) onRPC(content *Gate2RPCContent) {
	fc := &FuncContent{}
	json.Unmarshal(content.Json, fc)

	fun := this.out.MethodByName(fc.F)
	v_nil := reflect.Value{}
	if fun == v_nil {
		this.Error("fun == v_nil fc.F=" + fc.F)
	}

	//
	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(fc.NJ)
	params[1] = reflect.ValueOf(content)
	fun.Call(params)
}
