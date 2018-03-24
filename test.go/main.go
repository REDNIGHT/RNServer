package main

import "fmt"

//import "path/filepath"

//import "io/ioutil"
import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	s := struct{ n, n1 string }{"123", "345"}
	fmt.Println(s)
}
