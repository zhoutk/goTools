package db

type Idb interface{
	Retrieve(sql string, valsues [] interface{}) (rs map[string]interface{}, err error)
}
