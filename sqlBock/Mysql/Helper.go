package Mysql

import (
	"os"
	"encoding/json"
	"strconv"
	mysql "database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Query(tablename string, params map[string]interface{}, fields []string, ps ...interface{}) (map[string]interface{}, error) {
	return execute("select * from " + tablename, ps)
}

func execute(sql string, values []interface{}) (map[string]interface{}, error)  {
	var configs interface{}
	fr, err := os.Open("./configs.json")
	decoder := json.NewDecoder(fr)
	err = decoder.Decode(&configs)

	confs := configs.(map[string]interface{})
	dialect := confs["database_dialect"].(string)

	dbConf := confs["db_"+dialect+"_config"].(map[string]interface{})
	dbHost := dbConf["db_host"].(string)
	dbPort := strconv.FormatFloat(dbConf["db_port"].(float64), 'f', -1, 32)
	dbUser := dbConf["db_user"].(string)
	dbPass := dbConf["db_pass"].(string)
	dbName := dbConf["db_name"].(string)
	dbCharset := dbConf["db_charset"].(string)

	rs := make(map[string]interface{})
	dao, err := mysql.Open(dialect, dbUser + ":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)
	defer dao.Close()
	if err != nil {
		rs["code"] = 204
		return rs, err
	}
	stmt, err := dao.Prepare(sql)
	if err != nil {
		rs["code"] = 204
		return rs, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		rs["code"] = 204
		return rs, err
	}

	columns, err := rows.Columns()
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
	return rs, err
}

