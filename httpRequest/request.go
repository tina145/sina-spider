package httpRequest

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 通过 GET 请求获取请求 url 的 html 页面
func GetRequestByte(url string) []byte {
	if IsDisallow(url) {
		return nil
	}
	getHtmlClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	// 随即设置一个 useragent
	req.Header.Set("User-Agent", GetRandomUserAgent())

	res, err := getHtmlClient.Do(req)
	if err != nil {
		log.Println(nil)
		return nil
	}

	defer res.Body.Close()
	if err != nil {
		log.Println(err)
		return nil
	} else if res.StatusCode == 404 {
		log.Println("404 请求的页面不存在")
		return nil
	}
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return html
}
