package main

import (
	Bock "./sqlBock/Mysql"
	"fmt"
)

func main()  {
	table := Bock.Bock{
		Table: "role",
	}
	params := make(map[string] interface{})
	args := make(map[string] interface{})
	//db := make([]DB.IBock, 1)
	//db[0] = &table
	//rs, _ := db[0].Retrieve("", params, nil, nil)
	//db := &table
	//rs, _ := db.Retrieve("afataea", values)
	//fields := []string {"id", "name"}
	//args["fields"] = fields
	params["id"] = "1090912"
	params["name"] = "测试插入"
	params["is_rank"] = 8
	db := &table
	rs := db.Create(params, args)
	fmt.Println(rs)
}
