// test.go project main.go
package main

import (
	"fmt"
)

func main() {
	a := make([]int, 1, 100)
	var a2 = append(a, 1, 2, 3)

	fmt.Println(a)
	fmt.Println(a2)
}

func set(a [4]int) {
	a[0] = -1
}

