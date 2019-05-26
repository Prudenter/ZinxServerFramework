package zinxNet

import (
	"net"
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
	"io"
	"errors"
)

//实现具体的TCP链接模块
type ZinxConnection struct {
	zinxInterface.InterfaceConnection

	//当前链接的原生套接字
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//当前的链接状态
	IsClosed bool

	//多路由的消息管理模块
	MsgHandler zinxInterface.InterfaceMsgHandler
}

/*
	初始化链接方法
*/
func NewZinxConnection(conn *net.TCPConn, connId uint32, msgHandler zinxInterface.InterfaceMsgHandler) zinxInterface.InterfaceConnection {
	return &ZinxConnection{
		Conn:       conn,
		ConnID:     connId,
		IsClosed:   false,
		MsgHandler: msgHandler,
	}
}

//针对链接读业务的方法
func (zc *ZinxConnection) StartReader() {
	//从对端读数据
	fmt.Println("Reader go is startin...")
	defer fmt.Println("connId = ", zc.ConnID, "Reader is exit,remote addr is =", zc.GetRemoteAddr().String())
	defer zc.Stop()
	for {
		//创建拆包封包对象
		zdp := NewZinxDataPack()

		//读取客户端消息的头部
		headData := make([]byte, zdp.GetHeadLen())
		if _, err := io.ReadFull(zc.Conn, headData); err != nil {
			fmt.Println("read msg head err :", err)
			break
		}
		//根据头部数据的长度，进行第二次读取
		message, err := zdp.UnPack(headData)
		if err != nil {
			fmt.Println("UnPack err :", err)
			break
		}
		//根据长度，再次读取
		var data []byte
		if message.GetMsgLen() > 0 {
			//有内容
			data = make([]byte, message.GetMsgLen())
			if _, err := io.ReadFull(zc.Conn, data); err != nil {
				fmt.Println("read msg data err :", err)
				break
			}
		}
		//给data赋值
		message.SetData(data)
		//将读出来的message组装成一个request
		request := NewZinxRequest(zc, message)

		//调用用户传递进来的业务处理方法,即自定义router中的业务处理方法--模板设计模式
		//添加go程,防止阻塞
		//在此调用MsgHander的调度路由方法处理核心业务
		go zc.MsgHandler.DoMsgHandler(request)
	}
}

//实现抽象接口的方法

//启动链接
func (zc *ZinxConnection) Start() {
	fmt.Println("Conn start()...id=", zc.ConnID)
	//先进行读业务,添加go程,将读写进行分离
	go zc.StartReader()

	//TODO 进行写业务
}

//停止链接
func (zc *ZinxConnection) Stop() {
	fmt.Println("c.stop()...ConnId=", zc.ConnID)
	//回收工作
	if zc.IsClosed == true {
		return
	}
	zc.IsClosed = true
	//关闭原生的套接字
	zc.Conn.Close()
}

//获取链接ID
func (zc *ZinxConnection) GetConnID() uint32 {
	return zc.ConnID
}

//获取conn的原生socket套接字
func (zc *ZinxConnection) GetTCPConnection() *net.TCPConn {
	return zc.Conn
}

//获取远程客户端的ip地址
func (zc *ZinxConnection) GetRemoteAddr() net.Addr {
	return zc.Conn.RemoteAddr()
}

//发送数据给对方客户端
func (zc *ZinxConnection) Send(messageId uint32, messageData []byte) error {
	//判断当前链接是否关闭
	if zc.IsClosed == true {
		return errors.New("Connection is closed ...send Msg")
	}
	//创建拆包封包对象
	zdp := NewZinxDataPack()
	//将数据封包成二进制数据流形式
	binaryMessage, err := zdp.Pack(NewZinxMessage(messageId, messageData))
	if err != nil {
		fmt.Println("Pack error msg id =", messageId)
		return err
	}
	//将binaryMessage发送给对端
	if _, err := zc.Conn.Write(binaryMessage); err != nil {
		fmt.Println("Write message err:", err)
		return err
	}
	return nil
}
