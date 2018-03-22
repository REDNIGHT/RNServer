package RNCore

import ()

func mian() {
	root := NewRoot("test...")

	//root.Add(node)

	root.Init()
	root.Register()

	root.Run()

	//
	root.Close()
	root.Destroy()
}
