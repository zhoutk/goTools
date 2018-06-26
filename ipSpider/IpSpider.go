package ipSpider

import (
	"../common"
	"regexp"
)

func SpiderOnPage(url string, province string, ch chan <- string) error {
	ctx := common.HttpGet(url)
	//reg := regexp.MustCompile(`<li><a title="[\S]+" href='([^']+?)'>([^<]+?)</a></li>`)
	reg := regexp.MustCompile(`<dd><span class="v_l">([^<]+?)</span><span class="v_r">([^<]+?)</span><div class="clearfix"></div></dd>`)
	//<dd><span class="v_l">49.64.0.0</span><span class="v_r">49.95.255.255</span><div class="clearfix"></div></dd>
	ip := reg.FindAllStringSubmatch(string(ctx), -1)

	var vs [] interface{}
	for _, el := range ip {
		vs = append(vs, el[1], el[2], province)
	}

	ch <- url + " : " + ip[0][1] +","+ ip[0][2]
	return nil
}
