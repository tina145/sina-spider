package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

var targetIP string = "想连接的服务器 IP 及端口"

func main() {
	logFile, err := os.OpenFile(`.\\Logs.txt`, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	defer func() {
		logFile.Close()
	}()
	if err != nil {
		fmt.Println("文件错误")
		time.Sleep(time.Hour)
		return
	}

	// 设置输出文件
	log.SetOutput(logFile)

	// 使用时要换成自己的 IP 及端口号
	fmt.Println("尝试连接中...")
	connect, err := net.Dial("tcp", targetIP)

	log.SetPrefix(connect.LocalAddr().String() + "：")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
		time.Sleep(time.Hour)
		return
	}
	defer connect.Close()

	fmt.Println("连接成功!")
	fmt.Println("请选择执行的功能：")
	fmt.Println("\t\t输入 register 注册账号（使用邮箱）")
	fmt.Println("\t\t输入 login 登录")
	fmt.Println("\t\t输入 send 获取一次要闻")
	fmt.Println("\t\t输入 exit 退出")
	fmt.Println("\t\t输入 change password 修改密码")
	fmt.Println("\t\t输入 logout 注销")

	defer connect.Write([]byte("exit"))
	wg.Add(1)
	go send(connect)
	wg.Add(1)
	go receive(connect)

	wg.Wait()
	time.Sleep(time.Second * 3)
}

// 发送请求
func send(connect net.Conn) {
	defer wg.Done()
	obj := bufio.NewScanner(os.Stdin)

	for {
		obj.Scan()
		text := obj.Text()
		if len(text) > 128 {
			fmt.Println("字符太多，请重新输入")
			continue
		}

		err := connect.SetWriteDeadline(time.Now().Add(time.Second * 60))

		if err != nil {
			log.Println(err)
			fmt.Println(err)
			return
		}

		_, err = connect.Write([]byte(text))
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			return
		}

		if text == "close" || text == "break" || text == "exit" {
			break
		}
	}
}

// 接收数据
func receive(connect net.Conn) {
	defer wg.Done()
	for {
		text := make([]byte, 1024)

		err := connect.SetReadDeadline(time.Now().Add(time.Second * 60))
		if err != nil {
			log.Println(err)
			fmt.Println(err)
			return
		}
		n, err := connect.Read(text)

		if err != nil {
			fmt.Println(err)
			log.Println(err)
			return
		}
		text = text[:n]

		if string(text) == "Welcome to use next time!" || string(text) == "server was offline!" {
			fmt.Println(string(text))
			log.Println(string(text))
			break
		} else {
			fmt.Println(string(text))
			log.Println(string(text))
		}
	}
}
