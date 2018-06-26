package ipSpider

func SpiderOnPage(url string, province string, ch chan <- string) error {
	ch <- url + " : " + province
	return nil
}
