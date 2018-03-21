package main

import "os"
import "fmt"

//import "path/filepath"

//import "io/ioutil"
import "RNCore"

func main() {

	//execDir := filepath.Dir(os.Args[0])
	execDir, _ := os.Getwd()

	fmt.Println(execDir)

	path := execDir + "\\ttt123"
	//file := path + "\\test.txt"

	b, _ := RNCore.Exists(path)
	fmt.Println("Exists=", b)
	if b == false {
		os.MkdirAll(path, os.ModePerm)
	}

	b, _ = RNCore.Exists(path)
	fmt.Println("222Exists=", b)
	/*error := ioutil.WriteFile(file, []byte("fuck the world!"), os.ModeAppend)
	fmt.Println(error)

	b, e := ioutil.ReadFile(file)
	fmt.Println(string(b), e)*/
}
