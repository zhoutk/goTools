package backUp

import (
	"os"
	"fmt"
	"time"
	"strconv"
)

func ExportOne(fields DbConnFields, workDir string)  {
	name := workDir + fields.fileAlias + time.Now().Format("2006-01-02") + ".sql"
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		os.Exit(2)
	}
	defer fileObj.Close()
	var content string
	content = "/*   Mysql export" +
		 "\n\n		Host: " + fields.dbHost +
		"\n\n		Port: " + strconv.Itoa(fields.dbPort) +
		"\n\n		DataBase: " + fields.dbName +
		"\n\n		Date: " + time.Now().Format("2006-01-02 15:04:05") +
		"\n\n		Author: zhoutk@189.cn" +
		"\n\n		Copyright: tolw-2018" +
		"\n*/\n\n"
	if _,err := fileObj.WriteString(content);err == nil {
		fmt.Println("Successful writing to the file with .")
	}
}
