package tools

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
)

// 下载资源
func DownloadGet(url string, position ...string) {
	if len(position) > 1 {
		panic("too many args")
	}
	// 如果没有传入参数，则设置一个默认位置下载文件
	pos := ".\\files\\" + strconv.Itoa(rand.Intn(100000))
	if len(position) == 1 {
		// 有传入位置则使用传入的位置
		pos = position[0]
	}

	// Windows 下文件名称不能包含以下字符，去除不合法字符
	obj := regexp.MustCompile(`[/\:*?"<>|]`)
	pos = obj.ReplaceAllString(pos, "")

	// 正则表达式搜索文件所需文件夹是否存在
	obj = regexp.MustCompile(`([\s\S]+)\\`)
	fileroot := obj.FindAllStringSubmatch(pos, -1)[0][1]

	// 递归创建所有所需文件夹
	os.MkdirAll(fileroot, 0644)

	data := GetRequestByte(url)

	err := ioutil.WriteFile(pos, data, 0644)
	if err != nil {
		log.Println(err)
		return
	}
}

// 下载 html 页面中所有资源，参数为 url 地址和 保存的文件位置
func DownloadHtmlSource(url, rule string, position ...string) {
	// 参数检测，传入文件位置超过一个则 panic
	if len(position) > 1 {
		panic("too many args")
	}

	pos := ".\\files\\" + strconv.Itoa(rand.Intn(100000))
	if len(position) == 1 {
		pos = position[0]
	}

	// 资源文件匹配规则
	sourceRule := rule
	sources := RegexpHtml(url, sourceRule)
	for _, urls := range sources {
		DownloadGet(urls, pos+urls)
	}
}
