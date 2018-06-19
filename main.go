package main

import (
	"./backUp"
	"fmt"
)
func main()  {
	err := backUp.Export()
	if err != nil{
		fmt.Println(err)
	}
}