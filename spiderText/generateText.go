package spiderText

import (
	"GoProject/spider/httpRequest"
	"log"
	"regexp"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mysqlAccount string = "账号密码等..."

type urls struct {
	Url string `db:"url"`
}

type news struct {
	Url   string `db:"url"`
	Title string `db:"title"`
}

// 返回文章链接和标题
func Search() [][]string {
	// 连接 MySQL
	db, _ := sqlx.Open("mysql", mysqlAccount)
	defer db.Close()

	html := httpRequest.GetRequestByte("https://finance.sina.com.cn/stock/")

	// 匹配规则
	rule := `href="(https://finance.sina.com.cn/stock/[\S]+?html[\S]*?)">[\S]+?</a>`
	obj := regexp.MustCompile(rule)

	arr := obj.FindAllStringSubmatch(string(html), -1)

	// 查看 url 是否有重复
	datas := make([]urls, 0)
	db.Select(&datas, "SELECT url FROM urlinfo")
	rec := make(map[string]bool)
	for _, data := range datas {
		rec[data.Url] = true
	}

	for index := range arr {
		if _, ok := rec[arr[index][1]]; ok {
			// 已经存储的直接跳过
			continue
		}

		arr[index] = append(arr[index], httpRequest.RegexpHtml(arr[index][1], `<title>([\s\S]+?)</title>`)[0])
	}

	return arr
}

// 生成正文
func GenerateText(arr [][]string) string {
	db, _ := sqlx.Open("mysql", mysqlAccount)
	defer db.Close()
	text := ``

	for _, data := range arr {
		if len(data) < 3 {
			continue
		}
		text += `<h2>
		<a target="_blank" href="` + data[1] + `">` + data[2] + `</a>
		<h2>`

		// 保存到数据库中
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return ""
		}
		_, err = tx.Exec("INSERT INTO urlinfo values(?,?,?,?,?)", 0, "https://finance.sina.com.cn/stock/", data[1], 0, data[2])

		if err != nil {
			log.Println(err)
			return ""
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return ""
		}

		// SaveRedis(data[1])
	}

	return text
}

// 查找缓存是否存在
func FindFromCache(key, member string) bool {
	connect, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer connect.Close()
	reply, _ := redis.Bool(connect.Do("SISMEMBER", key, member))
	return reply
}

// 存储到缓存中
func SaveRedis(member ...string) string {
	connect, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer connect.Close()

	key := time.Now().String()
	reply, _ := redis.String(connect.Do("SADD", key, member))

	// 设置一天后过期
	connect.Do("expire", "key", 86400)
	return reply
}

func SelectFirst20() string {
	db := sqlx.MustOpen("mysql", mysqlAccount)
	defer db.Close()
	news := []news{}

	db.Select(&news, "SELECT url, title FROM urlinfo ORDER BY ID DESC LIMIT 20")

	text := ``

	for _, data := range news {
		text += `<h2>
		<a target="_blank" href="` + data.Url + `">` + data.Title + `</a>
		<h2>`
	}

	return text
}
