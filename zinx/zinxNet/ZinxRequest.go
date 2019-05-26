package zinxNet

import "ZinxServerFramework/zinx/zinxInterface"

/*
	定义InterfaceRequest模块的实现层
*/
type ZinxRequest struct {
	zinxInterface.InterfaceRequest
	//链接信息
	conn zinxInterface.InterfaceConnection
	//客户端发送的消息
	message zinxInterface.InterfaceMessage
}

//定义初始化请求的方法
func NewZinxRequest(conn zinxInterface.InterfaceConnection,message zinxInterface.InterfaceMessage )zinxInterface.InterfaceRequest  {
	return &ZinxRequest{
		conn:conn,
		message:message,
	}
}

//实现抽象接口的方法

//得到当前请求的链接
func (zrqst *ZinxRequest)GetConnection() zinxInterface.InterfaceConnection{
	return  zrqst.conn
}

//得到客户端发送的消息
func (zrqst *ZinxRequest)GetMessage() zinxInterface.InterfaceMessage{
	return zrqst.message
}
