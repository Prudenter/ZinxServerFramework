package zinxInterface

/*
	定义Server模块的抽象层
*/
type ZinxInterfaceServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Run()

	//提供一个得到链接管理模块的方法
	GetConnMgr() InterfaceConnManager

	//定义添加路由的方法,暴露给开发者
	AddRouter(messageId uint32, router InterfaceRouter)

	//注册 创建链接之后调用的Hook函数的方法
	AddOnConnStart(hookFunc func(conn InterfaceConnection))

	//注册 销毁链接之前调用的Hook函数 的方法
	AddOnConnStop(hookFunc func(conn InterfaceConnection))

	//调用 创建链接之后的HOOK函数的方法
	CallOnConnStart(conn InterfaceConnection)

	//调用 销毁链接之前调用的HOOk函数的方法
	CallOnConnStop(conn InterfaceConnection)
}
