package db

import (
	_ "github.com/go-sql-driver/mysql"
	mysql "database/sql"
)

func Query(sql string, values []interface{}) (map[string]interface{}, error) {
	return execute(sql, values)
}

func execute(sql string, values []interface{}) (map[string]interface{}, error)  {
	rs := make(map[string]interface{})
	dao, err := mysql.Open("mysql", "root:znhl2017UP@tcp(tlwl2020.mysql.rds.aliyuncs.com:3686)/fbox?charset=utf8mb4")
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
	rs["rows"] = result
	return rs, err
}
