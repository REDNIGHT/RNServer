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
	"runtime"
)

func rc() {
	_, f, l, _ := runtime.Caller(0)
	fmt.Printf("%v  %v", f, l)
}
func main() {
	rc()
}
