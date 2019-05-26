package zinxInterface

/*
	抽象InterfaceRequest,对一次性请求的数据进行封装
*/

type InterfaceRequest interface {
	//得到当前请求的链接
	GetConnection() InterfaceConnection

	//得到客户端发送的消息
	GetMessage() InterfaceMessage
}