package RNServer

import (
	"reflect"

	"../RNCore"
	//"time"
)

type RPC struct {
	RNCore.CallNode

	out reflect.Value
}

func NewRPC() *RPC {
	return &RPC{RNCore.NewCallNode(), reflect.Value{}}
}

func (this *RPC) In(content *Gate2RPCContent) {
	this.InCall() <- func() {
		this.onRPC(content)
	}
}
func (this *RPC) Out(out interface{}) {
	this.out = reflect.ValueOf(out).Elem()
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
