package main

import (
	"./db"
	"fmt"
	"encoding/json"
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
		data, _ := json.Marshal(rs)
		fmt.Println(string(data))
	}
}