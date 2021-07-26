package functions

import (
	"math/rand"
	"net/http"
	"project/Mail"
	"project/Text"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

func ToFuntions(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "function.html", nil)
}

func SendStock(ctx *gin.Context) {
	cookie, err := ctx.Cookie("cookie")
	if err != nil {
		ctx.String(http.StatusBadRequest, "请先登录")
		return
	}
	users := Mail.GetNewMail(cookie)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		rand.Seed(time.Now().UnixNano())
		picNum := strconv.Itoa(rand.Intn(18) + 1)
		users.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", Text.SelectFirst10WithPicture(picNum), gomail.NewMessage(), ".\\pic\\"+picNum+".png")
	}()

	ctx.String(http.StatusOK, "已发送，如果没有收到请检查垃圾箱。")
	wg.Wait()
}
