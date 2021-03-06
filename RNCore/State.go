//todo...
//保存的文件 每种最多只有20份 每周清空多出来的旧文件

package RNCore

import (
	"fmt"
	//"reflect"
	"io/ioutil"
	"os"
	"time"

	"github.com/robfig/cron"
)

type iState interface {
	AddNodeInfo(*NodeInfo)
}

var _State iState

func AddNodeInfo(nodeInfo *NodeInfo) {
	_State.AddNodeInfo(nodeInfo)
}

type State struct {
	MNode

	//每隔5秒执行一次："*/5 * * * * ?"
	//每隔1分钟执行一次："0 */1 * * * ?"
	//每天23点执行一次："0 0 23 * * ?"
	//每天凌晨1点执行一次："0 0 1 * * ?"
	//每月1号凌晨1点执行一次："0 0 1 1 * ?"
	//在26分、29分、33分执行一次："0 26,29,33 * * * ?"
	//每天的0点、13点、18点、21点都执行一次："0 0 0,13,18,21 * * ?"

	//
	stateInfos      []*StateInfo
	stateInfosProxy []*StateInfo
	stateInfoMap    map[string]*StateInfo

	//
	maxStateInfos []MaxStateInfo

	//
	nodeInfoMap map[string]*NodeInfo
}

type MaxStateInfo struct {
	Value uint
	Time  time.Time
}

type StateWarning struct {
	Name, Warning string
}

type IState interface {
	//Name() string

	GetStateInfo() *StateInfo

	GetStateWarning(func(Name, Warning string))
}

type StateInfo struct {
	RootName string
	NodeName string

	InTotal uint

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

func NewStateInfo(node IName, inTotal uint) *StateInfo {
	return &StateInfo{Root().Name(), node.Name(), inTotal, nil, nil}
}

func NewState(name string, stateTickerSpec string, saveMaxSpec string) *State {
	this := &State{NewMNode(name),

		make([]*StateInfo, 0),
		make([]*StateInfo, 0),
		make(map[string]*StateInfo),
		nil,

		make(map[string]*NodeInfo)}

	if _State != nil {
		this.Panic(nil, "_State != nil")
	}
	_State = this

	//
	//
	c := cron.New()
	c.Start()

	//
	if len(stateTickerSpec) <= 0 {
		stateTickerSpec = "0 */1 * * * ?" //每分钟执行一次
	}
	c.AddFunc(stateTickerSpec, func() { this.InCall() <- func(_ ICall) { this.stateTicker() } })

	//
	if len(saveMaxSpec) <= 0 {
		saveMaxSpec = "0 0 6 * * ?" //每天6点执行一次
	}
	c.AddFunc(saveMaxSpec, func() { this.InCall() <- func(_ ICall) { this.saveMax() } })

	//
	if RNCDebug {
		c.AddFunc(RNCStateWarningTickerSpec, func() { this.stateWarning() }) //每10s执行一次
	}

	return this
}

func (this *State) AddNodeInfo(nodeInfo *NodeInfo) {
	this.InCall() <- func(ICall) {
		if _, b := this.nodeInfoMap[nodeInfo.Name]; b == true {
			this.Error("b := this.nodeInfoMap[nodeInfo.Name]; b == true  nodeInfo.Name=%v", nodeInfo.Name)
		}
		this.nodeInfoMap[nodeInfo.Name] = nodeInfo
	}
}

func (this *State) stateWarning() {
	Root().ForEach(func(node interface{}) {
		if is, b := node.(IState); b == true {
			is.GetStateWarning(this.inStateWarning)
		}
	})
}
func (this *State) inStateWarning(name, warning string) {
	this.InCall() <- func(ICall) {
		this.saveStateWarning(name, warning)
	}
}

func (this *State) stateTicker() {
	this.stateInfos = make([]*StateInfo, 0)

	//
	this.stateInfos = append(this.stateInfos, this.stateInfosProxy...)
	for _, si := range this.stateInfosProxy {
		this.stateInfoMap[si.key()] = si
	}
	this.stateInfosProxy = make([]*StateInfo, 0)

	//
	go func() {
		Root().BroadcastMessage(func(node IMessage) {
			if is, b := node.(IState); b == true {
				this.inStateInfo(is.GetStateInfo())
			}
		})

		this.InCall() <- func(ICall) {
			this.save()
		}
	}()
}
func (this *State) inStateInfo(stateInfo *StateInfo) {
	this.InCall() <- func(ICall) {
		this.stateInfos = append(this.stateInfos, stateInfo)
		this.stateInfoMap[stateInfo.key()] = stateInfo
	}
}

func (this *State) InProxy(buffer []byte) {
	this.InCall() <- func(ICall) {
		this.inProxy(buffer)
	}
}
func (this *State) inProxy(buffer []byte) {
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
	if j.StateWarning != nil {
		this.saveStateWarning(j.StateWarning.Name, j.StateWarning.Warning)
	}
}

//
func stateWarningFileName() string {
	Time := time.Now()
	return fmt.Sprintf("%v\\%v-%v.%v.%v.stateWarning.csv", baseStatesPath(), Root().Name(), Time.Year(), Time.Month(), Time.Day())
}
func (this *State) saveStateWarning(name, warning string) {

	row := fmt.Sprintf("%v	%v	%v\n", time.Now(), name, warning)

	ioutil.WriteFile(stateWarningFileName(), []byte(row), os.ModeAppend)
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
		row[index] = fmt.Sprintf("%v", si.InTotal)
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
		firstRow[index] = fmt.Sprintf("%v.%v.%v", si.RootName, si.NodeName, "inTotal")
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
		if si.InTotal > this.maxStateInfos[index].Value {
			this.maxStateInfos[index].Value = si.InTotal
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

//----------------------------------------------------------------------------------------------------------------

type StateProxy struct {
	MNode

	//每隔5秒执行一次："*/5 * * * * ?"
	//每隔1分钟执行一次："0 */1 * * * ?"
	//每天23点执行一次："0 0 23 * * ?"
	//每天凌晨1点执行一次："0 0 1 * * ?"
	//每月1号凌晨1点执行一次："0 0 1 1 * ?"
	//在26分、29分、33分执行一次："0 26,29,33 * * * ?"
	//每天的0点、13点、18点、21点都执行一次："0 0 0,13,18,21 * * ?"

	//In chan *StateInfo

	//InNodeInfo chan *NodeInfo

	//InStateWarning chan *StateWarning

	Out func([]byte)
}

type proxyDataJ struct {
	StateInfo    *StateInfo
	NodeInfo     *NodeInfo
	StateWarning *StateWarning
}

func NewStateProxy(name, stateTickerSpec, saveMaxSpec string) *StateProxy {
	this := &StateProxy{NewMNode(name), nil}

	if _State != nil {
		this.Panic(nil, "_State != nil")
	}
	_State = this

	//
	c := cron.New()
	c.Start()

	//
	if len(stateTickerSpec) <= 0 {
		stateTickerSpec = "0 */1 * * * ?" //每分钟执行一次
	}
	c.AddFunc(stateTickerSpec, this.stateTicker)

	//
	if RNCDebug {
		c.AddFunc(RNCStateWarningTickerSpec, func() { this.stateWarning() })
	}

	return this
}

func (this *StateProxy) AddNodeInfo(nodeInfo *NodeInfo) {
	buffer, err := json.Marshal(&proxyDataJ{nil, nodeInfo, nil})
	if err == nil {
		this.Out(buffer)
	} else {
		this.Error("json.Marshal(nodeInfo)  err=%v", err)
	}
}

func (this *StateProxy) stateWarning() {
	Root().ForEach(func(node interface{}) {
		if is, b := node.(IState); b == true {
			is.GetStateWarning(this.inStateWarning)
		}
	})
}

func (this *StateProxy) inStateWarning(name, warning string) {
	this.InCall() <- func(ICall) {
		buffer, err := json.Marshal(&proxyDataJ{nil, nil, &StateWarning{name, warning}})
		if err == nil {
			this.Out(buffer)
		} else {
			this.Error("json.Marshal(nodeInfo)  err=%v", err)
		}
	}
}

func (this *StateProxy) stateTicker() {
	this.InCall() <- func(ICall) {
		Root().BroadcastMessage(func(node IMessage) {
			if is, b := node.(IState); b == true {
				this.inStateInfo(is.GetStateInfo())
			}
		})
	}
}
func (this *StateProxy) inStateInfo(stateInfo *StateInfo) {
	this.InCall() <- func(ICall) {
		buffer, err := json.Marshal(&proxyDataJ{stateInfo, nil, nil})
		if err == nil {
			this.Out(buffer)
		} else {
			this.Error("json.Marshal(stateInfo)  err=%v", err)
		}
	}
}
