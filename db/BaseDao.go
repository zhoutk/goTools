package db

type BaseDao struct {
	Table string
}

func (b *BaseDao) Retrieve(sql string, values [] interface{}) (map[string]interface{}, error) {
	return Query(sql, values)
}