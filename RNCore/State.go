package RNCore

import (
	"fmt"
	//"reflect"
	"encoding/json"
	"github.com/robfig/cron"
	"io/ioutil"
	"os"
	"time"
)

var inStateInfo chan *StateInfo = nil
var inNodeInfo chan *NodeInfo = nil
var inChanOverload chan *ChanOverload = nil

func InStateInfo() chan *StateInfo {
	return inStateInfo
}
func InNodeInfo() chan *NodeInfo {
	return inNodeInfo
}
func InChanOverload() chan *ChanOverload {
	return inChanOverload
}

type State struct {
	Node

	//stateTicker time.Duration //= 60
	stateTickerSpec string
	saveMaxSpec     string
	//每隔5秒执行一次："*/5 * * * * ?"
	//每隔1分钟执行一次："0 */1 * * * ?"
	//每天23点执行一次："0 0 23 * * ?"
	//每天凌晨1点执行一次："0 0 1 * * ?"
	//每月1号凌晨1点执行一次："0 0 1 1 * ?"
	//在26分、29分、33分执行一次："0 26,29,33 * * * ?"
	//每天的0点、13点、18点、21点都执行一次："0 0 0,13,18,21 * * ?"

	//
	In      chan *StateInfo
	InProxy chan []byte

	stateInfos      []*StateInfo
	stateInfosProxy []*StateInfo
	stateInfoMap    map[string]*StateInfo

	//
	maxStateInfos []MaxStateInfo

	//
	InNodeInfo  chan *NodeInfo
	nodeInfoMap map[string]*NodeInfo

	//
	InChanOverload chan *ChanOverload
}

type MaxStateInfo struct {
	Value uint
	Time  time.Time
}

type ChanOverload struct {
	Name    string
	ChanLen int
}

type IState interface {
	Name() string
	OnStateInfo(counts ...*uint) *StateInfo
}

type StateInfo struct {
	RootName string
	NodeName string

	InCount uint

	Values    map[string]uint
	StrValues map[string]string
	//DebugValues map[string]string
}

func (this *StateInfo) key() string {
	return this.RootName + "." + this.NodeName
}

type NodeInfo struct {
	Name     string
	OutNames []string
}

func NewStateInfo(node INode, inCount uint) *StateInfo {
	return &StateInfo{Root().Name(), node.Name(), inCount, nil, nil}
}

func NewState(name string, stateTickerSpec string, saveMaxSpec string) *State {
	state := &State{NewNode(name),
		stateTickerSpec,
		saveMaxSpec,
		make(chan *StateInfo, InChanCount),
		make(chan []byte, InChanCount),

		make([]*StateInfo, 0),
		make([]*StateInfo, 0),
		make(map[string]*StateInfo),
		nil,

		make(chan *NodeInfo, InChanCount),
		make(map[string]*NodeInfo),

		make(chan *ChanOverload, InChanCount)}

	if inStateInfo != nil {
		panic("inStateInfo != nil")
	}
	inStateInfo = state.In

	if inNodeInfo != nil {
		panic("inNodeInfo != nil")
	}
	inNodeInfo = state.InNodeInfo

	if inChanOverload != nil {
		panic("inChanOverload != nil")
	}
	inChanOverload = state.InChanOverload

	return state
}

func (this *State) Run() {

	c := cron.New()
	c.Start()

	//
	stateTicker := make(chan bool)
	if len(this.stateTickerSpec) <= 0 {
		this.stateTickerSpec = "0 */1 * * * ?" //每分钟执行一次
	}
	c.AddFunc(this.stateTickerSpec, func() { stateTicker <- true })

	//
	saveMax := make(chan bool)
	if len(this.saveMaxSpec) <= 0 {
		this.saveMaxSpec = "0 0 6 * * ?" //每天6点执行一次
	}
	c.AddFunc(this.saveMaxSpec, func() { saveMax <- true })

	//
	debugChanStateTicker := make(chan bool)
	if RNCDebug {
		c.AddFunc(RNCDebugStateTickerSpec, func() { debugChanStateTicker <- true }) //每10s执行一次
	}

	//
	save := make(chan bool)

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case <-stateTicker:
			this.stateInfos = make([]*StateInfo, 0)

			//
			this.stateInfos = append(this.stateInfos, this.stateInfosProxy...)
			for _, si := range this.stateInfosProxy {
				this.stateInfoMap[si.key()] = si
			}
			this.stateInfosProxy = make([]*StateInfo, 0)

			//
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

		case <-debugChanStateTicker:
			Root().DebugChanState()
			continue

		//
		case stateInfo := <-this.In:

			this.stateInfos = append(this.stateInfos, stateInfo)
			this.stateInfoMap[stateInfo.key()] = stateInfo

			continue

		case nodeInfo := <-this.InNodeInfo:
			this.nodeInfoMap[nodeInfo.Name] = nodeInfo
			continue

		case chanOverload := <-this.InChanOverload:
			this.saveChanOverload(chanOverload)
			continue

			//
		case buffer := <-this.InProxy:
			this.OnInProxy(buffer)
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

func (this *State) OnInProxy(buffer []byte) {
	j := &proxyDataJ{}
	json.Unmarshal(buffer, j)

	if j.StateInfo != nil && j.NodeInfo != nil {
		this.Error("j.StateInfo != nil && j.NodeInfo != nil")
		return
	}

	if j.StateInfo != nil {
		this.stateInfosProxy = append(this.stateInfosProxy, j.StateInfo)
		return
	}
	if j.NodeInfo != nil {
		this.nodeInfoMap[j.NodeInfo.Name] = j.NodeInfo
	}
	if j.ChanOverload != nil {
		this.saveChanOverload(j.ChanOverload)
	}
}

//
func chanOverloadFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v\\%v-%v.%v.%v.chanOverload.csv", baseStatesPath(), Root().Name(), Time.Year(), Time.Month(), Time.Day())
}
func (this *State) saveChanOverload(chanOverload *ChanOverload) {

	row := fmt.Sprintf("%v	%v	%v\n", time.Now(), chanOverload.Name, chanOverload.ChanLen)

	ioutil.WriteFile(chanOverloadFileName(), []byte(row), os.ModeAppend)
}

//
func (this *State) save() {
	count := this.csvRowCount()

	//
	//----------------------------------------------------------------------
	//csv文件名
	csvFile := stateFileName()
	if b, _ := Exists(csvFile); b == false {
		//新csv文件第一排数据
		firstRow := this.getFirstRow(count)
		buffer := getRowBuffer(firstRow)
		ioutil.WriteFile(csvFile, []byte(buffer), os.ModeAppend)
	}

	//
	//----------------------------------------------------------------------
	//csv内容
	row := this.csvRow(count)
	buffer := getRowBuffer(row)
	ioutil.WriteFile(csvFile, []byte(buffer), os.ModeAppend)

	//
	//max
	//----------------------------------------------------------------------
	this.max(count)
}

func (this *State) saveMax() {

	//csv文件名
	csvMaxFile := maxFileName()
	if b, _ := Exists(csvMaxFile); b == false {
		//新csv文件第一排数据
		firstRow := this.getFirstRow(len(this.maxStateInfos))
		buffer := getRowBuffer(firstRow)
		ioutil.WriteFile(csvMaxFile, []byte(buffer), os.ModeAppend)
	}

	//
	//每天保持一次最大值
	t_row := make([]string, len(this.maxStateInfos))
	v_row := make([]string, len(this.maxStateInfos))
	for i, v := range this.maxStateInfos {
		//一行时间
		//一行最大值
		t_row[i] = v.Time.String()
		v_row[i] = fmt.Sprintf("%v", v.Value)
	}
	buffer := getRowBuffer(t_row)
	ioutil.WriteFile(csvMaxFile, []byte(buffer), os.ModeAppend)
	buffer = getRowBuffer(v_row)
	ioutil.WriteFile(csvMaxFile, []byte(buffer), os.ModeAppend)

	//
	//
	//clear
	this.maxStateInfos = nil
}

func baseStatesPath() string {
	return AutoNewPath(ExecPath() + "\\states")
}

func stateFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v\\%v-%v.%v.%v.state.csv", baseStatesPath(), Root().Name(), Time.Year(), Time.Month(), Time.Day())
}
func maxFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v\\%v.%v.max_state.csv", baseStatesPath(), Root().Name(), Time.Year())
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
func getRowBuffer(row []string) string {
	buffer := ""
	for _, r := range row {
		buffer += r + "\t"
	}
	buffer += "\n"
	return buffer
}

//
func (this *State) OnStateInfo(counts ...*uint) *StateInfo {
	return NewStateInfo(this, *counts[0])
}

func (this *State) DebugChanState() {
	this.OnDebugChanState("In", len(this.In))
	this.OnDebugChanState("InProxy", len(this.InProxy))
	this.OnDebugChanState("InNodeInfo", len(this.InNodeInfo))
	this.OnDebugChanState("InChanOverload", len(this.InChanOverload))
}

//----------------------------------------------------------------------------------------------------------------

type StateProxy struct {
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

	InNodeInfo chan *NodeInfo

	InChanOverload chan *ChanOverload

	Out func([]byte)
}

type proxyDataJ struct {
	StateInfo    *StateInfo
	NodeInfo     *NodeInfo
	ChanOverload *ChanOverload
}

func NewStateProxy(name string, stateTicker time.Duration, saveMaxSpec string) *StateProxy {
	state := &StateProxy{NewNode(name), stateTicker, saveMaxSpec, make(chan *StateInfo, InChanCount), make(chan *NodeInfo, InChanCount), make(chan *ChanOverload, InChanCount), nil}

	if inStateInfo != nil {
		panic("inStateInfo != nil")
	}
	inStateInfo = state.In

	if inNodeInfo != nil {
		panic("inNodeInfo != nil")
	}
	inNodeInfo = state.InNodeInfo

	if inChanOverload != nil {
		panic("inChanOverload != nil")
	}
	inChanOverload = state.InChanOverload
	return state
}

func (this *StateProxy) Run() {

	//
	c := cron.New()
	c.Start()

	debugChanStateTicker := make(chan bool)
	if RNCDebug {
		c.AddFunc(RNCDebugStateTickerSpec, func() { debugChanStateTicker <- true })
	}

	//
	var inCount uint = 0
	for {
		inCount++

		//
		select {

		case <-time.After(time.Second * this.stateTicker):
			go func() { Root().State() }()
			continue

			//
		case stateInfo := <-this.In:

			buffer, err := json.Marshal(&proxyDataJ{stateInfo, nil, nil})
			if err == nil {
				this.Out(buffer)
			} else {
				this.Error("json.Marshal(stateInfo)  err=%v", err)
			}
			continue

			//
		case nodeInfo := <-this.InNodeInfo:

			buffer, err := json.Marshal(&proxyDataJ{nil, nodeInfo, nil})
			if err == nil {
				this.Out(buffer)
			} else {
				this.Error("json.Marshal(nodeInfo)  err=%v", err)
			}
			continue

			//
		case chanOverload := <-this.InChanOverload:

			buffer, err := json.Marshal(&proxyDataJ{nil, nil, chanOverload})
			if err == nil {
				this.Out(buffer)
			} else {
				this.Error("json.Marshal(nodeInfo)  err=%v", err)
			}
			continue

		case <-debugChanStateTicker:
			Root().DebugChanState()
			continue

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

func (this *StateProxy) DebugChanState() {
	this.OnDebugChanState("In", len(this.In))
	this.OnDebugChanState("InNodeInfo", len(this.InNodeInfo))
	this.OnDebugChanState("InChanOverload", len(this.InChanOverload))
}
