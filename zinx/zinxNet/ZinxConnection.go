package zinxNet

import (
	"net"
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
	"ZinxServerFramework/zinx/utils"
)

//实现具体的TCP链接模块
type ZinxConnection struct {
	zinxInterface.InterfaceConnection

	//当前链接的原生套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前的链接状态
	IsClosed bool

	//当前链接所绑定的Router
	Router zinxInterface.InterfaceRouter
}

/*
	初始化链接方法
*/
func NewZinxConnection(conn *net.TCPConn, connId uint32, router zinxInterface.InterfaceRouter) zinxInterface.InterfaceConnection {
	return &ZinxConnection{
		Conn:     conn,
		ConnID:   connId,
		IsClosed: false,
		Router:   router,
	}
}

//针对链接读业务的方法
func (conn *ZinxConnection) StartReader() {
	//从对端读数据
	fmt.Println("Reader go is startin...")
	defer fmt.Println("connId = ", conn.ConnID, "Reader is exit,remote addr is =", conn.GetRemoteAddr().String())
	defer conn.Stop()
	for {
		buf := make([]byte, utils.Globj.MaxPackage)
		n, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buff err:", err)
			break
		}
		//将当前一次性得到的客户端请求数据封装成一个Request
		request := NewZinxRequest(conn, buf, n)

		//调用用户传递进来的业务处理方法,即自定义router中的业务处理方法--模板设计模式
		go func() { //添加go程,防止阻塞
			//在此处理核心业务
			conn.Router.PreHandle(request)
			conn.Router.Handle(request)
			conn.Router.PostHandle(request)
		}()
	}
}

//实现抽象接口的方法

//启动链接
func (conn *ZinxConnection) Start() {
	fmt.Println("Conn start()...id=", conn.ConnID)
	//先进行读业务,添加go程,将读写进行分离
	go conn.StartReader()

	//TODO 进行写业务
}

//停止链接
func (conn *ZinxConnection) Stop() {
	fmt.Println("c.stop()...ConnId=", conn.ConnID)
	//回收工作
	if conn.IsClosed == true {
		return
	}
	conn.IsClosed = true
	//关闭原生的套接字
	conn.Conn.Close()
}

//获取链接ID
func (conn *ZinxConnection) GetConnID() uint32 {
	return conn.ConnID
}

//获取conn的原生socket套接字
func (conn *ZinxConnection) GetTCPConnection() *net.TCPConn {
	return conn.Conn
}

//获取远程客户端的ip地址
func (conn *ZinxConnection) GetRemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

//发送数据给对方客户端
func (conn *ZinxConnection) Send(data []byte, cnt int) error {
	if _, err := conn.Conn.Write(data[:cnt]); err != nil {
		fmt.Println("send buf err:", err)
	}
	return nil
}
