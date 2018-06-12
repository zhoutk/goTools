package main

import (
	"fmt"
	"os"
	"strconv"
	"github.com/Luxurioust/excelize"
	"encoding/json"
	"math/rand"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

func init(){
	//以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}


func main() {
	xlsx, err := excelize.OpenFile("./qs.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get value from cell by given sheet index and axis.
	//cell := xlsx.GetCellValue("Sheet1", "B2")
	//fmt.Println(cell)
	// Get sheet index.
	index := xlsx.GetSheetIndex("Sheet1")
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("Sheet" + strconv.Itoa(index))
	var vs string
	var vss string
	for i, row := range rows {
		fmt.Println(i)
		if i == 0 {
			continue
		}
		var answers []map[string]interface{}
		answer := make(map[string]interface{})

		answer["right"] = true
		answer["answer"] = row[3]
		answers = append(answers, answer)

		answer = make(map[string]interface{})
		answer["right"] = false
		answer["answer"] = row[4]
		answers = append(answers, answer)

		answer = make(map[string]interface{})
		answer["right"] = false
		answer["answer"] = row[5]
		answers = append(answers, answer)

		answer = make(map[string]interface{})
		answer["right"] = false
		answer["answer"] = row[6]
		answers = append(answers, answer)

		for j := 0; j < 10; j++ {
			sp := rand.Intn(4)
			tmp := answers[sp]
			answers[sp] = answers[0]
			answers[0] = tmp
		}

		data, err := json.Marshal(answers)
		if err == nil {
			dd := string(data)
			vs = vs + " ('" + row[1]+ "', '" + dd + "' ),"
		}
		vss += "(?,?),"
	}
	vs = vs[0:len(vs) -1]
	vss = vss[0:len(vss) -1]
	db, err := sql.Open("mysql", "root:znhl2017UP@tcp(tlwl2020.mysql.rds.aliyuncs.com:3686)/policy?charset=utf8")
	sqlstr := "insert into questions2 (name, answer_json) values " + vs
	defer db.Close()
	fmt.Printf("%s\n", sqlstr)
	res, err := db.Exec(sqlstr)
	fmt.Println(res)
	fmt.Println(err)
}