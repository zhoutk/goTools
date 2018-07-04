package main

import (
	Bock "./sqlBock/Mysql"
	"fmt"
)

func main()  {
	table := Bock.Bock{
		Table: "books",
	}

	/*//insert one record
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
	params["status"] = 3
	rs = db.Update(params, args)
	fmt.Println(rs)

	//delete one record
	args = make(map[string] interface{})
	args["id"] = 6
	rs = db.Delete(nil, args)
	fmt.Println(rs)

	//execSql
	values := make([]interface{}, 0)
	values = append(values, "我是手写sql")
	values = append(values, 1)
	rs = db.ExecSql("update books set name = ? where id = ? ", values)
	fmt.Println(rs)
	*/

	//insertBatch
	vs := make([]map[string]interface{}, 0)

	params := make(map[string] interface{})
	params["name"] = "golang批量11213"
	params["isbn"] = "4156s5"
	params["status"] = 5
	params["id"] = 9
	vs = append(vs, params)

	params = make(map[string] interface{})
	params["name"] = "golang批量22af24"
	params["isbn"] = "xxfqwt325rqrf45"
	params["status"] = 2
	params["id"] = 10
	vs = append(vs, params)

	db := &table
	rs := db.InsertBatch("books", vs)
	fmt.Println(rs)
}
