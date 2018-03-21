package RNCore

import (
	"os"
	//"os/exec"
	//"path/filepath"
)

func ExecPath() string {
	//execPath, _ := exec.LookPath(os.Args[0])
	//return filepath.Dir(os.Args[0])//exe所在路径

	path, _ := os.Getwd() //代码所在路径
	return path
}
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func NewPath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func AutoNewPath(path string) string {
	if b, _ := Exists(path); b == false {
		error := NewPath(path)
		if error != nil {
			panic(error)
		}
	}

	return path
}
