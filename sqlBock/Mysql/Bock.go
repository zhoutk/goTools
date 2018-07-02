package Mysql

type Bock struct {
	Table string
}

func (b *Bock) Retrieve(params map[string]interface{}, args ...interface{}) map[string]interface{} {
	_, fields, _ := parseArgs(args)
	return Query(b.Table, params, fields)
}

func (b *Bock) Create(params map[string]interface{}, args ...interface{}) map[string]interface{} {
	_, _, session := parseArgs(args)
	uId := session["userid"].(string)
	params["u_id"] = uId
	return Insert(b.Table, params)
}

func (b *Bock) Update(params map[string]interface{}, args ...interface{}) map[string]interface{} {
	id, _, _ := parseArgs(args)
	if len(id) == 0 {
		rs := make(map[string]interface{})
		rs["code"] = 301
		rs["err"] = "Id must be input."
		return rs
	}
	return Update(b.Table, params)
}

func (b *Bock) Delete(params map[string]interface{}, args ...interface{}) map[string]interface{} {
	id, _, _ := parseArgs(args)
	if len(id) == 0 {
		rs := make(map[string]interface{})
		rs["code"] = 301
		rs["err"] = "Id must be input."
		return rs
	}
	return Delete(b.Table, params)
}

func parseArgs(args []interface{}) (string, []string, map[string]interface{}) {
	var id string
	var fields []string
	var session map[string]interface{}
	for _, vs := range args {
		switch vs.(type) {
		case map[string]interface{}:
			for k, v := range vs.(map[string]interface{}) {
				if k == "id" {
					id = v.(string)
				}
				if k == "fields" {
					fields = v.([]string)
				}
				if k == "session" {
					session = v.(map[string]interface{})
				}
			}
		default:
		}
	}
	return id, fields, session
}