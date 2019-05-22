package zinxNet

import (
	"fmt"
	"net"
	"io"
	"ZinxServerFramework/zinx/zinxInterface"
)

/*
	定义Server模块的实现层
 */
type Server struct {
	zinxInterface.ZinxInterfaceServer

	//服务器ip
	IP string
	//ip版本
	IPVersion string
	//服务器端口
	Port int
	//服务器名称
	Name string
}

//定义初始化服务器的方法
func Init(name string) zinxInterface.ZinxInterfaceServer {
	server := &Server{
		IP:        "0.0.0.0",
		IPVersion: "tcp4",
		Port:      7777,
		Name:      name,
	}
	return server
}

//实现抽象接口的方法
//启动服务器,实现服务器监听---使用原生socket 服务器编程
func (server *Server) Start() {
	fmt.Printf("[start] Server Linstenner at IP:%s,Port:%d,is starting..\n", server.IP, server.Port)
	//1.创建套接字,得到一个TCP的addr
	addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
		return
	}

	//2.监听服务器地址
	listenner, err := net.ListenTCP(server.IPVersion, addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
		return
	}

	//3.阻塞等待客户端发送请求
	go func() {//如果不加go程,则主go程会一直阻塞,无法执行主go程的其他扩展
		for {
			//阻塞等待客户端发送请求
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("ListenTCP err:", err)
				continue
			}
			//此时已经和客户端建立了连接
			go func() {//如果不加go程,则当前子go程会一直阻塞,无法进行并发访问,不能同时处理多个客户端的请求
				//4.客户端有数据请求，循环处理客户端业务,包括读和写
				for {
					buf := make([]byte,512)
					n,err := conn.Read(buf)
					if err!=nil && err==io.EOF{
						fmt.Println("请求数据读完了")
						break
					}
					fmt.Printf("receive client buf %s,n = %d\n",buf,n)

					//回显功能,进行业务处理
					_,err = conn.Write(buf[:n])
					if err!=nil {
						fmt.Println("Write err:",err)
						continue
					}
				}
			}()
		}
	}()
}

//停止服务器
func (server *Server) Stop() {
	//TODO 将一些服务器资源进行回收
}

//运行服务器
func (server *Server) Run() {
	//启动server的监听功能
	server.Start()

	//TODO 做一些其他的扩展
	//阻塞 //告诉CPU不再需要处理，节省cpu资源
	select{ //保证main函数不退出

	}
}
