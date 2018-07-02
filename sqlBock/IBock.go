package sqlBock

type IBock interface{
	Retrieve(params map[string]interface{}, args ...interface{}) (rs map[string]interface{}, err error)
}

