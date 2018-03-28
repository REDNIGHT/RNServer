package main

import (
	"RNCore"
)

func main() {
	root := RNCore.NewRoot("test...")

	//root.Add(node)

	root.Run()

	//
	root.Close()
}
