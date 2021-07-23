package functions

import (
	"1/Users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindPassword(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	userInfo := &Users.User{
		MailAccount: userName,
	}

	// 发送验证码
	userInfo.Verification()

	code := ctx.PostForm("code")
	if code != userInfo.GetVerificationCode() {
		ctx.String(http.StatusOK, "验证码错误！")
		return
	}

	newPassword := ctx.PostForm("newPassword")

	userInfo.ChangePassword(newPassword)

	ctx.String(http.StatusOK, "修改完成！")
}
