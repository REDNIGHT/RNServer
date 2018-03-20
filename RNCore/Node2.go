package RNCore

import (
	"reflect"
	"time"
)

type Node2 struct {
	Node

	OnCloseSig func(i interface{})

	selectCases []reflect.SelectCase
	funcs       []func(i interface{})
}

func NewNode2(name string) Node2 {
	return Node2{NewNode(name), nil, make([]reflect.SelectCase, 0), make([]func(i interface{}), 0)}
}

func (this *Node2) AddIn(in interface{}, onFunc func(i interface{})) {
	this.selectCases = append(this.selectCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(in)})
	this.funcs = append(this.funcs, onFunc)
}

func (this *Node2) Run() {

	if this.OnCloseSig == nil {
		this.OnCloseSig = this.closeSigHandle
	}
	this.AddIn(this.CloseSig, this.closeSigHandle)

	this.AddIn(nil, nil)
	this.setNextState()

	for {
		chosen, recv, recvOk := reflect.Select(this.selectCases)

		if recvOk {
			this.funcs[chosen](recv.Interface())
		}
	}
}

func (this *Node2) closeSigHandle(v interface{}) {
	this.CloseSig <- true
}

func (this *Node2) setNextState() {
	index := len(this.funcs) - 1
	this.selectCases[index] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(time.After(time.Second * StateTime))}
	this.funcs[index] = this.stateHandle
}
func (this *Node2) stateHandle(v interface{}) {
	this.setNextState()

	this.State()
}

//
type _Node2StateInfo struct {
	StateInfo
	InNames []string
	InLens  []int
}

func (this *_Node2StateInfo) GetInNames() []string {
	return this.InNames
}
func (this *_Node2StateInfo) GetInLens() []int {
	return this.InLens
}

func (this *Node2) OnState() IStateInfo {
	l := len(this.selectCases) - 2 //去掉CloseSig,State
	inNames := make([]string, l)
	inLens := make([]int, l)
	for i := 0; i < l; i++ {
		inNames[i] = reflect.TypeOf(this.funcs[i]).Name()
		inLens[i] = this.selectCases[i].Chan.Elem().Len()
	}

	return &_Node2StateInfo{StateInfo{this}, inNames, inLens}
}
