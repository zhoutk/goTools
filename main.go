package main

import (
	"os"
	"log"
	"fmt"
	"./backUp"
	"./common"
)

func main()  {
	var flag common.OpFlag
	if len(os.Args) > 1 {
		flag = common.OpFlag{
			Tables: false,
			Datum: false,
			Views: false,
			Funcs: false,
		}
		switch os.Args[1] {
		case "table":
			flag.Tables = true
		case "data":
			flag.Tables = true
			flag.Datum = true
		case "views":
			flag.Views = true
		case "funcs":
			flag.Funcs = true
		default:
			log.Fatal("You arg must be in : table, data, views or funcs.")
		}
	}else{
		flag = common.OpFlag{
			Tables: true,
			Datum: true,
			Views: true,
			Funcs: true,
		}
	}
	err := backUp.Export(flag)
	if err != nil{
		fmt.Println(err)
	}
}