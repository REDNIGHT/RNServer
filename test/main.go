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
	"runtime/debug"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	/*value := 111
	zero := 0
	value = value / zero*/
	panic("fuck this world!")
}
