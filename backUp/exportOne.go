package backUp

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"../common"
	"../db"
	"reflect"
)

func ExportOne(fields common.DbConnFields, workDir string) {
	var fileName string
	if fields.FileAlias == "" {
		fileName = workDir + fields.DbName + "-" + time.Now().Format("2006-01-02") + ".sql"
	}else{
		fileName = workDir + fields.FileAlias + "-" + time.Now().Format("2006-01-02") + ".sql"
	}

	content := "/*   Mysql export" +
		"\n\n		Host: " + fields.DbHost +
		"\n\n		Port: " + strconv.Itoa(fields.DbPort) +
		"\n\n		DataBase: " + fields.DbName +
		"\n\n		Date: " + time.Now().Format("2006-01-02 15:04:05") +
		"\n\n		Author: zhoutk@189.cn" +
		"\n\n		Copyright: tlwl-2018" +
		"\n*/\n\n"
	writeToFile(fileName, content, false)
	writeToFile(fileName, "SET FOREIGN_KEY_CHECKS=0;\n\n", true)
	sqlStr := "select CONSTRAINT_NAME,TABLE_NAME,COLUMN_NAME,REFERENCED_TABLE_SCHEMA," +
		"REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME from information_schema.`KEY_COLUMN_USAGE` " +
		"where REFERENCED_TABLE_SCHEMA = ? "
	values := make([]interface{}, 0)
	values = append(values, fields.DbName)
	rs, err := db.ExecuteWithDbConn(sqlStr, values, fields)
	if err != nil{
		fmt.Print(err)
		return
	}
	for key, value := range rs{
		fmt.Println(reflect.TypeOf(value).String())
		fmt.Println(key, value)
	}
}

func writeToFile(name string, content string, append bool)  {
	var fileObj *os.File
	var err error
	if append{
		fileObj, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	}else{
		fileObj, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}
	if err != nil {
		fmt.Println("Failed to open the file", err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	if _, err := fileObj.WriteString(content); err == nil {
		fmt.Println("Successful writing to the file.")
	}
}