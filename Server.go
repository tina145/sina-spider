package main

import (
	"GoProject/spider/serverExe"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// 日志文件
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
	// 服务器输出的日志
	log.SetPrefix("Server：")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 服务器开始监听
	server, err := net.Listen("tcp", ":8888")
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

		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(connect.RemoteAddr().String() + " 已建立连接")
		log.Println(connect.RemoteAddr().String() + " 已建立连接")
		// 新开一个协程处理请求
		go serverExe.Solve(connect)
	}
}
