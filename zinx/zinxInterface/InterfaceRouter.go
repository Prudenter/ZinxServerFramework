package zinxInterface

/*
	定义抽象的路由模块
*/

type InterfaceRouter interface {
	//处理业务之前要调用的方法
	PreHandle(requset InterfaceRequest)
	//真正处理业务要调用的方法
	Handle(requset InterfaceRequest)
	//处理业务之后要调用的方法
	PostHandle(requset InterfaceRequest)
}
