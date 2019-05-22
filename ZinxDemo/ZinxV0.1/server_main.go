package main

import "ZinxServerFramework/zinx/zinxNet"

/**
	定义基于zinx框架的服务器
*/
func main() {
	//创建一个zinx server对象
	s := zinxNet.Init("zinxV0.1 ")
	//让server对象启动服务
	s.Run()
}
