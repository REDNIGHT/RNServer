//todo...
package RNCore

import (
	//"fmt"
	//"reflect"
	"time"
)

var InState chan *StateInfo = nil

type State struct {
	Node

	stateTicker time.Duration //= 60

	In chan *StateInfo

	stateInfos   []interface{}
	stateInfoMap map[string]interface{}
}

type IState interface {
	Name() string
	OnStateInfo(counts ...*uint) *StateInfo
	GetOutNodeInfos() []string
}

type StateInfo struct {
	RootName     string
	NodeName     string
	OutNodeInfos []string

	InCount uint

	Map map[string]interface{}
}

func NewStateInfo(node INode, inCount uint) *StateInfo {
	return &StateInfo{Root().Name(), node.Name(), node.GetOutNodeInfos(), inCount, nil}
}

/*func (this *StateInfo) GetInNames() []string {
	return []string{}
}
func (this *StateInfo) GetInCounts() []uint {
	return []uint{}
}*/

func NewState(name string, stateTicker time.Duration) *State {
	return &State{NewNode(name), stateTicker, make(chan *StateInfo, InChanCount), make([]interface{}, 0), make(map[string]interface{})}
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
		case stateInfo := <-this.In:

			this.stateInfos = append(this.stateInfos, stateInfo)
			this.stateInfoMap[stateInfo.NodeName] = stateInfo

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

func (this *State) _save(stateInfo *StateInfo) {
	/*
		//----------------------------------------------------------------------
		//csv内容
		row := make([]string, 0)
		infos := reflect.ValueOf(stateInfo).Elem()
		for i := 0; i < infos.NumField(); i++ {
			row = append(row, fmt.Sprintf("v%", infos.Field(i).Interface()))
		}

		inCounts := stateInfo.GetInCounts()
		for i := 0; i < len(inCounts); i++ {
			row = append(row, fmt.Sprintf("v%", inCounts[i]))
		}

		_ = row
		//todo...
		//row 往csv文件尾部添加

		//----------------------------------------------------------------------
		//csv文件名
		Time := time.Now()
		csvFileName := fmt.Sprintf("%v.%v.%v-%v.%v.%v.csv", Root().Name(), reflect.TypeOf(stateInfo).String(), stateInfo.NodeName(), Time.Year(), Time.Month(), Time.Day())

		_ = csvFileName
		//todo...
		//csv文件名 一天一个csv文件

		//----------------------------------------------------------------------
		//新csv文件第一排数据
		firstRow := make([]string, 0)
		infosT := reflect.TypeOf(stateInfo)
		for i := 0; i < infosT.NumField(); i++ {
			firstRow = append(firstRow, infosT.Field(i).Name)
		}

		inNames := stateInfo.GetInNames()
		firstRow = append(firstRow, inNames...)

		_ = firstRow
		//新csv文件第一排数据
		//todo...
	*/
}

//
func (this *State) OnStateInfo(counts ...*uint) *StateInfo {
	return NewStateInfo(this, *counts[0])
}
