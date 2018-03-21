//todo...
package RNCore

import (
	"fmt"
	//"reflect"
	"github.com/robfig/cron"
	"time"
)

var InState chan *StateInfo = nil

type State struct {
	Node

	stateTicker time.Duration //= 60
	saveMaxSpec string
	//每隔5秒执行一次："*/5 * * * * ?"
	//每隔1分钟执行一次："0 */1 * * * ?"
	//每天23点执行一次："0 0 23 * * ?"
	//每天凌晨1点执行一次："0 0 1 * * ?"
	//每月1号凌晨1点执行一次："0 0 1 1 * ?"
	//在26分、29分、33分执行一次："0 26,29,33 * * * ?"
	//每天的0点、13点、18点、21点都执行一次："0 0 0,13,18,21 * * ?"

	In chan *StateInfo

	stateInfos   []*StateInfo
	stateInfoMap map[string]*StateInfo

	maxStateInfos []MaxStateInfo
}

type MaxStateInfo struct {
	Value uint
	Time  time.Time
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

	Values    map[string]uint
	StrValues map[string]string
}

func NewStateInfo(node INode, inCount uint) *StateInfo {
	return &StateInfo{Root().Name(), node.Name(), node.GetOutNodeInfos(), inCount, nil, nil}
}

func NewState(name string, stateTicker time.Duration, saveMaxSpec string) *State {
	return &State{NewNode(name), stateTicker, saveMaxSpec, make(chan *StateInfo, InChanCount), make([]*StateInfo, 0), make(map[string]*StateInfo), nil}
}

/*func (this *State) SetOut(outs []*chan<- interface{}) {
	this.outs = outs
}*/

func (this *State) Run() {

	//
	saveMax := make(chan bool)
	c := cron.New()
	if len(this.saveMaxSpec) <= 0 {
		this.saveMaxSpec = "0 0 6 * * ?" //每天6点执行一次
	}
	c.AddFunc(this.saveMaxSpec, func() { saveMax <- true })
	c.Start()

	//
	save := make(chan bool)

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case <-time.After(time.Second * this.stateTicker):
			this.stateInfos = make([]*StateInfo, 0)

			go func() {
				Root().State()
				save <- true
			}()
			continue
		case <-save:
			this.save()
			continue

		case <-saveMax:
			this.saveMax()
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
	count := this.csvRowCount()

	//
	//----------------------------------------------------------------------
	//csv内容
	_ = this.csvRow(count)
	//todo...
	//row 往csv文件尾部添加

	//
	//----------------------------------------------------------------------
	//csv文件名
	_ = csvFileName()
	//todo...
	//csv文件名 一天一个csv文件

	//
	//----------------------------------------------------------------------
	//新csv文件第一排数据
	_ = this.getFirstRow(count)

	//
	//max
	//----------------------------------------------------------------------
	this.max(count)
}

func (this *State) saveMax() {
	//todo...
	//每天保持一次最大值
	v_row := make([]string, len(this.maxStateInfos))
	t_row := make([]string, len(this.maxStateInfos))
	for i, v := range this.maxStateInfos {
		//一行时间
		//一行最大值
		v_row[i] = fmt.Sprintf("%v", v.Value)
		t_row[i] = v.Time.String()
	}
	_, _ = v_row, t_row
	//todo...
	//往csv文件尾部添加

	//
	//----------------------------------------------------------------------
	//csv文件名
	_ = csvMaxFileName()
	//todo...
	//一年一个csv文件

	//
	//----------------------------------------------------------------------
	//新csv文件第一排数据
	_ = this.getFirstRow(len(this.maxStateInfos))

	//
	//
	//clear
	this.maxStateInfos = nil
}

func basePath() string {
	return "states" //todo... 加上exe所在的目录
}

func csvFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v/%v-%v.%v.%v.state.csv", basePath(), Root().Name(), Time.Year(), Time.Month(), Time.Day())
}
func csvMaxFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v/%v.%v.max_state.csv", basePath(), Root().Name(), Time.Year())
}
func (this *State) csvRowCount() int {
	count := 0
	for _, si := range this.stateInfos {
		count++
		if si.Values != nil {
			count += len(si.Values)
		}
		if si.StrValues != nil {
			count += len(si.StrValues)
		}
		count++ //space
	}
	return count
}
func (this *State) csvRow(count int) []string {
	row := make([]string, count)
	index := 0
	for _, si := range this.stateInfos {
		row[index] = fmt.Sprintf("%v", si.InCount)
		index++

		if si.Values != nil {
			for _, v := range si.Values {
				row[index] = fmt.Sprintf("%v", v)
				index++
			}
		}
		if si.StrValues != nil {
			for _, v := range si.StrValues {
				row[index] = v
				index++
			}
		}
		count++ //space
	}
	return row
}

func (this *State) getFirstRow(count int) []string {

	firstRow := make([]string, count)
	index := 0
	for _, si := range this.stateInfos {
		firstRow[index] = fmt.Sprintf("%v.%v.%v", si.RootName, si.NodeName, "inCount")
		index++

		if si.Values != nil {
			for k, _ := range si.Values {
				firstRow[index] = k
				index++
			}
		}
		if si.StrValues != nil {
			for k, _ := range si.StrValues {
				firstRow[index] = k
				index++
			}
		}
		count++ //space
	}

	return firstRow
}

func (this *State) max(count int) {
	if this.maxStateInfos == nil {
		this.maxStateInfos = make([]MaxStateInfo, count)
	}

	index := 0
	for _, si := range this.stateInfos {
		if si.InCount > this.maxStateInfos[index].Value {
			this.maxStateInfos[index].Value = si.InCount
			this.maxStateInfos[index].Time = time.Now()
		}
		index++

		if si.Values != nil {
			for _, v := range si.Values {
				if v > this.maxStateInfos[index].Value {
					this.maxStateInfos[index].Value = v
					this.maxStateInfos[index].Time = time.Now()
				}
				index++
			}
		}

		if si.StrValues != nil {
			index += len(si.StrValues)
		}
		count++ //space
	}
}

//
func (this *State) OnStateInfo(counts ...*uint) *StateInfo {
	return NewStateInfo(this, *counts[0])
}
