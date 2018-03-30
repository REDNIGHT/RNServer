package main

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func main() {
	fmt.Println("-----------------")
	f()
	fmt.Println("-----------------")
}
func f() {
	Debug("test!!!")
}
func Debug(format string, a ...interface{}) {
	doPrintf(format, a)
}

func doPrintf(format string, a ...interface{}) {
	s := string(debug.Stack())
	s = removeTop3(s)
	s = remove_nl(s)
	fmt.Println(s)
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
