package sqlBock

type IBock interface{
	Retrieve(params map[string]interface{}, args ...interface{}) map[string]interface{}
	Create(params map[string]interface{}, args ...interface{}) map[string]interface{}
	Update(params map[string]interface{}, args ...interface{}) map[string]interface{}
	Delete(params map[string]interface{}, args ...interface{}) map[string]interface{}
	QuerySql(sql string, values []interface{}, params map[string]interface{}) map[string]interface{}
	ExecSql(sql string, values []interface{}) map[string]interface{}
	InsertBatch(tablename string, els []map[string]interface{}) map[string]interface{}
	TransGo(objs map[string]interface{}) map[string]interface{}
}

