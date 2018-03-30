/*package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(10)
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(0))
}*/
package main

import (
	"fmt"
	"time"
)

func catchPanic(v string) {
	if err := recover(); err != nil {
		//debug.PrintStack()
		fmt.Printf("%v  %v\n", v, err)
		return
	}
	fmt.Println("catchPanic=" + v)
}
func main() {
	defer catchPanic("111")

	fmt.Println("begin")

	go func() {
		//defer catchPanic("222")
		panic("aaaaaaaaa fuck this world!")
	}()
	go func() {
		time.Sleep(time.Second)

		//defer catchPanic("333")
		panic("bbbbbbbbb fuck this world!")
	}()

	time.Sleep(time.Second * 2)
	fmt.Println("end")

	panic("ccccccccc fuck this world!")

	//fmt.Println("end")
}
