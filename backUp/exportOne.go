package backUp

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"../common"
	"../db"
	"strings"
)

func ExportOne(fields common.DbConnFields, workDir string) {
	var fileName string
	if fields.FileAlias == "" {
		fileName = workDir + fields.DbName + "-" + time.Now().Format("2006-01-02") + ".sql"
	}else{
		fileName = workDir + fields.FileAlias + "-" + time.Now().Format("2006-01-02") + ".sql"
	}

	content := "/*   Mysql export \n" +
		"\n		 Host: " + fields.DbHost +
		"\n		 Port: " + strconv.Itoa(fields.DbPort) +
		"\n		 DataBase: " + fields.DbName +
		"\n		 Date: " + time.Now().Format("2006-01-02 15:04:05") +
		"\n\n		 Author: zhoutk@189.cn" +
		"\n		 Copyright: tlwl-2018" +
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
	rows := rs["rows"].([]map[string]interface{})
	FKEYS := make(map[string]interface{})
	for i := 0; i < len(rows); i++ {
		if _, ok := FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)]; !ok {
			FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)] = map[string]interface{}{
				"constraintName": rows[i]["CONSTRAINT_NAME"],
				"sourceCols":     make([]string, 0),
				"schema":         rows[i]["REFERENCED_TABLE_SCHEMA"],
				"tableName":      rows[i]["REFERENCED_TABLE_NAME"],
				"targetCols":     make([]string, 0),
			}
		}
		FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)].(map[string]interface{})["sourceCols"] =
			append(FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)].(map[string]interface{})["sourceCols"].([]string), rows[i]["COLUMN_NAME"].(string))
		FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)].(map[string]interface{})["targetCols"] =
			append(FKEYS[rows[i]["TABLE_NAME"].(string)+"."+rows[i]["CONSTRAINT_NAME"].(string)].(map[string]interface{})["targetCols"].([]string), rows[i]["REFERENCED_COLUMN_NAME"].(string))
	}

	sqlStr = "select TABLE_NAME,ENGINE,ROW_FORMAT,AUTO_INCREMENT,TABLE_COLLATION,CREATE_OPTIONS,TABLE_COMMENT" +
		" from information_schema.`TABLES` where TABLE_SCHEMA = ? and TABLE_TYPE = ? order by TABLE_NAME"
	values = make([]interface{}, 0)
	values = append(values, fields.DbName, "BASE TABLE")
	rs, err = db.ExecuteWithDbConn(sqlStr, values, fields)
	if err != nil{
		fmt.Print(err)
		return
	}
	tbRs := rs["rows"].([]map[string]interface{})
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
		colRs := rs["rows"].([]map[string]interface{})
		tableName := tbAl["TABLE_NAME"].(string)
		tableEngine := tbAl["ENGINE"].(string)
		//tableRowFormat := tbAl["ROW_FORMAT"]
		tableAutoIncrement := tbAl["AUTO_INCREMENT"].(string)
		tableCollation := tbAl["TABLE_COLLATION"].(string)
		tableCharset := strings.Split(tableCollation, "_")[0]
		tableCreateOptions := tbAl["CREATE_OPTIONS"].(string)
		tableComment := tbAl["TABLE_COMMENT"].(string)

		strExport := "DROP TABLE IF EXISTS `" + tbAl["TABLE_NAME"].(string) + "`;\n"
		strExport += "CREATE TABLE `" + tableName + "` (\n"

		priKey := make(map[string]interface{})
		colKey := make(map[string]interface{})
		mulKey := make(map[string]interface{})
		theTableColSet := make(map[string]int)
		var allFields []string
		var defaultValue string
		for _, colAl := range colRs{
			if _, ok := theTableColSet[colAl["COLUMN_NAME"].(string)]; !ok {
				theTableColSet[colAl["COLUMN_NAME"].(string)] = 1
				allFields = append(allFields, colAl["COLUMN_NAME"].(string))
				if colAl["COLUMN_DEFAULT"] != nil && len(colAl["COLUMN_DEFAULT"].(string)) > 0{
					if colAl["COLUMN_DEFAULT"] == "CURRENT_TIMESTAMP"{
						defaultValue = colAl["COLUMN_DEFAULT"].(string)
					}else{
						defaultValue = "'" + colAl["COLUMN_DEFAULT"].(string) + "'"
					}
				}
				var charSet string
				if colAl["CHARACTER_SET_NAME"] != nil && colAl["CHARACTER_SET_NAME"] != tableCharset {
					charSet = " CHARACTER SET " + colAl["CHARACTER_SET_NAME"].(string)
				}
				var collation string
				if colAl["COLLATION_NAME"] != nil && colAl["COLLATION_NAME"] != tableCollation {
					collation = " COLLATE " + colAl["COLLATION_NAME"].(string)
				}
				var nullStr string
				if colAl["IS_NULLABLE"] != nil && colAl["IS_NULLABLE"] == "NO" {
					nullStr = " NOT NULL"
				}
				if colAl["COLUMN_DEFAULT"] != nil && len(colAl["COLUMN_DEFAULT"].(string)) > 0 {
					defaultValue = " DEFAULT " + defaultValue
				}else{
					if colAl["IS_NULLABLE"] != nil && colAl["IS_NULLABLE"] == "NO"{
						defaultValue = ""
					}else{
						defaultValue = " DEFAULT NULL"
					}
				}
				var space string
				if colAl["EXTRA"] != nil && len(colAl["EXTRA"].(string)) > 0 {
					space = " " + colAl["EXTRA"].(string)
				}else{
					space = ""
				}
				var cstr string
				if colAl["COLUMN_COMMENT"] != nil && len(colAl["COLUMN_COMMENT"].(string)) > 0 {
					cstr = " COMMENT '" + escape(colAl["COLUMN_COMMENT"].(string)) + "'"
				}
				strExport += "  `" + colAl["COLUMN_NAME"].(string) + "` " + colAl["COLUMN_TYPE"].(string) + charSet + collation +
					nullStr + defaultValue + space + cstr + ",\n"
			}
			if colAl["INDEX_NAME"] != nil && colAl["INDEX_NAME"].(string) == "PRIMARY" {
				if _, ok := priKey[colAl["INDEX_NAME"].(string)]; !ok {
					priKey[colAl["INDEX_NAME"].(string)] = make([]string,0)
				}
				priKey[colAl["INDEX_NAME"].(string)] = append(priKey[colAl["INDEX_NAME"].(string)].([]string), colAl["COLUMN_NAME"].(string))
			}else if colAl["INDEX_NAME"] != nil && colAl["NON_UNIQUE"] == "0" {
				if _, ok := colKey[colAl["INDEX_NAME"].(string)]; !ok {
					colKey[colAl["INDEX_NAME"].(string)] = make([]string,0)
				}
				colKey[colAl["INDEX_NAME"].(string)] = append(colKey[colAl["INDEX_NAME"].(string)].([]string), colAl["COLUMN_NAME"].(string))
			}else if colAl["INDEX_NAME"] != nil && colAl["NON_UNIQUE"] == "1" {
				if _, ok := mulKey[colAl["INDEX_NAME"].(string)]; !ok {
					mulKey[colAl["INDEX_NAME"].(string)] = make([]string,0)
				}
				mulKey[colAl["INDEX_NAME"].(string)] = append(mulKey[colAl["INDEX_NAME"].(string)].([]string), colAl["COLUMN_NAME"].(string))
			}
		}
		for _, v := range priKey {
			strExport += "  PRIMARY KEY (`" + strings.Join(v.([]string), "`,`") + "`),\n"
		}
		for k, v := range colKey {
			strExport += "  UNIQUE KEY `"+k+"` (`" + strings.Join(v.([]string), "`,`") + "`),\n"
		}
		for k, v := range mulKey {
			strExport += "  KEY `"+k+"` (`" + strings.Join(v.([]string), "`,`") + "`),\n"
		}

		for k, v := range FKEYS {
			if strings.HasPrefix(k, tableName+".") {
				strExport += "  CONSTRAINT `" + v.(map[string]interface{})["constraintName"].(string) + "` FOREIGN KEY (`" +
					strings.Join(v.(map[string]interface{})["sourceCols"].([]string), "`,`") + "`) REFERENCES `" +
					v.(map[string]interface{})["tableName"].(string) + "` (`" +
					strings.Join(v.(map[string]interface{})["targetCols"].([]string), "`,`") + "`),\n"
			}
		}
		if strings.HasSuffix(strExport, ",\n") {
			strExport = strExport[:len(strExport)-2]
		}

		var incr string
		if len(tableAutoIncrement) > 0 {
			incr = " AUTO_INCREMENT=" + tableAutoIncrement
		}
		var colla string
		if len(tableCollation) > 0 {
			colla = " COLLATE=" + tableCollation
		}
		strExport += "\n) ENGINE=" + tableEngine + incr + " DEFAULT CHARSET=" +
			tableCharset + colla + " " + tableCreateOptions + " COMMENT='" + tableComment + "';\n\n"

		writeToFile(fileName, strExport, true)			//表结构导出

		//表中数据导出
		sqlStr = "select " + strings.Join(allFields, ",") + " from " + tableName
		rs, err = db.ExecuteWithDbConn(sqlStr, make([]interface{}, 0), fields)
		recordsRs := rs["rows"].([]map[string]interface{})
		for _, ele := range recordsRs {
			strExport = "INSERT INTO `" + tableName + "` ("//+strings.Join(allFields, ",")+") VALUES ("
			var ks []string
			var vs []string
			for k, v := range ele {
				ks = append(ks, k)
				elStr := "''"
				if v == nil {
					elStr = "null"
				}else if len(v.(string)) > 0 {
					elStr = "'" + escape(v.(string)) + "'"
				}
				vs = append(vs, elStr)
			}
			strExport += strings.Join(ks, ",") + ") VALUES (" + strings.Join(vs, ",") + ");\n"
			writeToFile(fileName, strExport, true)
		}
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

func escape(source string) string {
	var j int
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
		case '\n':
			flag = true
			escape = '\n'
		case '\\':
			flag = true
			escape = '\\'
		case '\'':
			flag = true
			escape = '\''
		case '"':
			flag = true
			escape = '"'
		case '\032':
			flag = true
			escape = 'Z'
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}