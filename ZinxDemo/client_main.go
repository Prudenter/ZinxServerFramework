package main

import (
	"fmt"
	"time"
	"net"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("client start...")
	//等待一下,给个缓冲
	time.Sleep(1*time.Second)
	//直接connect服务器得到一个已经建立好的conn句柄
	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("Dial err :",err)
		return
	}

	for{
		//写数据到服务器
		_,err := conn.Write([]byte("Hello ZinxServer.."))
		if err != nil {
			fmt.Println("Write err:",err)
			return
		}
		//读取服务器回显的数据
		buf := make([]byte,512)
		n,err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read err:",err)
			return
		}
		fmt.Printf("service call back :%s,n=%d\n",buf,n)
		//等待一下,给个缓冲
		time.Sleep(1*time.Second)
	}
}
