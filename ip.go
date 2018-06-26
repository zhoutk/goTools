package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"regexp"
	"./ipSpider"
	"fmt"
)

func main()  {
	res, err := http.Get("http://ips.chacuo.net/")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile(`<li><a title="[\S]+" href='([^']+?)'>([^<]+?)</a></li>`)
	ss := reg.FindAllStringSubmatch(string(robots), -1)

	ch := make(chan string)

	for _, el := range ss {
		go ipSpider.SpiderOnPage(el[1], el[2], ch)
	}

	for range ss {
		fmt.Println( <- ch )
	}
}
