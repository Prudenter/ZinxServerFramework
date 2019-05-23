package zinxNet

import "ZinxServerFramework/zinx/zinxInterface"

/*
	定义InterfaceRequest模块的实现层
*/
type ZinxRequest struct {
	zinxInterface.InterfaceRequest
	//链接信息
	conn zinxInterface.InterfaceConnection
	//数据内容
	data []byte
	//数据长度
	length int
}

//定义初始化请求的方法
func NewZinxRequest(conn zinxInterface.InterfaceConnection,data []byte,length int)zinxInterface.InterfaceRequest  {
	return &ZinxRequest{
		conn:conn,
		data:data,
		length:length,
	}
}

//实现抽象接口的方法

//得到当前请求的链接
func (zrqst *ZinxRequest)GetConnection() zinxInterface.InterfaceConnection{
	return  zrqst.conn
}

//得到链接的数据
func (zrqst *ZinxRequest)GetData() []byte{
	return zrqst.data
}

//得到数据的长度
func (zrqst *ZinxRequest)GetDataLen() int{
	return zrqst.length
}