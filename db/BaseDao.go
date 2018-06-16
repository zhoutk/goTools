package db

type BaseDao struct {
	Table string
}

func (b *BaseDao) Retrieve(sql string, values [] interface{}) (map[string]interface{}, error) {
	rs := make(map[string]interface{})
	rs["code"] = "200"
	rs["table"] = sql
	return rs, nil
}