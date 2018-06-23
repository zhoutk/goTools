package backUp

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"time"
	"../common"
	"../db"
)

func exportFuncs(fileName string, fields common.DbConnFields) error {
	sqlStr := "select name,type,param_list,returns,body from mysql.proc where db = ? "
	values := make([]interface{}, 0)
	values = append(values, fields.DbName)

	rs, err := db.ExecuteWithDbConn(sqlStr, values, fields)
	fRs := rs["rows"].([]map[string]interface{})

	for _, cstAl := range fRs {
		var rets string
		if cstAl["returns"] != nil && len(cstAl["returns"].(string)) > 0 {
			rets = " RETURNS " + cstAl["returns"].(string)
		}
		sqlStr = "DROP PROCEDURE IF EXISTS `" + cstAl["name"].(string) + "`;\nDELIMITER ;;\n" +
			"CREATE DEFINER=`root`@`%` " + cstAl["type"].(string) + " `" + cstAl["name"].(string) +
			"`(" + cstAl["param_list"].(string) + ")" + rets + "\n" + cstAl["body"].(string) + "\n" +
			";;\nDELIMITER ;\n\n"
		writeToFile(fileName, sqlStr, true)
	}
	return err
}

func exportViews(fileName string, fields common.DbConnFields) error {
	sqlStr := "select TABLE_NAME, VIEW_DEFINITION from information_schema.VIEWS where TABLE_SCHEMA = ? "
	values := make([]interface{}, 0)
	values = append(values, fields.DbName)
	rs, err := db.ExecuteWithDbConn(sqlStr, values, fields)
	vRs := rs["rows"].([]map[string]interface{})
	ps := make(map[string]string)
	vName := make([]string, 0)
	for _, v := range vRs {
		ps["`"+v["TABLE_NAME"].(string)+"`"] = v["VIEW_DEFINITION"].(string)
		vName = append(vName, "`"+v["TABLE_NAME"].(string)+"`")
	}
	rely1 := processRely(ps, &vName)
	rely := processRely(ps, &rely1)
	for _, al := range rely {
		viewStr := strings.Replace(ps[al], "`"+fields.DbName+"`.", "", -1)
		sqlStr = "DROP VIEW IF EXISTS " + al + ";\n" + "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` " +
			" SQL SECURITY DEFINER VIEW " + al + " AS " + viewStr + ";\n\n"
		writeToFile(fileName, sqlStr, true)
	}
	return err
}

func setSqlHeader(fields common.DbConnFields, fileName string) {
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
}

func processRely(params map[string]string, relyOld *[]string) []string {
	rely := make([]string, 0)
	for _, k := range *relyOld {
		for bl := range params {
			if strings.Index(params[k], bl) > -1 {
				if findInArray(&rely, bl) < 0 {
					if findInArray(&rely, k) < 0 {
						rely = append(rely, bl)
					}else{
						i := findInArray(&rely, k)
						lastStr := make([]string, len(rely) - i)
						copy(lastStr, rely[i:])
						rely = append(rely[:i], bl)
						rely = append(rely, lastStr...)
					}
				}else{
					if findInArray(&rely, k) > -1 {
						i := findInArray(&rely, k)
						j := findInArray(&rely, bl)
						if i < j {
							rely = append(rely[:j], rely[j+1:]...)
							lastStr := make([]string, len(rely) - i)
							copy(lastStr, rely[i:])
							rely = append(rely[:i], bl)
							rely = append(rely, lastStr...)
						}
					}
				}
			}
		}
		if findInArray(&rely, k) < 0 {
			rely = append(rely, k)
		}
	}
	return rely
}

func findInArray(arry*[]string, value string) int{
	if arry == nil {
		return -1
	}else{
		for index, v := range *arry {
			if v == value {
				return index
			}
		}
		return -1
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
	if _, err := fileObj.WriteString(content); err != nil {
		fmt.Println(err)
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
