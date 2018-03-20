package RNCore

import (
	"os"
	"os/signal"
)

func mian() {
	root := Root()

	root.Init()
	root.Register()

	root.Run()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	Print(root, "closing down (signal: %v)", sig)

	//
	root.Close()
	root.Destroy()
}
