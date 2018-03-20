//todo...
package RNCore

import (
	"fmt"
	"reflect"
	"time"
)

var InState chan IStateInfo = nil

type State struct {
	Node

	In chan IStateInfo

	stateInfos map[string]interface{}
}

type IState interface {
	Name() string
	OnState() IStateInfo
	GetOutNodeInfos() []string
}
type IStateInfo interface {
	NodeName() string

	OutNodeInfos() []string

	GetInNames() []string
	GetInLens() []int
}

type StateInfo struct {
	IState IState
}

func (this *StateInfo) NodeName() string {
	return this.IState.Name()
}
func (this *StateInfo) OutNodeInfos() []string {
	return this.IState.GetOutNodeInfos()
}

func (this *StateInfo) GetInNames() []string {
	return []string{}
}
func (this *StateInfo) GetInLens() []int {
	return []int{}
}

func NewState(name string, _default bool) *State {
	return &State{NewNode(name), make(chan IStateInfo, InChanCount), make(map[string]interface{})}
}

/*func (this *State) SetOut(outs []*chan<- interface{}) {
	this.outs = outs
}*/

func (this *State) Run() {

	//
	for {
		//
		select {
		case iStateInfo := <-this.In:

			this.stateInfos[iStateInfo.NodeName()] = iStateInfo

			this.save(iStateInfo)

			continue

			//
		case <-time.After(time.Second * StateTime):
			this.State()
			continue
		case <-this.CloseSig:
			this.CloseSig <- true
			return
		}
	}
}

func (this *State) save(iStateInfo IStateInfo) {

	//----------------------------------------------------------------------
	//csv内容
	row := make([]string, 0)
	infos := reflect.ValueOf(iStateInfo).Elem()
	for i := 0; i < infos.NumField(); i++ {
		row = append(row, fmt.Sprintf("v%", infos.Field(i).Interface()))
	}

	inLens := iStateInfo.GetInLens()
	for i := 0; i < len(inLens); i++ {
		row = append(row, fmt.Sprintf("v%", inLens[i]))
	}

	_ = row
	//todo...
	//row 往csv文件尾部添加

	//----------------------------------------------------------------------
	//csv文件名
	Time := time.Now()
	csvFileName := fmt.Sprintf("%v.%v.%v-%v.%v.%v.csv", ServerName, reflect.TypeOf(iStateInfo).String(), iStateInfo.NodeName(), Time.Year(), Time.Month(), Time.Day())

	_ = csvFileName
	//todo...
	//csv文件名 一天一个csv文件

	//----------------------------------------------------------------------
	//新csv文件第一排数据
	firstRow := make([]string, 0)
	infosT := reflect.TypeOf(iStateInfo)
	for i := 0; i < infosT.NumField(); i++ {
		firstRow = append(firstRow, infosT.Field(i).Name)
	}

	inNames := iStateInfo.GetInNames()
	firstRow = append(firstRow, inNames...)

	_ = firstRow
	//新csv文件第一排数据
	//todo...
}

//
type _StateStateInfo struct {
	StateInfo
	In int
}

func (this *State) OnState() IStateInfo {
	return &_StateStateInfo{StateInfo{this}, len(this.In)}
}
