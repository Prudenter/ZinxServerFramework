package zinxInterface
/*
	定义抽象的消息接口
*/
type InterfaceMessage interface {
	//getter方法
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	//setter方法
	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
