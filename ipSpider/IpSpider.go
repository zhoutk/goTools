package ipSpider

import (
	"../common"
	"regexp"
	"os"
	"encoding/json"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	mysql "database/sql"
)

func SpiderOnPage(url string, province string, ch chan <- string) error {
	ctx := common.HttpGet(url)
	//reg := regexp.MustCompile(`<li><a title="[\S]+" href='([^']+?)'>([^<]+?)</a></li>`)
	reg := regexp.MustCompile(`<dd><span class="v_l">([^<]+?)</span><span class="v_r">([^<]+?)</span><div class="clearfix"></div></dd>`)
	//<dd><span class="v_l">49.64.0.0</span><span class="v_r">49.95.255.255</span><div class="clearfix"></div></dd>
	ip := reg.FindAllStringSubmatch(string(ctx), -1)

	if len(ip) == 0 {
		ch <- "There are no data exist."
		return nil
	}

	var vs [] interface{}
	var vss string
	for _, el := range ip {
		vs = append(vs, el[1], el[2], province)
		vss += "(?,?,?),"
	}
	vss = vss[0:len(vss) -1]
	var configs interface{}
	fr, err := os.Open("./configs.json")
	if err != nil {
		ch <- err.Error()
		return err
	}
	decoder := json.NewDecoder(fr)
	err = decoder.Decode(&configs)

	confs := configs.(map[string]interface{})
	dialect := confs["database_dialect"].(string)

	dbConf := confs["db_"+dialect+"_config"].(map[string]interface{})
	dbHost := dbConf["db_host"].(string)
	dbPort := strconv.FormatFloat(dbConf["db_port"].(float64), 'f', -1, 32)
	dbUser := dbConf["db_user"].(string)
	dbPass := dbConf["db_pass"].(string)
	dbName := dbConf["db_name"].(string)
	dbCharset := dbConf["db_charset"].(string)

	dao, err := mysql.Open(dialect, dbUser + ":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)
	defer dao.Close()
	if err != nil {
		ch <- err.Error()
		return err
	}
	sqlstr := "insert into ip_addr_info (ip_addr,ip_mask,ip_comp) values " + vss +
		" ON DUPLICATE KEY UPDATE ip_addr = values(ip_addr), ip_mask = values(ip_mask), ip_comp = values(ip_comp)"
	stmt, err := dao.Prepare(sqlstr)
	rs, err := stmt.Exec(vs...)
	if err != nil {
		ch <- err.Error()
		return err
	}else {
		affect, _ := rs.RowsAffected()
		ch <- "Province: " + province + ", affect: " + strconv.FormatInt(affect, 10)
		return nil
	}

}
