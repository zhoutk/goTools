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
	db := &table
	values = append(values, 548)
	rs, err := db.Retrieve("select * from aaa where index_id = ? ", values)
	if err != nil{
		fmt.Println(err)
	}else {
		fmt.Println(rs)
	}
}