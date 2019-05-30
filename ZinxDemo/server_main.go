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
func (this *PingRouter) PreHandle(requset zinxInterface.InterfaceRequest) {
	fmt.Println("The PreHandle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(200, []byte("before ping ..."))
	if err != nil {
		fmt.Println("before Send err:", err)
		return
	}
}

//真正处理业务要调用的方法
func (this *PingRouter) Handle(requset zinxInterface.InterfaceRequest) {
	fmt.Println("The Handle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("Send err :", err)
		return
	}
}

//处理业务之后要调用的方法
func (this *PingRouter) PostHandle(requset zinxInterface.InterfaceRequest) {
	fmt.Println("The PostHandle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(200, []byte("after ping ..."))
	if err != nil {
		fmt.Println("after Send err:", err)
		return
	}
}

type HelloRouter struct {
	zinxNet.ZinxRouter
}

//提供自定义的处理业务方法,重写父类方法

//真正处理业务要调用的方法
func (this *HelloRouter) Handle(requset zinxInterface.InterfaceRequest) {
	fmt.Println("The Handle Is Running...")
	//给客户端回写一个数据
	err := requset.GetConnection().Send(201, []byte("Hello zinx!"))
	if err != nil {
		fmt.Println("Send err :", err)
		return
	}
}

//创建链接之后的执行的钩子函数
func DoConntionBegin(conn zinxInterface.InterfaceConnection) {
	fmt.Println("===> DoConntionBegin ...")
	//链接一旦创建成功,就给用户返回一个消息
	if err := conn.Send(202, []byte("Hello welcome to zinx...")); err != nil {
		fmt.Println(err)
		return
	}
	//当用户一旦创建链接成功,就给链接绑定一些属性
	fmt.Println("set conn property...")
	conn.SetProperty("Name","ASlowPerson")
	conn.SetProperty("address","TBD...")
	conn.SetProperty("time","2019/5/31")
}

//链接销毁之前执行的钩子函数
func DoConntionLost(conn zinxInterface.InterfaceConnection) {
	fmt.Println("===> DoConntionLost ...")
	fmt.Println("Conn id", conn.GetConnID(), "is Lost!")

	//获取属性
	fmt.Println("Get conn Property...")
	//获取name
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ",name)
	}
	//获取address
	if address, err := conn.GetProperty("address"); err == nil {
		fmt.Println("address = ",address)
	}
	//获取time
	if time, err := conn.GetProperty("time"); err == nil {
		fmt.Println("time = ",time)
	}
}

func main() {
	//创建一个zinx server对象
	zinxServer := zinxNet.NewServer()

	//注册一个创建链接之后的方法业务
	zinxServer.AddOnConnStart(DoConntionBegin)
	//注册一个销毁链接之前的方法业务
	zinxServer.AddOnConnStop(DoConntionLost)

	// 添加自定义路由到server中,真正处理核心业务的方法在自定义路由里
	//添加不同的自定义路由,处理不同的业务
	zinxServer.AddRouter(1, &PingRouter{})
	zinxServer.AddRouter(2, &HelloRouter{})

	//让server对象启动服务
	zinxServer.Run()
	return
}
