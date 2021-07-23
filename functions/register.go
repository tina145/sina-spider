package functions

import (
	"1/Users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Regiter(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	passWord := ctx.PostForm("passWord")

	userInfo := &Users.User{
		MailAccount:  userName,
		MailPassword: passWord,
	}

	if userInfo.CheckUserExist() {
		ctx.String(http.StatusOK, "用户已存在")
		return
	}

	userInfo.Verification()

	code := ctx.PostForm("code")

	if code != userInfo.GetVerificationCode() {
		ctx.String(http.StatusOK, "验证码错误！")
		return
	}

	ctx.String(http.StatusOK, userInfo.Register())
}
