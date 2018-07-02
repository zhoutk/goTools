package Mysql

type Bock struct {
	Table string
}

func (b *Bock) Retrieve(params map[string]interface{}, args ...interface{}) (rs map[string]interface{}, err error) {
	_, fields, _ := parseArgs(args)
	return Query(b.Table, params, fields)
}

func parseArgs(args ...interface{}) (string, []string, map[string]interface{}) {
	var id string
	var fields []string
	var session map[string]interface{}
	for _, v := range args {
		switch v.(type) {
		case map[string]interface{}:
			arg := v.(map[string]interface{})
			if v, ok := arg["id"]; ok {
				id = v.(string)
			}
			if v, ok := arg["fields"]; ok {
				fields = v.([]string)
			}
			if v, ok := arg["session"]; ok {
				session = v.(map[string]interface{})
			}
		default:
		}
	}
	return id, fields, session
}