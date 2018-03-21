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

	stateTicker time.Duration //= 60

	In chan IStateInfo

	stateInfos   []interface{}
	stateInfoMap map[string]interface{}
}

type IState interface {
	Name() string
	OnStateInfo(counts ...*uint) IStateInfo
	GetOutNodeInfos() []string
}
type IStateInfo interface {
	NodeName() string

	OutNodeInfos() []string

	GetInNames() []string
	GetInCounts() []uint
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
func (this *StateInfo) GetInCounts() []uint {
	return []uint{}
}

func NewState(name string, stateTicker time.Duration) *State {
	return &State{NewNode(name), stateTicker, make(chan IStateInfo, InChanCount), make([]interface{}, 0), make(map[string]interface{})}
}

/*func (this *State) SetOut(outs []*chan<- interface{}) {
	this.outs = outs
}*/

func (this *State) Run() {

	save := make(chan bool)

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case <-time.After(time.Second * this.stateTicker):
			this.stateInfos = make([]interface{}, 0)

			go func() {
				Root().State()
				save <- true
			}()
			continue
		case <-save:
			this.save()
			continue

			//
		case iStateInfo := <-this.In:

			this.stateInfos = append(this.stateInfos, iStateInfo)
			this.stateInfoMap[iStateInfo.NodeName()] = iStateInfo

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

//
func (this *State) save() {
}

func (this *State) _save(iStateInfo IStateInfo) {

	//----------------------------------------------------------------------
	//csv内容
	row := make([]string, 0)
	infos := reflect.ValueOf(iStateInfo).Elem()
	for i := 0; i < infos.NumField(); i++ {
		row = append(row, fmt.Sprintf("v%", infos.Field(i).Interface()))
	}

	inCounts := iStateInfo.GetInCounts()
	for i := 0; i < len(inCounts); i++ {
		row = append(row, fmt.Sprintf("v%", inCounts[i]))
	}

	_ = row
	//todo...
	//row 往csv文件尾部添加

	//----------------------------------------------------------------------
	//csv文件名
	Time := time.Now()
	csvFileName := fmt.Sprintf("%v.%v.%v-%v.%v.%v.csv", Root().Name(), reflect.TypeOf(iStateInfo).String(), iStateInfo.NodeName(), Time.Year(), Time.Month(), Time.Day())

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
	InCount uint
}

func (this *State) OnStateInfo(counts ...*uint) IStateInfo {
	return &_StateStateInfo{StateInfo{this}, *counts[0]}
}
