package RNCore

import (
	"reflect"
	//"time"
)

type Node2 struct {
	Node

	selectCases []reflect.SelectCase
	funcs       []func(i interface{})

	inCount uint
	//inCounts []uint
}

func NewNode2(name string) Node2 {
	return Node2{NewNode(name), make([]reflect.SelectCase, 0), make([]func(i interface{}), 0), 0}
}

func (this *Node2) AddIn(in interface{}, onFunc func(i interface{})) {
	this.selectCases = append(this.selectCases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(in)})
	this.funcs = append(this.funcs, onFunc)
}

func (this *Node2) Run() {

	this.AddIn(this.StateSig, this._OnState)
	this.AddIn(this.CloseSig, nil)

	//this.inCounts = make([]uint, len(this.selectCases)-2+1) //-2是去掉CloseSig,State +1是InCount
	this.inCount = 0
	for {
		//this.inCounts[len(this.inCounts)-1]++
		this.inCount++

		//
		chosen, recv, recvOk := reflect.Select(this.selectCases)

		if recvOk {
			if chosen == len(this.funcs)-1 {
				this.CloseSig <- true
				return
			} else {
				//this.inCounts[chosen]++

				this.funcs[chosen](recv.Interface())
			}
		}
	}
}

func (this *Node2) _OnState(v interface{}) {
	this.OnState()

	//this.inCounts = make([]uint, len(this.selectCases)-2+1) //-2是去掉CloseSig,State +1是InCount
	this.inCount = 0
	this.StateSig <- true
}

//
type _Node2StateInfo struct {
	StateInfo
	InCount uint
	//InNames  []string
	//InCounts []uint
}

/*func (this *_Node2StateInfo) GetInNames() []string {
	return this.InNames
}
func (this *_Node2StateInfo) GetInLens() []uint {
	return this.InCounts
}*/

func (this *Node2) OnStateInfo(counts ...*uint) IStateInfo {
	/*l := len(this.selectCases) - 2 + 1 //-2是去掉CloseSig,State +1是InCount
	inNames := make([]string, l)
	for i := 0; i < l; i++ {
		inNames[i] = reflect.TypeOf(this.funcs[i]).Name()
	}
	inNames[len(inNames)-1] = "InCount"*/

	return &_Node2StateInfo{StateInfo{this}, this.inCount}
}
