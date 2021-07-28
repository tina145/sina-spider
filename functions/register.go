package functions

import (
	"net/http"
	"project/Users"

	"github.com/gin-gonic/gin"
)

func ToRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", nil)
}

func ToVerificationCode(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "verificationCode.html", nil)
}

func Register(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	passWord := ctx.PostForm("passWord")

	if userName == "" {
		ctx.String(http.StatusOK, "邮箱不能为空")
		return
	} else if passWord == "" {
		ctx.String(http.StatusOK, "密码不能为空")
		return
	}

	userInfo := &Users.User{
		MailAccount:  userName,
		MailPassword: passWord,
	}

	if userInfo.CheckUserExist() {
		ctx.String(http.StatusOK, "用户已存在")
		return
	}

	code := ctx.PostForm("code")

	if code != userInfo.GetVerificationCode() {
		ctx.String(http.StatusOK, "验证码错误！")
		return
	}

	userInfo.Register()

	// domain 是域名，path 是域名，合起来限制可以被哪些 url 访问
	ctx.SetCookie("cookie", userName, 86400, "/", "localhost:8080", false, true)
	ctx.HTML(http.StatusOK, "function.html", nil)
}

func SendCode(ctx *gin.Context) {
	userName := ctx.PostForm("userName")
	if userName == "" {
		ctx.String(http.StatusOK, "邮箱不能为空")
		return
	}
	userInfo := &Users.User{
		MailAccount: userName,
	}

	if userInfo.CheckUserExist() {
		ctx.String(http.StatusBadRequest, "用户已存在")
		return
	}

	err := userInfo.Verification()
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "register.html", nil)
}
