package functions

import (
	"1/Users"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "file", nil)

	userName := ctx.PostForm("userName")
	passWord := ctx.PostForm("passWord")

	userInfo := &Users.User{
		MailAccount:  userName,
		MailPassword: passWord,
	}

	statu := userInfo.Login()

	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	if statu == "success" {
		conn.Do("set", userInfo.MailAccount, "login")
	}

	ctx.String(http.StatusOK, statu)
}
