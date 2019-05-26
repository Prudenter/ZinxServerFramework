package zinxInterface

/*
	定义抽象的消息管理模块,存放router集合的
*/

type InterfaceMsgHandler interface {
	//添加路由到map集合中
	AddRouter(messageId uint32,router InterfaceRouter)
	//根据messageId 调度路由
	DoMsgHandler(request InterfaceRequest)
}
