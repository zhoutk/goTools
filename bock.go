package main

import (
	Bock "./sqlBock/Mysql"
	"fmt"
)

func main()  {
	table := Bock.Bock{
		Table: "books",
	}

	//insert one record
	params := make(map[string] interface{})
	args := make(map[string] interface{})
	session := make(map[string] interface{})
	session["userid"] = "112"
	args["session"] = session
	params["name"] = "golang实战3443"
	params["status"] = 0
	db := &table
	rs := db.Create(params, args)
	fmt.Println(rs)

	//update one record
	params = make(map[string] interface{})
	args = make(map[string] interface{})
	args["id"] = 2
	params["name"] = "update 2"
	params["status"] = 9
	rs = db.Update(params, args)
	fmt.Println(rs)
}
