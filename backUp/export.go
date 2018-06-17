package backUp

import (
	"strings"
	"encoding/json"
	"os"
)

type DbConnFields struct {
	dbHost    string
	dbPort    int
	dbUser    string
	dbPass    string
	dbName    string
	dbCharset string
	fileAlias string
}

func Export() (bool, error) {
	var configs interface{}
	fr, err := os.Open("./configs.json")
	if err != nil {
		return false, err
	}
	decoder := json.NewDecoder(fr)
	err = decoder.Decode(&configs)
	if err != nil {
		return false, err
	}
	confs := configs.(map[string]interface{})
	workDir := confs["workDir"].(string)
	for key, value := range confs {
		if strings.HasPrefix(key, "db_") {
			dbConf := value.(map[string]interface{})
			dbConn := DbConnFields{
				dbHost:    dbConf["db_host"].(string),
				dbPort:    int(dbConf["db_port"].(float64)),
				dbUser:    dbConf["db_user"].(string),
				dbPass:    dbConf["db_pass"].(string),
				dbName:    dbConf["db_name"].(string),
				dbCharset: dbConf["db_charset"].(string),
			}
			if dbConf["file_alias"] != nil {
				dbConn.fileAlias = dbConf["file_alias"].(string)
			}
			ExportOne(dbConn, workDir)
		}
	}
	return true, nil
}
