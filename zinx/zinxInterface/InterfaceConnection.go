package zinxInterface

import "net"

/*
	定义抽象的Tcp链接层,,用于将链接请求和用户的业务进行绑定
*/
type InterfaceConnection interface {
	//启动链接
	Start()

	//停止链接
	Stop()

	//获取链接ID
	GetConnID() uint32

	//获取conn的原声socket套接字
	GetTCPConnection() *net.TCPConn

	//获取远程客户端的ip地址
	GetRemoteAddr() net.Addr

	//发送数据给对方客户端
	Send(messageId uint32,messageData []byte)error

	//设置属性
	SetProperty(key string,value interface{})

	//获取属性
	GetProperty(key string)(interface{},error)

	//删除属性
	RemoveProperty(key string)
}

//定义抽象的业务处理方法,将函数指针定义在抽象层,符合依赖倒转设计原则
type HandleFunc func(request InterfaceRequest) error