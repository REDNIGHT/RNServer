//todo...
//保存的文件 最多只有20份 每周清空多出来的旧文件

package RNCore

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	printLogLevel   = "log"
	printWarnLevel  = "warn"
	printErrorLevel = "error"
	printDebugLevel = "debug"
)

func Print(node IName, format string, a ...interface{}) {
	inLog <- doPrintf(node, printLogLevel, format, a)
}
func Warn(node IName, format string, a ...interface{}) {
	inLog <- doPrintf(node, printWarnLevel, format, a)
}
func Error(node IName, format string, a ...interface{}) {
	inLog <- doPrintf(node, printErrorLevel, format, a)
}
func Debug(node IName, format string, a ...interface{}) {
	inLog <- doPrintf(node, printDebugLevel, format, a)
}

func doPrintf(node IName, printLevel string, format string, a ...interface{}) *logData {
	return &logData{time.Now(), Root().Name(), node.Type_Name(), printLevel, fmt.Sprintf(format, a...)}
}

/*func getNodeName(node interface{}) string {
	v_nil := reflect.Value{}

	nv := reflect.ValueOf(node).Elem()
	if nv == v_nil {
		return ""
	}

	fun := nv.MethodByName("Name")
	if fun == v_nil {
		return ""
	}

	//
	return fun.Call(nil)[0].String()
}*/

//--------------------------------------------------------------------------------------------------------
type Log struct {
	Node

	In      chan *logData
	InProxy chan []byte
}

var inLog chan *logData = nil

func InLog() chan *logData {
	return inLog
}

type logData struct {
	Time     time.Time
	RootName string
	NodeName string
	Level    string
	Log      string
}

func NewLog(name string) *Log {
	log := &Log{NewNode(name), make(chan *logData, InChanLen), make(chan []byte, InChanLen)}
	if inLog == nil {
		inLog = log.In
	}
	return log
}

func (this *Log) Run() {

	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case logData := <-this.In:
			fmt.Println("%v>%v>%v", logData.NodeName, logData.Level, logData.Log)
			this.save(logData)
			continue

		case buffer := <-this.InProxy:
			logData := &logData{}
			json.Unmarshal(buffer, logData)
			this.save(logData)
			continue

		case f := <-this.messageChan:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

func baseLogPath() string {
	return AutoNewPath(ExecPath() + "\\log")
}

func (this *Log) save(logData *logData) {
	csvFileName := fmt.Sprintf("%v\\%v.%v.%v.log.csv", baseLogPath(), logData.Time.Year(), logData.Time.Month(), logData.Time.Day())

	buffer := fmt.Sprintf("%v	%v	%v	%v	%v\n", logData.Time, logData.RootName, logData.NodeName, logData.Level, logData.Log)
	ioutil.WriteFile(csvFileName, []byte(buffer), os.ModeAppend)
}

//
func (this *Log) DebugChanState(chanOverload chan *ChanOverload) {
	this.TestChanOverload(chanOverload, "In", len(this.In))
}

//----------------------------------------------------------------------------------------------------
type LogProxy struct {
	Node

	In  chan *logData
	Out func([]byte)
}

func NewLogProxy(name string) *LogProxy {
	logProxy := &LogProxy{NewNode(name), make(chan *logData, InChanLen), nil}
	inLog = logProxy.In
	return logProxy
}

func (this *LogProxy) Run() {

	var inCount uint = 0
	for {
		inCount++

		//
		select {
		case logData := <-this.In:

			fmt.Println("%v>%v>%v", logData.NodeName, logData.Level, logData.Log)

			buffer, err := json.Marshal(logData)
			if err == nil {
				this.Out(buffer)
			} else {
				this.Error("json.Marshal(logData)  err=%v", err)
			}
			continue

		case f := <-this.messageChan:
			if this.OnMessage(f) == true {
				return
			}
		}
	}
}

func (this *LogProxy) DebugChanState(chanOverload chan *ChanOverload) {
	this.TestChanOverload(chanOverload, "In", len(this.In))
}
