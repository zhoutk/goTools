package main

import (
	"./db"
	"fmt"
)

func main()  {
	table := db.BaseDao{
		Table: "tablename",
	}
	var values [] interface{}
	db := make([]db.Idb, 1)
	db[0] = &table
	rs, _ := db[0].Retrieve("sfaf", values)
	//db := &table
	//rs, _ := db.Retrieve("afataea", values)
	fmt.Println(rs)
}