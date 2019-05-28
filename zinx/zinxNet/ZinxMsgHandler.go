package zinxNet

import (
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
	"ZinxServerFramework/zinx/utils"
)

/*
	实现具体的消息管理类
*/
type ZinxMsgHandler struct {
	zinxInterface.InterfaceMsgHandler
	//存放路由集合的map,存放开发者全部的业务,key是消息id,value是用户自定义的路由
	MsgRouter map[uint32]zinxInterface.InterfaceRouter
	//负责存储消息队列的切片,一个worker对应一个任务队列
	TaskQueue []chan zinxInterface.InterfaceRequest
	//worker工作池的worker数量
	WorkerPoolSize uint32
}

//定义初始化方法
func NewZinxMsgHandler() zinxInterface.InterfaceMsgHandler {
	return &ZinxMsgHandler{
		//给map开辟空间
		MsgRouter:      make(map[uint32]zinxInterface.InterfaceRouter),
		WorkerPoolSize: utils.Globj.WorkerPoolSiz,
		//给消息队列的切片初始化
		TaskQueue:      make([]chan zinxInterface.InterfaceRequest, utils.Globj.WorkerPoolSiz),
	}
}

//实现抽象借口的所有方法

//添加路由到map集合中
func (zmh *ZinxMsgHandler) AddRouter(messageId uint32, router zinxInterface.InterfaceRouter) {
	//1.判断messageId是否已经存在
	if _, ok := zmh.MsgRouter[messageId]; ok {
		fmt.Println("repeat key messsageId=", messageId)
		return
	}
	//2.添加messageId和router的对应关系
	zmh.MsgRouter[messageId] = router
	fmt.Println("Add MsgRouter messageId = ", messageId, "success!")
}

//根据messageId 调度路由
func (zmh *ZinxMsgHandler) DoMsgHandler(request zinxInterface.InterfaceRequest) {
	//从request中获取当前的messageId
	messageId := request.GetMessage().GetMsgId()
	//判断MsgRouter中是否已添加messageId
	router, ok := zmh.MsgRouter[messageId];
	if !ok {
		fmt.Println("Handler MsgRouter messageId=", messageId, "Not Found!Need Add!")
		return
	}
	//根据当前messageId,调度其对应的路由处理业务
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

//定义启动worker工作池的方法
func (zmh *ZinxMsgHandler) StartWorkerPool() {
	fmt.Println("WorkPool is started..")
	//根据WorkPoolSize创建worker goroutine
	for i := 0; i < int(zmh.WorkerPoolSize); i++ {
		//1.给当前的Worker所绑定的消息channel对象开辟空间,第i个worker就用第i个channel
		zmh.TaskQueue[i] = make(chan zinxInterface.InterfaceRequest, utils.Globj.MaxWorkerTaskLen)
		//2.启动一个Worker goroutine,阻塞等地消息从对应的管道传进来
		go zmh.startOneWorker(i, zmh.TaskQueue[i])
	}
}

//一个worker真正处理业务的goroutine函数
func (zmh *ZinxMsgHandler) startOneWorker(workerId int, taskQueue chan zinxInterface.InterfaceRequest) {
	fmt.Println("worker Id=", workerId, "is starting...")
	//不断的从对应的管道等待数据传入
	for {
		select {
		case request := <-taskQueue:
			zmh.DoMsgHandler(request)
		}
	}
}

//将消息添加到worker工作池中,并发送给对应的消息队列
func (zmh *ZinxMsgHandler) SendMsgToTaskQueue(request zinxInterface.InterfaceRequest) {
	//1.将消息平均分配给worker,确定当前的request到底要给哪个worker来处理
	//1个客户端绑定一个worker来处理任务
	workerId := request.GetConnection().GetConnID() % zmh.WorkerPoolSize

	//2.直接将request发送给对应的worker的TaskQueue
	zmh.TaskQueue[workerId] <- request
}
