package zinxNet

import (
	"ZinxServerFramework/zinx/zinxInterface"
	"bytes"
	"encoding/binary"
)

/*
	实现具体的数据封装类
*/

type ZinxDataPack struct {
	zinxInterface.InterfaceDataPack
}

//初始化一个DataPack对象
func NewZinxDataPack() *ZinxDataPack {
	return &ZinxDataPack{}
}

//实现抽象借口的所有方法

//获取二进制包的头部长度,固定返回8字节
func (zdp *ZinxDataPack) GetHeadLen() uint32 {
	//Datalen uint32（4字节) + ID uint32（4字节)
	return 8
}

//数据封包方法---将数据打包成|dataLen|dataID|data|形式
func (zdp *ZinxDataPack)Pack(message zinxInterface.InterfaceMessage)([]byte, error){
	//创建一个存放二进制的字节缓冲区
	dataBuffer := bytes.NewBuffer([]byte{})
	//将dataLen写入buffer中
	if err := binary.Write(dataBuffer,binary.LittleEndian,message.GetMsgLen());err!=nil{
		return nil,err
	}
	//将dataId写入buffer中
	if err := binary.Write(dataBuffer,binary.LittleEndian,message.GetMsgId());err!=nil{
		return nil,err
	}
	//将data写入buffer中
	if err := binary.Write(dataBuffer,binary.LittleEndian,message.GetMsgData());err!=nil{
		return nil,err
	}
	//返回封装好的数据包
	return dataBuffer.Bytes(),nil
}

//拆包方法---将|dataLen|dataID|data|形式的二进制数据拆解到ZinxMessage结构体中
func (zdp *ZinxDataPack) UnPack(binaryData []byte)(zinxInterface.InterfaceMessage, error){
	//解包的时候分2次解压,第一次读取数据头部,包括dataLen和dataID,第二次再根据dataLen读取全部的数据
	message := &ZinxMessage{}
	//创建一个读取二进制数据流的io.Reader
	dataReader := bytes.NewReader(binaryData)
	//读取二进制数据流,将dataLen赋值给message的dataLen
	if err := binary.Read(dataReader,binary.LittleEndian,&message.DataLen);err!=nil{
		return nil,err
	}
	//读取二进制数据流,将dataId赋值给message的Id
	if err := binary.Read(dataReader,binary.LittleEndian,&message.Id);err!=nil{
		return nil,err
	}
	//返回第一次拆包后的数据
	return message,nil
}
