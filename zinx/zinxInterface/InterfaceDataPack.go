package zinxInterface

/*
	定义抽象的数据封装借口
*/

type InterfaceDataPack interface {

	//获取二进制包的头部长度,固定返回8字节
	GetHeadLen() uint32

	//数据封包方法---将数据打包成|dataLen|dataID|data|形式
	Pack(im InterfaceMessage)([]byte,error)

	//拆包方法---将|dataLen|dataID|data|形式的数据拆解到ZinxMessage结构体中
	UnPack([]byte)(InterfaceMessage,error)
}
