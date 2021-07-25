package main

import (
	"project/Text"
	"project/functions"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	go countTime()

	router.LoadHTMLGlob("HTML/*")
	router.GET("/", functions.ToLogin)
	router.GET("/toRegister", functions.ToRegister)
	router.GET("/toChangePassword", functions.ToChangePassword)
	router.GET("/toFindPassword", functions.FindPassword)
	router.GET("/toFunctions", functions.ToFuntions)
	router.GET("/toverificationCode", functions.ToVerificationCode)

	router.POST("/register", functions.Register)
	router.POST("/login", functions.Login)
	router.POST("/changePassword", functions.ChangePassWord)
	router.POST("/findPassword", functions.FindPassword)
	router.POST("/sendCode", functions.SendCode)
	router.GET("/sendStock", functions.SendStock)

	router.Run()
}

func countTime() {
	Text.GenerateText()
	for {
		// 两小时更新一次
		<-time.After(time.Hour * 2)
		Text.GenerateText()
	}
}
