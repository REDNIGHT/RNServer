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

type iLog interface {
	Log(*logData)
}

var _log iLog

//
func Print(node IName, format string, a ...interface{}) {
	_log.Log(doPrintf(node, printLogLevel, format, a))
}
func Warn(node IName, format string, a ...interface{}) {
	_log.Log(doPrintf(node, printWarnLevel, format, a))
}
func Error(node IName, format string, a ...interface{}) {
	_log.Log(doPrintf(node, printErrorLevel, format, a))
}
func Debug(node IName, format string, a ...interface{}) {
	_log.Log(doPrintf(node, printDebugLevel, format, a))
}

func doPrintf(node IName, printLevel string, format string, a ...interface{}) *logData {
	return &logData{time.Now(), Root().Name(), node.Type_Name(), printLevel, fmt.Sprintf(format, a...)}
}

//--------------------------------------------------------------------------------------------------------
type Log struct {
	Node
}

type logData struct {
	Time     time.Time
	RootName string
	NodeName string
	Level    string
	Log      string
}

func NewLog(name string) *Log {
	l := &Log{NewNode(name)}
	if _log != nil {
		l.Panic("_log != nil")
	}
	_log = l
	return l
}

//
func baseLogPath() string {
	return AutoNewPath(ExecPath() + "\\log")
}

func (this *Log) Log(logData *logData) {
	this.SendCall() <- func(IMessage) {
		this.log(logData)
	}
}

func (this *Log) LogByProxy(buffer []byte) {
	this.SendCall() <- func(IMessage) {
		logData := &logData{}
		json.Unmarshal(buffer, logData)
		this.log(logData)
	}
}

func (this *Log) log(logData *logData) {
	fmt.Println("%v>%v>%v", logData.NodeName, logData.Level, logData.Log)

	csvFileName := fmt.Sprintf("%v\\%v.%v.%v.log.csv", baseLogPath(), logData.Time.Year(), logData.Time.Month(), logData.Time.Day())

	buffer := fmt.Sprintf("%v	%v	%v	%v	%v\n", logData.Time, logData.RootName, logData.NodeName, logData.Level, logData.Log)
	ioutil.WriteFile(csvFileName, []byte(buffer), os.ModeAppend)
}

//----------------------------------------------------------------------------------------------------
type LogProxy struct {
	Node

	Out func([]byte)
}

func NewLogProxy(name string) *LogProxy {
	l := &LogProxy{NewNode(name), nil}
	if _log != nil {
		l.Panic("_log != nil")
	}
	_log = l
	return l
}

func (this *LogProxy) Log(logData *logData) {
	this.SendCall() <- func(_ IMessage) {
		this.log(logData)
	}
}

func (this *LogProxy) log(logData *logData) {
	fmt.Println("%v>%v>%v", logData.NodeName, logData.Level, logData.Log)

	buffer, err := json.Marshal(logData)
	if err == nil {
		this.Out(buffer)
	} else {
		this.Error("json.Marshal(logData)  err=%v", err)
	}
}
