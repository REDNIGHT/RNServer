package RNCore

func main() {
	defer SavePanic()

	root := NewRoot("test...")

	//root.Add(node)

	root.Run()

	//
	root.Close()
}
