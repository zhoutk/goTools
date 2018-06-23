package backUp

import (
	"strings"
	"encoding/json"
	"os"
	"../common"
	"fmt"
)

func Export(flag common.OpFlag) (error) {
	var configs interface{}
	fr, err := os.Open("./configs.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(fr)
	err = decoder.Decode(&configs)
	if err != nil {
		return err
	}
	confs := configs.(map[string]interface{})
	workDir := confs["workDir"].(string)
	ch := make(chan string)
	for key, value := range confs {
		if strings.HasPrefix(key, "db_") {
			dbConf := value.(map[string]interface{})
			dbConn := common.DbConnFields{
				DbHost:    dbConf["db_host"].(string),
				DbPort:    int(dbConf["db_port"].(float64)),
				DbUser:    dbConf["db_user"].(string),
				DbPass:    dbConf["db_pass"].(string),
				DbName:    dbConf["db_name"].(string),
				DbCharset: dbConf["db_charset"].(string),
			}
			if dbConf["file_alias"] != nil {
				dbConn.FileAlias = dbConf["file_alias"].(string)
			}
			go ExportOne(dbConn, workDir, ch, flag)
		}
	}
	for key := range confs {
		if strings.HasPrefix(key, "db_") {
			fmt.Print( <- ch )
		}
	}
	return nil
}
