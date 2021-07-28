package functions

import (
	"fmt"
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

func SendStock(ctx *gin.Context) {
	cookie, err := ctx.Cookie("cookie")
	if err != nil {
		ctx.String(http.StatusBadRequest, "请先登录")
		return
	}
	users := Mail.GetNewMail(cookie)
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	picNum := strconv.Itoa(rand.Intn(18) + 1)
	err = users.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", Text.SelectFirst10WithPicture(picNum), gomail.NewMessage(), "pic/"+picNum+".png")
	if err != nil {
		fmt.Println(err)
	}

	ctx.String(http.StatusOK, "已发送，如果没有收到请检查垃圾箱。")
	wg.Wait()
}
