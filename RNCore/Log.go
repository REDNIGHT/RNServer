//todo...
//保存的文件 最多只有20份 每周清空多出来的旧文件

package RNCore

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

const (
	logLevel   = "log"
	warnLevel  = "warn"
	errorLevel = "error"
	debugLevel = "debug"
	panicLevel = "panic"
)

type iLog interface {
	Log(*_LogData)
}

var _log iLog

//
func Print(node IName, format string, a ...interface{}) {
	_log.Log(newLogData(false, node, logLevel, format, a))
}
func Warn(node IName, format string, a ...interface{}) {
	_log.Log(newLogData(true, node, warnLevel, format, a))
}
func Error(node IName, format string, a ...interface{}) {
	_log.Log(newLogData(true, node, errorLevel, format, a))
}
func Debug(node IName, format string, a ...interface{}) {
	_log.Log(newLogData(true, node, debugLevel, format, a))
}
func Panic(node IName, format string, a ...interface{}) {
	ld := newLogData(true, node, panicLevel, format, a)
	_log.Log(ld)
	panic(ld)
}

func newLogData(stack bool, node IName, printLevel string, format string, a ...interface{}) *_LogData {
	s := ""
	if stack {
		s = string(debug.Stack())
		s = removeTop3(s)
	}
	return &_LogData{time.Now(), Root().Name(), node.Type_Name(), printLevel, fmt.Sprintf(format, a...), s}
}

func removeTop3(s string) string {
	ss := strings.Split(s, "\n")
	ss2 := ss[7:]
	s = ""
	for i, v := range ss2 {
		s += v
		if i < len(ss2)-1 {
			s += "\n"
		}
	}
	s += ss[0]

	s = strings.Replace(s, "\n\t", "\t", -1)
	s = strings.Replace(s, "\n", "\\n", -1)
	return s
}

func CatchPanic() {
	if r := recover(); r != nil {
		_, b := r.(*_LogData)
		if b == false {
			ld := newLogData(false, _log.(IName), panicLevel, "panic:%v", r)
			_log.Log(ld)
		}
	}

	Root().Close()

	/*log0, b0 := _log.(*Log)
	log1, b1 := _log.(*LogProxy)
	if b0 {
		for len(log0.InCall()) > 0 || len(log0.InMessage()) > 0 {
			time.Sleep(time.Millisecond * 10)
		}
	}
	if b1 {
		for len(log1.InCall()) > 0 || len(log1.InMessage()) > 0 {
			time.Sleep(time.Millisecond * 10)
		}
	}*/
}

//--------------------------------------------------------------------------------------------------------
type Log struct {
	MNode
}

type _LogData struct {
	Time     time.Time
	RootName string
	NodeName string
	Level    string
	Log      string
	Stack    string
}

func NewLog(name string) *Log {
	l := &Log{NewMNode(name)}
	if _log != nil {
		l.Panic("_log != nil")
	}
	_log = l
	return l
}

func (this *Log) Close() {
	//this.SendMessage(nil)//禁止退出 运行都最后
	//close(this.inMessage)
}

//
func baseLogPath() string {
	return AutoNewPath(ExecPath() + "\\log")
}

func (this *Log) Log(logData *_LogData) {
	this.InCall() <- func(IMessage) {
		this.log(logData)
	}
}

func (this *Log) LogByProxy(buffer []byte) {
	this.InCall() <- func(IMessage) {
		logData := &_LogData{}
		json.Unmarshal(buffer, logData)
		this.log(logData)
	}
}

func (this *Log) log(logData *_LogData) {
	save(logData)
}
func save(logData *_LogData) {
	fmt.Printf("%v>%v>%v\n", logData.NodeName, logData.Level, logData.Log)

	csvFileName := fmt.Sprintf("%v\\%v.%v.%v.log.csv", baseLogPath(), logData.Time.Year(), logData.Time.Month(), logData.Time.Day())

TO:
	buffer := fmt.Sprintf("%v	%v	%v	%v	%v	%v\n", logData.Time, logData.RootName, logData.NodeName, logData.Level, logData.Log, logData.Stack)
	ioutil.WriteFile(csvFileName, []byte(buffer), os.ModeAppend)

	if logData.Level == panicLevel {
		logData.Level += "!"
		csvFileName = fmt.Sprintf("%v\\%v.%v.%v.panic.csv", baseLogPath(), logData.Time.Year(), logData.Time.Month(), logData.Time.Day())
		goto TO
	}
}

//----------------------------------------------------------------------------------------------------
type LogProxy struct {
	MNode

	Out func([]byte)
}

func NewLogProxy(name string) *LogProxy {
	l := &LogProxy{NewMNode(name), nil}
	if _log != nil {
		l.Panic("_log != nil")
	}
	_log = l
	return l
}

func (this *LogProxy) Close() {
	//this.SendMessage(nil)//禁止退出 运行都最后
	//close(this.inMessage)
}

func (this *LogProxy) Log(logData *_LogData) {
	this.InCall() <- func(_ IMessage) {
		this.log(logData)
	}
}

func (this *LogProxy) log(logData *_LogData) {
	save(logData)

	buffer, err := json.Marshal(logData)
	if err == nil {
		this.Out(buffer)
	} else {
		this.Error("json.Marshal(logData)  err=%v", err)
	}
}
