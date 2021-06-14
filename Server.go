package main

import (
	"GoProject/spider/tools"
	"GoProject/spider/users"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/gomodule/redigo/redis"
)

func main() {
	// 日志文件
	logFile, err := os.OpenFile(`.\\Logs.txt`, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	defer logFile.Close()
	if err != nil {
		fmt.Println("文件错误")
		time.Sleep(time.Hour)
		return
	}
	// 设置输出文件
	log.SetOutput(logFile)
	// 服务器输出的日志
	log.SetPrefix("Server：")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 服务器开始监听
	server, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
		time.Sleep(time.Hour)
		return
	}

	// 服务器关闭时执行的操作
	defer server.Close()

	fmt.Println("服务器监听中...")
	log.Println("服务器监听中...")

	for {
		// 循环接收连接请求
		connect, err := server.Accept()
		defer connect.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(connect.RemoteAddr().String() + " 已建立连接")
		log.Println(connect.RemoteAddr().String() + " 已建立连接")
		// 新开一个协程处理请求
		go solve(connect)
	}
}

func solve(connect net.Conn) {
	// 连接 redis
	redisConn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	isLogin := false

	user := &users.User{}

	defer func() {
		connect.Write([]byte("Welcome to use next time!"))
		fmt.Println(connect.RemoteAddr().String() + " 已断开连接")
		log.Println(connect.RemoteAddr().String() + " 已断开连接")
		redisConn.Close()
	}()

	for {
		temp, err := readText(connect)
		if err != nil {
			return
		}
		text := temp

		if string(text) == "close" || string(text) == "exit" || string(text) == "break" {
			break
		} else if string(text) == "send" { // 发送要闻
			if !isLogin {
				connect.Write([]byte("please sign in"))
				continue
			}
			fmt.Println("正在向 " + connect.RemoteAddr().String() + "发送中...")
			log.Println(connect.RemoteAddr().String() + " 请求 " + string(text))

			tools.GenerateText(tools.Search())
			// 正文内容
			texts := tools.SelectFirst20()

			connect.Write([]byte("delivering..."))

			// 构造邮件对象
			mail := &tools.Mail{
				// 根据自己需求修改
				SenderAccount:  xxx@xxx,
				SenderPassword: yyy,
				Receiver:       user.MailAccount,
				ServerAddr:     "smtp.office365.com",
				ServerPort:     587,
			}

			// 发送消息
			mail.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", texts, gomail.NewMessage())

			connect.Write([]byte("complete!"))

			fmt.Println(connect.RemoteAddr().String() + " 发送完成")
			log.Println(connect.RemoteAddr().String() + " 发送完成")
		} else if string(text) == "register" { // 注册功能
			connect.Write([]byte("please input the account:"))
			temp, err := readText(connect)
			if err != nil {
				connect.Write([]byte("error"))
				return
			}
			registerAccount := temp

			if user.CheckUserExist(string(registerAccount)) {
				connect.Write([]byte("already exist user"))
				continue
			}

			// 校验邮箱
			for {
				connect.Write([]byte("The verification code has been sent to your email, please check and enter the verification code."))

				connect.Write([]byte("input verify to resend the verification code."))

				user.Verification(string(registerAccount))

				temp, err = readText(connect)
				if err != nil {
					return
				}

				verificationCode := string(temp)

				if verificationCode == "exit" {
					return
				} else if verificationCode == "resend" {
					user.Verification(string(registerAccount))
					continue
				}

				code := user.GetVerificationCode(string(registerAccount))
				if code == "not exist" {
					connect.Write([]byte("the verification code has expired."))
				} else if code != verificationCode {
					connect.Write([]byte("verification code wrong"))
				} else {
					break
				}

				connect.Write([]byte("please input the account:"))
				temp, err = readText(connect)
				if err != nil {
					connect.Write([]byte("error"))
					return
				}
				registerAccount = temp

				if user.CheckUserExist(string(registerAccount)) {
					connect.Write([]byte("already exist user"))
					break
				}
			}

			connect.Write([]byte("please input the password:"))
			temp, err = readText(connect)
			if err != nil {
				return
			}
			registerPassword := temp

			status := user.Register(string(registerAccount), string(registerPassword))

			connect.Write([]byte(status))
		} else if string(text) == "login" { // 登录功能
			if isLogin {
				connect.Write([]byte("already login"))
			} else {
				connect.Write([]byte("please input the account:"))

				temp, err := readText(connect)
				if err != nil {
					return
				}
				user.MailAccount = string(temp)
				connect.Write([]byte("please input the password:"))

				temp, err = readText(connect)
				if err != nil {
					return
				}
				user.MailPassword = string(temp)

				status := user.Login()
				if status == "success" {
					isLogin = true
					redisConn.Do("set", connect.RemoteAddr().String(), "login", "ex", 300)
				}
				connect.Write([]byte(status))
			}
		} else if string(text) == "change password" {
			if !isLogin {
				connect.Write([]byte("please sign in"))
				continue
			}
			connect.Write([]byte("please input new password"))

			temp, err := readText(connect)
			if err != nil {
				log.Println(err)
				return
			}

			newpassword := string(temp)
			connect.Write([]byte(user.ChangePassword(newpassword)))
		} else if string(text) == "logout" {
			isLogin = false
			connect.Write([]byte("success"))
		} else {
			connect.Write([]byte("the function is under development..."))
		}
	}
}

func readText(connect net.Conn) ([]byte, error) {
	text := make([]byte, 1024)

	// 60 秒超时
	connect.SetDeadline(time.Now().Add(time.Second * 60))
	n, err := connect.Read(text)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return nil, err
	}
	return text[:n], nil
}
