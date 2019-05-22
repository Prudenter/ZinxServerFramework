package zinxNet

import (
	"net"
	"ZinxServerFramework/zinx/zinxInterface"
)

//实现具体的TCP链接模块
type Connection struct {

	zinxInterface.InterfaceConnection

	//当前链接的原声套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前的链接状态
	IsClosed bool

	//当前链接所绑定的业务处理方法
	HandleAPI zinxInterface.HandleFunc
}

/*
	初始化链接方法
*/
func NewConnection(conn *net.TCPConn, connId uint32, callBack_api zinxInterface.HandleFunc) zinxInterface.InterfaceConnection {
	return &Connection{
		Conn:conn,
		ConnID:connId,
		HandleAPI:callBack_api,
		IsClosed:false,
	}
}

//启动链接
func (conn *Connection)Start(){

}

//停止链接
func (conn *Connection)Stop(){

}

//获取链接ID
func (conn *Connection)GetConnID() uint32{
	return 0
}

//获取conn的原声socket套接字
func (conn *Connection)GetTCPConnection() *net.TCPConn{
	return nil
}

//获取远程客户端的ip地址
func (conn *Connection)GetRemoteAddr() *net.Addr{
	return nil
}

//发送数据给对方客户端
func (conn *Connection)Send(data[]byte)error{
	return nil
}