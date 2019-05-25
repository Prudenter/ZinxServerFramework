package zinxNet

import "ZinxServerFramework/zinx/zinxInterface"

/*
	实现具体的消息模块
*/
type ZinxMessage struct {

	zinxInterface.InterfaceMessage

	Id uint32
	DataLen uint32
	Data []byte
}

//实现抽象接口的所有方法

//getter方法
func (zm *ZinxMessage)GetMsgId() uint32{
	return zm.Id
}

func (zm *ZinxMessage)GetMsgLen() uint32{
	return zm.DataLen
}

func (zm *ZinxMessage)GetMsgData() []byte{
	return zm.Data
}

//setter方法
func (zm *ZinxMessage) SetMsgId(id uint32){
	zm.Id = id
}

func (zm *ZinxMessage) SetData(data []byte){
	zm.Data = data
}

func (zm *ZinxMessage) SetDataLen(len uint32){
	zm.DataLen = len
}