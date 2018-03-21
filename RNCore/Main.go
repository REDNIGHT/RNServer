package RNCore

import ()

func mian() {
	root := NewRoot("testS")

	root.Init()
	root.Register()

	root.Run()

	//
	root.Close()
	root.Destroy()
}
