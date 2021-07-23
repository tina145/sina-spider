package functions

import (
	"1/Text"
	"1/httpRequest"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

var wg sync.WaitGroup

func Send(ctx *gin.Context) {
	userEmail := ctx.PostForm("un")
	Text.GenerateText()

	// 正文内容
	texts := Text.SelectFirst20()
	mail := httpRequest.GetNewMail(userEmail)

	// 发送消息
	go func() {
		wg.Add(1)
		mail.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", texts, gomail.NewMessage())
		wg.Done()
	}()

	ctx.String(http.StatusOK, "已发送，如果没有收到请检查垃圾箱。")
	wg.Wait()
}
