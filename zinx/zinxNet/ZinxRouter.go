package zinxNet

import "ZinxServerFramework/zinx/zinxInterface"

/*
	实现具体的router模块
*/
type ZinxRouter struct {

	zinxInterface.InterfaceRouter
}

//实现抽象借口的所有方法,将InterfaceRouter接口的方法全部实现,目的是让用户自定义的router继承ZinxRouter后,可以直接重写任意一个方法,而不是要去实现InterfaceRouter接口

//处理业务之前要调用的方法
func (router *ZinxRouter)PreHandle(requset zinxInterface.InterfaceRequest){

}

//真正处理业务要调用的方法
func (router *ZinxRouter)Handle(requset zinxInterface.InterfaceRequest){

}
//处理业务之后要调用的方法
func (router *ZinxRouter)PostHandle(requset zinxInterface.InterfaceRequest){

}
