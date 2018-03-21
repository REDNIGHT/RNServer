package RNCore

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	//"strings"
	"time"
)

const (
	printLogLevel   = "log"
	printWarnLevel  = "warn"
	printErrorLevel = "error"
	printDebugLevel = "debug"
)

func Print(node interface{}, format string, a ...interface{}) {
	InLog <- doPrintf(node, printLogLevel, format, a)
}
func Warn(node interface{}, format string, a ...interface{}) {
	InLog <- doPrintf(node, printWarnLevel, format, a)
}
func Error(node interface{}, format string, a ...interface{}) {
	InLog <- doPrintf(node, printErrorLevel, format, a)
}
func Debug(node interface{}, format string, a ...interface{}) {
	InLog <- doPrintf(node, printDebugLevel, format, a)
}

func doPrintf(node interface{}, printLevel string, format string, a ...interface{}) *LogData {

	return &LogData{time.Now(), reflect.TypeOf(node).String() + "." + getNodeName(node), printLevel, fmt.Sprintf(format, a...)}
}

func getNodeName(node interface{}) string {
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
}

//--------------------------------------------------------------------------------------------------------
type Log struct {
	Node
	In chan *LogData
}

var InLog chan *LogData = nil

type LogData struct {
	Time time.Time
	//ServerName string
	NodeName string
	Level    string
	Log      string
}

func NewLog(name string, _default bool) *Log {
	return &Log{NewNode(name), make(chan *LogData, InChanCount)}
}

/*func (this *Log) SetOut(outs []*chan<- interface{}) {
	this.outs = outs
}*/

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

func baseLogPath() string {
	return AutoNewPath(ExecPath() + "\\log")
}

func (this *Log) save(logData *LogData) {

	csvFileName := fmt.Sprintf("%v\\%v.%v.%v.log.csv", baseLogPath(), logData.Time.Year(), logData.Time.Month(), logData.Time.Day())

	buffer := fmt.Sprintf("%v	%v	%v	%v	%v\n", logData.Time, Root().Name(), logData.NodeName, logData.Level, logData.Log)
	ioutil.WriteFile(csvFileName, []byte(buffer), os.ModeAppend)
}

//
func (this *Log) OnStateInfo(counts ...*uint) *StateInfo {
	return NewStateInfo(this, *counts[0])
}
