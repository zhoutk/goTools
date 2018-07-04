package Mysql

import (
	"os"
	"encoding/json"
	"strconv"
	mysql "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func query(tablename string, params map[string]interface{}, fields []string, sql string, vaules []interface{}) map[string]interface{} {
	if vaules == nil {
		vaules = make([]interface{},0)
	}
	rs := execQeury("select "+ strings.Join(fields, ",")+" from " + tablename, vaules)
	return rs
}

func Query(tablename string, params map[string]interface{}, fields []string ) map[string]interface{} {
	return query(tablename, params, fields, "", nil)
}

func Insert(tablename string, params map[string]interface{}) map[string]interface{} {
	sql := "Insert into " + tablename
	values := make([]interface{},0)
	return execute(sql, values)
}

func Update(tablename string, params map[string]interface{}) map[string]interface{} {
	sql := "Update " + tablename + " set "
	values := make([]interface{},0)
	return execute(sql, values)
}

func Delete(tablename string, params map[string]interface{}) map[string]interface{} {
	sql := "Delete from " + tablename + " where"
	values := make([]interface{},0)
	return execute(sql, values)
}

func execute(sql string, values []interface{}) map[string]interface{}  {
	rs := make(map[string]interface{})
	rs["code"] = 200
	return rs
}

func execQeury(sql string, values []interface{}) (rs map[string]interface{})  {
	var configs interface{}
	rs = make(map[string]interface{})
	defer func() {
		if r := recover(); r != nil {
			rs["code"] = 500
			rs["err"] = "Exception, " + r.(error).Error()
		}
	}()

	fr, err := os.Open("./configs.json")
	decoder := json.NewDecoder(fr)
	err = decoder.Decode(&configs)
	if err != nil {
		rs["code"] = 204
		rs["err"] = "Open database config error, " + err.Error()
		return rs
	}

	confs := configs.(map[string]interface{})
	dialect := confs["database_dialect"].(string)

	dbConf := confs["db_"+dialect+"_config"].(map[string]interface{})
	dbHost := dbConf["db_host"].(string)
	dbPort := strconv.FormatFloat(dbConf["db_port"].(float64), 'f', -1, 32)
	dbUser := dbConf["db_user"].(string)
	dbPass := dbConf["db_pass"].(string)
	dbName := dbConf["db_name"].(string)
	dbCharset := dbConf["db_charset"].(string)

	dao, err := mysql.Open(dialect, dbUser + ":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)
	defer dao.Close()
	if err != nil {
		rs["code"] = 204
		rs["err"] = err.Error()
		return rs
	}
	stmt, err := dao.Prepare(sql)
	if err != nil {
		rs["code"] = 204
		return rs
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		rs["code"] = 204
		rs["err"] = err.Error()
		return rs
	}

	columns, err := rows.Columns()
	if err != nil {
		rs["code"] = 204
		rs["err"] = err.Error()
		return rs
	}
	vs := make([]mysql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))

	for i := range vs {
		scans[i] = &vs[i]
	}

	var result []map[string]string
	for rows.Next() {
		_ = rows.Scan(scans...)
		each := make(map[string]string)

		for i, col := range vs {
			each[columns[i]] = string(col)
		}

		result = append(result, each)
	}
	rs["code"] = 200
	//data, _ := json.Marshal(result)
	rs["rows"] = result
	return rs
}



