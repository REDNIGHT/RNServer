package main

import (
	"fmt"
	"os"
	"os/signal"
)

const version = "0.0.1"

func Run() {

	fmt.Println("RNServer %v starting up", version)

	//Init...

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("RNServer closing down (signal: %v)", sig)

	//Destroy...
}
