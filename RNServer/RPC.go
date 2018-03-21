package RNServer

import (
	"RNCore"
	"encoding/json"
	"reflect"
	//"time"
)

type RPC struct {
	RNCore.Node

	InGate2RPCContent chan *Gate2RPCContent

	functions reflect.Value
}

func NewRPC(name string) *RPC {
	return &RPC{RNCore.NewNode(name), make(chan *Gate2RPCContent, RNCore.InChanCount), reflect.Value{}}
}

func (this *RPC) SetOut(functions interface{}, outNodeInfos ...string) {
	this.functions = reflect.ValueOf(functions).Elem()

	//
	this.SetOutNodeInfos("out", outNodeInfos...)
}

func (this *RPC) Run() {

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case content := <-this.InGate2RPCContent:
			this.onRPC(content)
			continue

			//
		case <-this.StateSig:
			this.OnState(&inCount)
			this.StateSig <- true
			continue

		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

type FuncContent struct {
	F  string
	NJ []byte
}

func (this *RPC) onRPC(content *Gate2RPCContent) {
	fc := &FuncContent{}
	json.Unmarshal(content.Json, fc)

	fun := this.functions.MethodByName(fc.F)
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

//
func (this *RPC) OnStateInfo(counts ...*uint) *RNCore.StateInfo {
	return RNCore.NewStateInfo(this, *counts[0])
}
