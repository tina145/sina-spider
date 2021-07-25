package functions

import (
	"net/http"
	"project/Mail"
	"project/Text"
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
	err = users.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", Text.SelectFirst20(), gomail.NewMessage())
	if err != nil {
		ctx.String(http.StatusBadRequest, "发送失败")
		return
	}
	ctx.String(http.StatusOK, "已发送，如果没有收到请检查垃圾箱。")
}
