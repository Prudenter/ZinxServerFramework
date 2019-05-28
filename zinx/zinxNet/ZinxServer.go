package zinxNet

import (
	"fmt"
	"net"
	"ZinxServerFramework/zinx/zinxInterface"
	"ZinxServerFramework/zinx/utils"
)

/*
	定义Server模块的实现层
 */
type ZinxServer struct {
	zinxInterface.ZinxInterfaceServer

	//服务器ip
	IP string

	//ip版本
	IPVersion string

	//服务器端口
	Port int

	//服务器名称
	Name string

	//多路由的消息管理模块
	MsgHandler zinxInterface.InterfaceMsgHandler
}

//定义初始化服务器的方法
func NewServer() zinxInterface.ZinxInterfaceServer {
	server := &ZinxServer{
		IP:         utils.Globj.Host,
		IPVersion:  "tcp4",
		Port:       utils.Globj.Port,
		Name:       utils.Globj.Name,
		MsgHandler: NewZinxMsgHandler(),
	}
	return server
}

//实现抽象接口的方法

//启动服务器,实现服务器监听---使用原生socket 服务器编程
func (server *ZinxServer) Start() {
	fmt.Printf("[start] %s Linstenner at IP:%s,Port:%d,is starting..\n", server.Name, server.IP, server.Port)

	//1.在监听之前启动工作池
	server.MsgHandler.StartWorkerPool()

	//2.创建套接字,得到一个TCP的addr
	addr, err := net.ResolveTCPAddr(server.IPVersion, fmt.Sprintf("%s:%d", server.IP, server.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
		return
	}

	//3.监听服务器地址
	listenner, err := net.ListenTCP(server.IPVersion, addr)
	if err != nil {
		fmt.Println("ListenTCP err:", err)
		return
	}

	//生成connid的累加器
	var connId uint32
	connId = 0

	//4.阻塞等待客户端发送请求
	go func() { //如果不加go程,Start()会一直阻塞,则主go程也会阻塞,无法执行主go程的其他扩展
		for {
			//阻塞等待客户端发送请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("ListenTCP err:", err)
				continue
			}

			//创建一个ZinxConnection对象,并传入当前消息的管理模块
			dealConn := NewZinxConnection(conn, connId, server.MsgHandler)
			connId ++
			//启动链接,进行业务处理
			go dealConn.Start() // 如果不加go程,则当前子go程会一直阻塞,无法进行并发访问,不能同时处理多个客户端的请求
		}
	}()
}

//停止服务器
func (server *ZinxServer) Stop() {
	//TODO 将一些服务器资源进行回收
}

//运行服务器
func (server *ZinxServer) Run() {
	//启动server的监听功能
	server.Start()

	//TODO 做一些其他的扩展

	//main函数阻塞在这
	select { //保证main函数不退出

	}
}

//添加路由的方法
func (server *ZinxServer) AddRouter(messageId uint32, router zinxInterface.InterfaceRouter) {
	//添加路由和messageid到msghandler中
	server.MsgHandler.AddRouter(messageId, router)
	fmt.Println("Add Router success! messageId = ", messageId)
}
