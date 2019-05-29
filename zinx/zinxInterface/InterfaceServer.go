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
}
