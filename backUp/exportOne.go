package backUp

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"../common"
	"../db"
	"encoding/json"
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
	var values []interface{}
	values = append(values, fields.DbName)
	rs, err := db.ExecuteWithDbConn(sqlStr, values, fields)
	if err != nil{
		fmt.Print(err)
		return
	}
	rows := rs["rows"].([]map[string]string)
	FKEYS := make(map[string]interface{})
	for i := 0; i < len(rows); i++ {
		if _, ok := FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]]; !ok {
			FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]] = map[string]interface{}{
				"constraintName": rows[i]["CONSTRAINT_NAME"],
				"sourceCols":     make([]interface{}, 0),
				"schema":         rows[i]["REFERENCED_TABLE_SCHEMA"],
				"tableName":      rows[i]["REFERENCED_TABLE_NAME"],
				"targetCols":     make([]interface{}, 0),
			}
		}
		FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]].(map[string]interface{})["sourceCols"] =
			append(FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]].(map[string]interface{})["sourceCols"].([]interface{}), rows[i]["COLUMN_NAME"])
		FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]].(map[string]interface{})["targetCols"] =
			append(FKEYS[rows[i]["TABLE_NAME"]+"."+rows[i]["CONSTRAINT_NAME"]].(map[string]interface{})["targetCols"].([]interface{}), rows[i]["COLUMN_NAME"])
	}
	//data, _ := json.Marshal(FKEYS)
	//fmt.Print(string(data))

	sqlStr = "select TABLE_NAME,ENGINE,ROW_FORMAT,AUTO_INCREMENT,TABLE_COLLATION,CREATE_OPTIONS,TABLE_COMMENT" +
		" from information_schema.`TABLES` where TABLE_SCHEMA = ? and TABLE_TYPE = ? order by TABLE_NAME"
	values = make([]interface{}, 0)
	values = append(values, fields.DbName, "BASE TABLE")
	rs, err = db.ExecuteWithDbConn(sqlStr, values, fields)
	if err != nil{
		fmt.Print(err)
		return
	}
	tbRs := rs["rows"].([]map[string]string)
	data, _ := json.Marshal(tbRs)
	fmt.Print(string(data))
	for _, tbAl := range tbRs{
		sqlStr = "SELECT	`COLUMNS`.COLUMN_NAME,`COLUMNS`.COLUMN_TYPE,`COLUMNS`.IS_NULLABLE," +
					"`COLUMNS`.CHARACTER_SET_NAME,`COLUMNS`.COLUMN_DEFAULT,`COLUMNS`.EXTRA," +
                    "`COLUMNS`.COLUMN_KEY,`COLUMNS`.COLUMN_COMMENT,`STATISTICS`.TABLE_NAME," +
                    "`STATISTICS`.INDEX_NAME,`STATISTICS`.SEQ_IN_INDEX,`STATISTICS`.NON_UNIQUE," +
                    "`COLUMNS`.COLLATION_NAME " +
                    "FROM information_schema.`COLUMNS` " +
                    "LEFT JOIN information_schema.`STATISTICS` ON " +
                    "information_schema.`COLUMNS`.TABLE_NAME = `STATISTICS`.TABLE_NAME " +
                    "AND information_schema.`COLUMNS`.COLUMN_NAME = information_schema.`STATISTICS`.COLUMN_NAME " +
                    "AND information_schema.`STATISTICS`.table_schema = ? " +
                    "where information_schema.`COLUMNS`.TABLE_NAME = ? and `COLUMNS`.table_schema = ?"
		values = make([]interface{}, 0)
		values = append(values, fields.DbName, tbAl["TABLE_NAME"],fields.DbName)
		rs, err = db.ExecuteWithDbConn(sqlStr, values, fields)
		colRs := rs["rows"].([]map[string]string)
		//data, _ := json.Marshal(colRs)
		//fmt.Print(string(data))
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