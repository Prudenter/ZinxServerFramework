package zinxNet

import (
	"fmt"
	"net"
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

//定义一个具体的回显业务,针对 type HandleFunc func(*net.TCPConn,[]byte,int) error
func CallBackFunc(conn *net.TCPConn, data []byte, cnt int) error {
	//处理回显业务
	fmt.Println("[conn Handle] CallBack...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back err:", err)
		return err
	}
	return nil
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

	//生成connid的累加器
	var connId uint32
	connId = 0

	//3.阻塞等待客户端发送请求
	go func() { //如果不加go程,Start()会一直阻塞,则主go程也会阻塞,无法执行主go程的其他扩展
		for {
			//阻塞等待客户端发送请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("ListenTCP err:", err)
				continue
			}

			//创建一个Connection对象
			dealConn := NewConnection(conn, connId, CallBackFunc)
			connId ++
			//启动链接,进行业务处理
			go dealConn.Start()    	// 如果不加go程,则当前子go程会一直阻塞,无法进行并发访问,不能同时处理多个客户端的请求
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
	select { //保证main函数不退出

	}
}
