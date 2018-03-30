package main

import (
	"../RNCore"
)

func main() {
	defer RNCore.CatchPanic()

	root := RNCore.NewRoot("test...")

	//root.Add(node)

	root.Run()

	//
	root.Close()
}
