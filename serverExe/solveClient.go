package serverExe

import (
	"GoProject/spider/httpRequest"
	"GoProject/spider/spiderText"
	"GoProject/spider/spiderUsers"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-gomail/gomail"
)

func Solve(connect net.Conn) {
	isLogin := false

	user := &spiderUsers.User{}

	defer func() {
		connect.Write([]byte("Welcome to use next time!"))
		fmt.Println(connect.RemoteAddr().String() + " 已断开连接")
		log.Println(connect.RemoteAddr().String() + " 已断开连接")
	}()

	for {
		temp, err := ReadText(connect)
		if err != nil {
			break
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

			spiderText.GenerateText()
			// 正文内容
			texts := spiderText.SelectFirst20()

			connect.Write([]byte("delivering..."))

			// 构造邮件对象
			mail := httpRequest.GetNewMail(user.MailAccount)

			// 发送消息
			mail.Send(time.Now().String()[:19]+" "+time.Now().Weekday().String()+"：每日要闻", texts, gomail.NewMessage())

			connect.Write([]byte("complete!"))

			fmt.Println(connect.RemoteAddr().String() + " 发送完成")
			log.Println(connect.RemoteAddr().String() + " 发送完成")
		} else if string(text) == "register" { // 注册功能
			connect.Write([]byte("please input the account:"))
			temp, err := ReadText(connect)
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

				temp, err = ReadText(connect)
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
				temp, err = ReadText(connect)
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
			temp, err = ReadText(connect)
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

				temp, err := ReadText(connect)
				if err != nil {
					return
				}
				user.MailAccount = string(temp)
				connect.Write([]byte("please input the password:"))

				temp, err = ReadText(connect)
				if err != nil {
					return
				}
				user.MailPassword = string(temp)

				status := user.Login()
				if status == "success" {
					isLogin = true
				}
				connect.Write([]byte(status))
			}
		} else if string(text) == "change password" {
			if !isLogin {
				connect.Write([]byte("please sign in"))
				continue
			}
			connect.Write([]byte("please input new password"))

			temp, err := ReadText(connect)
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

func ReadText(connect net.Conn) ([]byte, error) {
	text := make([]byte, 1024)

	connect.SetDeadline(time.Now().Add(time.Second * 60))
	n, err := connect.Read(text)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return nil, err
	}
	return text[:n], nil
}
