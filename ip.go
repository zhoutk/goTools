package main

import (
	"regexp"
	"./common"
	"./ipSpider"
	"fmt"
)

func main()  {
	ctx := common.HttpGet("http://ips.chacuo.net/")

	reg := regexp.MustCompile(`<li><a title="[\S]+" href='([^']+?)'>([^<]+?)</a></li>`)
	ips := reg.FindAllStringSubmatch(string(ctx), -1)

	ch := make(chan string)

	for _, el := range ips {
		go ipSpider.SpiderOnPage(el[1], el[2], ch)
	}

	for range ips {
		fmt.Println(<-ch)
	}
}

