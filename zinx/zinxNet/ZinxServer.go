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

	//链接管理模块
	ConnMgr zinxInterface.InterfaceConnManager

	//该server创建链接之后自动调用的Hook函数
	OnConnStart func(conn zinxInterface.InterfaceConnection)

	//该server销毁链接之前自动调用的Hook函数
	OnConnStop func(conn zinxInterface.InterfaceConnection)
}

//定义初始化服务器的方法
func NewServer() zinxInterface.ZinxInterfaceServer {
	server := &ZinxServer{
		IP:         utils.Globj.Host,
		IPVersion:  "tcp4",
		Port:       utils.Globj.Port,
		Name:       utils.Globj.Name,
		MsgHandler: NewZinxMsgHandler(),
		ConnMgr:    NewZinxConnManager(),
	}
	return server
}

//实现抽象接口的方法

//启动服务器,实现服务器监听---使用原生socket 服务器编程
func (zs *ZinxServer) Start() {
	fmt.Printf("[start] %s Linstenner at IP:%s,Port:%d,is starting..\n", zs.Name, zs.IP, zs.Port)

	//1.在监听之前启动工作池
	zs.MsgHandler.StartWorkerPool()

	//2.创建套接字,得到一个TCP的addr
	addr, err := net.ResolveTCPAddr(zs.IPVersion, fmt.Sprintf("%s:%d", zs.IP, zs.Port))
	if err != nil {
		fmt.Println("ResolveTCPAddr err:", err)
		return
	}

	//3.监听服务器地址
	listenner, err := net.ListenTCP(zs.IPVersion, addr)
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

			//判断当前的server链接数量是否已经到最大值,如果达到最大值,则不在添加链接---先判断再添加
			if zs.ConnMgr.GetLen() >= int(utils.Globj.MaxConn) {
				//当前链接已经满了
				fmt.Println("--->Too many Connection,Maxconn = ", utils.Globj.MaxConn)
				conn.Close()
				continue
			}

			//创建一个ZinxConnection对象,并传入当前消息的管理模块
			dealConn := NewZinxConnection(conn, connId, zs.MsgHandler, zs)

			connId ++
			//启动链接,进行业务处理
			go dealConn.Start() // 如果不加go程,则当前子go程会一直阻塞,无法进行并发访问,不能同时处理多个客户端的请求
		}
	}()
}

//停止服务器
func (zs *ZinxServer) Stop() {
	//服务器停止时,清空当前全部的链接
	zs.ConnMgr.ClearConn()
}

//运行服务器
func (zs *ZinxServer) Run() {
	//启动server的监听功能
	zs.Start()

	//TODO 做一些其他的扩展

	//main函数阻塞在这
	select { //保证main函数不退出

	}
}

//添加路由的方法
func (zs *ZinxServer) AddRouter(messageId uint32, router zinxInterface.InterfaceRouter) {
	//添加路由和messageid到msghandler中
	zs.MsgHandler.AddRouter(messageId, router)
	fmt.Println("Add Router success! messageId = ", messageId)
}

//得到链接管理模块的方法
func (zs *ZinxServer) GetConnMgr() zinxInterface.InterfaceConnManager {
	return zs.ConnMgr
}

//注册 创建链接之后调用的Hook函数的方法
func (zs *ZinxServer) AddOnConnStart(hookFunc func(conn zinxInterface.InterfaceConnection)){
	zs.OnConnStart = hookFunc
}

//注册 销毁链接之前调用的Hook函数 的方法
func (zs *ZinxServer)AddOnConnStop(hookFunc func(conn zinxInterface.InterfaceConnection)){
	zs.OnConnStop = hookFunc
}

//调用 创建链接之后的HOOK函数的方法
func (zs *ZinxServer)CallOnConnStart(conn zinxInterface.InterfaceConnection){
	if zs.OnConnStart != nil {
		fmt.Println("--->Call OnConnStart()...")
		zs.OnConnStart(conn)
	}
}

//调用 销毁链接之前调用的HOOk函数的方法
func (zs *ZinxServer)CallOnConnStop(conn zinxInterface.InterfaceConnection){
	if zs.OnConnStop != nil {
		fmt.Println("--->Call OnConnStop()...")
		zs.OnConnStop(conn)
	}
}