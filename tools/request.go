package tools

import (
	"io/ioutil"
	"log"
	"net/http"
)

// 获取请求 url 的 html 页面
func GetRequestByte(url string) []byte {
	res, err := http.Get(url)
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
