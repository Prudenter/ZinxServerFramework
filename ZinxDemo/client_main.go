package main

import (
	"fmt"
	"time"
	"net"
	"ZinxServerFramework/zinx/zinxNet"
	"io"
)

/*
	模拟客户端
*/
func main() {
	fmt.Println("client start...")
	//等待一下,给个缓冲
	time.Sleep(1 * time.Second)
	//直接connect服务器得到一个已经建立好的conn句柄
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Dial err :", err)
		return
	}

	for {
		//写数据到服务器
		//先封包再发包,定义封包拆包对象
		zdp := zinxNet.NewZinxDataPack()
		//把数据封包成二进制数据流形式
		binaryMessage, err := zdp.Pack(zinxNet.NewZinxMessage(0, []byte("this is client test message")))
		if err != nil {
			fmt.Println("Pack err :", err)
			break
		}
		if _, err := conn.Write(binaryMessage); err != nil {
			fmt.Println("Write err :", err)
			break
		}
		//获取服务器回发的数据,进行拆包
		binaryHead := make([]byte, zdp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Clien ReadFull binaryHead err :", err)
			break
		}
		//将头部数据进行拆包
		message, err := zdp.UnPack(binaryHead)
		//判断数据是否为空
		if message.GetMsgLen() > 0 {
			//读取包体
			msg := message.(*zinxNet.ZinxMessage)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("Clien ReadFull msg.Data err :", err)
				break
			}
			fmt.Println("--->Receive Server message:id=", msg.Id, "len=", msg.GetMsgLen(), "data=", string(msg.GetMsgData()))
		}
		//等待一下,给个缓冲
		time.Sleep(1 * time.Second)
	}
}
