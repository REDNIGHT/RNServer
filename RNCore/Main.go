package RNCore

func main() {
	defer CatchPanic()

	root := NewRoot("test...")

	//root.Add(node)

	root.Run()

	//
	root.Close()
}
