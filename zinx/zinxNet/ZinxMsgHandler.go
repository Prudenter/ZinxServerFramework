package zinxNet

import (
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
)

/*
	实现具体的消息管理类
*/
type ZinxMsgHandler struct {
	zinxInterface.InterfaceMsgHandler
	//存放路由集合的map,存放开发者全部的业务,key是消息id,value是用户自定义的路由
	MsgRouter map[uint32]zinxInterface.InterfaceRouter
}

//定义初始化方法
func NewZinxMsgHandler() zinxInterface.InterfaceMsgHandler {
	return &ZinxMsgHandler{
		//给map开辟空间
		MsgRouter:make(map[uint32]zinxInterface.InterfaceRouter),
	}
}

//实现抽象借口的所有方法

//添加路由到map集合中
func (zmh *ZinxMsgHandler)AddRouter(messageId uint32,router zinxInterface.InterfaceRouter){
	//1.判断messageId是否已经存在
	if _, ok := zmh.MsgRouter[messageId]; ok {
		fmt.Println("repeat key messsageId=",messageId)
		return
	}
	//2.添加messageId和router的对应关系
	zmh.MsgRouter[messageId] = router
	fmt.Println("Add MsgRouter messageId = ",messageId,"success!")
}

//根据messageId 调度路由
func(zmh *ZinxMsgHandler)DoMsgHandler(request zinxInterface.InterfaceRequest){
	//从request中获取当前的messageId
	messageId := request.GetMessage().GetMsgId()
	//判断MsgRouter中是否已添加messageId
	router, ok := zmh.MsgRouter[messageId];
	if !ok {
		fmt.Println("Handler MsgRouter messageId=",messageId,"Not Found!Need Add!")
		return
	}
	//根据当前messageId,调度其对应的路由处理业务
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}