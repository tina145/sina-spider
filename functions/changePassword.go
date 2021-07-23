package functions

import (
	"1/Users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChangePassWord(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	oldPassword := ctx.PostForm("oldPassword")
	newPassword := ctx.PostForm("newPassword")
	newPasswordConfirm := ctx.PostForm("newPasswordConfirm")

	if newPasswordConfirm != newPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致！")
		return
	}

	userInfo := &Users.User{
		MailAccount:  userName,
		MailPassword: oldPassword,
	}
	userInfo.ChangePassword(newPassword)
	ctx.String(http.StatusOK, "修改完成！")
}
