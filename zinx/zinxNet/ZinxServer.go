package zinxNet

import (
	"fmt"
	"ZinxServerFramework/zinx/zinxInterface"
)

/*
	定义Server模块的实现层
 */
type Server struct {
	zinxInterface.ZinxInterfaceServer

	//服务器ip
	IP string
	//ip版本
	IPVersion string
	//服务器端口
	Port int
	//服务器名称
	Name string
}

//定义初始化服务器的方法
func Init(name string) zinxInterface.ZinxInterfaceServer {
	server := &Server{
		IP:        "0.0.0.0",
		IPVersion: "tcp4",
		Port:      7777,
		Name:      name,
	}
	return server
}

//实现抽象接口的方法
//启动服务器
func (server *Server) Start() {
	fmt.Println("服务器已启动")
}

//停止服务器
func (server *Server) Stop() {
	//TODO 将一些服务器资源进行回收
}

//运行服务器
func (server *Server) Run() {
	//启动服务器
	server.Start()
	fmt.Println("服务器运行中...")
}
