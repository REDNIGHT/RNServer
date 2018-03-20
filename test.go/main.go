package main

import "fmt"

func main() {

	messages := make(chan string, 100)

	messages <- "buffered"
	messages <- "channel"

	fmt.Println(len(messages))
	fmt.Println(<-messages)
	fmt.Println(len(messages))
	fmt.Println(<-messages)
	fmt.Println(len(messages))
}
