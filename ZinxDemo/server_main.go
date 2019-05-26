package main

import (
	"ZinxServerFramework/zinx/zinxNet"
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
)

/*
	用户自定义router,用来处理与客户端的读写业务
*/

//继承ZinxRouter
type PingRouter struct {
	zinxNet.ZinxRouter
}


//提供自定义的处理业务方法,重写父类方法

//处理业务之前要调用的方法
func (this *PingRouter)PreHandle(requset zinxInterface.InterfaceRequest){
	fmt.Println("The PreHandle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(1,[]byte("before ping ..."))
	if err != nil {
		fmt.Println("before Send err:",err)
		return
	}
}

//真正处理业务要调用的方法
func (this *PingRouter)Handle(requset zinxInterface.InterfaceRequest){
	fmt.Println("The Handle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(1,[]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("Send err :",err)
		return
	}
}

//处理业务之后要调用的方法
func (this *PingRouter)PostHandle(requset zinxInterface.InterfaceRequest){
	fmt.Println("The PostHandle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(1,[]byte("after ping ..."))
	if err != nil {
		fmt.Println("after Send err:",err)
		return
	}
}

func main() {
	//创建一个zinx server对象
	zinxServer := zinxNet.NewServer()

	//TODO 注册一些自定义的业务
	// 添加自定义路由到server中,真正处理核心业务的方法在自定义路由里
	zinxServer.AddRouter(&PingRouter{})

	//让server对象启动服务
	zinxServer.Run()
	return
}
