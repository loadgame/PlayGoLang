package main

import (
	"fmt"
	"io/ioutil"
	//"os"
)

func ListDir(dir string) {

	entries, err := ioutil.ReadDir(dir)

	if err == nil {
		for _, entr := range entries {

			if entr.IsDir() {
				if entr.Name() != ".git" {
					ListDir(dir + "\\" + entr.Name())
				}

			} else {
				fmt.Println(dir+"\\"+entr.Name(), "==", entr.Size())
			}
		}
	}
}
func main() {
	// if len(os.Args) < 2 {
	//	fmt.Println("Need directory as parameter")
	//	return
	// }
	// dir := os.Args[1]

	dir := "D:\\source.go"
	ListDir(dir)

	//entries, err := ioutil.ReadDir(dir)

	//if err == nil {
	//	for _, entr := range entries {
	//		fmt.Println(entr.Name(), entr.Size())

	//	}
	//}
	//fmt.Println("err", err)
}
