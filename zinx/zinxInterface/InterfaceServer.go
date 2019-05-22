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
}

