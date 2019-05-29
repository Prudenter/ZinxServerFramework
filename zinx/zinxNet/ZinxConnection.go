package zinxNet

import (
	"net"
	"ZinxServerFramework/zinx/zinxInterface"
	"fmt"
	"io"
	"errors"
	"ZinxServerFramework/zinx/utils"
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

	//添加一个channel,用于Reader和Writer之间通信
	messageChan chan []byte

	//添加一个channel,用于Reader通知writer conn已经关闭,writer需要退出
	writerExitChan chan bool

	//当前链接是属于哪个server创建的
	server zinxInterface.ZinxInterfaceServer
}

/*
	初始化链接方法
*/
func NewZinxConnection(conn *net.TCPConn, connId uint32, msgHandler zinxInterface.InterfaceMsgHandler,
	server zinxInterface.ZinxInterfaceServer) zinxInterface.InterfaceConnection {
	zc := &ZinxConnection{
		Conn:           conn,
		ConnID:         connId,
		IsClosed:       false,
		MsgHandler:     msgHandler,
		messageChan:    make(chan []byte),
		writerExitChan: make(chan bool),
		server:         server,
	}

	//当已经成功创建一个链接的时候,将链接添加到链接管理器中
	zc.server.GetConnMgr().Add(zc)

	return zc
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
		//在此调用MsgHander的调度路由方法处理核心业务
		//将request交给worker工作池来处理
		if utils.Globj.WorkerPoolSiz > 0 {
			zc.MsgHandler.SendMsgToTaskQueue(request)
		} else {
			go zc.MsgHandler.DoMsgHandler(request)
		}
	}
}

/*
	写消息的goroutine,负责专门给客户端发送消息
*/
func (zc *ZinxConnection) StartWriter() {
	fmt.Println("[Writer Gotoutine is Statted...]")
	defer fmt.Println("[Writer Goroutine Stop...]")
	//循环监听channel中的数据流动,有数据就写给客户端
	for {
		select { //IO多路复用
		case data := <-zc.messageChan:
			if _, err := zc.Conn.Write(data); err != nil {
				fmt.Println("write data err :", err)
				return
			}
		case <-zc.writerExitChan:
			//代表reader已经退出了,writer也要退出
			return
		}
	}
}

//实现抽象接口的方法

//启动链接
func (zc *ZinxConnection) Start() {
	fmt.Println("Conn start()...id=", zc.ConnID)
	//先进行读业务,添加go程,将读写进行分离
	go zc.StartReader()
	//进行写业务
	go zc.StartWriter()

	//调用创建链接之后用户自定义的Hook
	zc.server.CallOnConnStart(zc)
}

//停止链接
func (zc *ZinxConnection) Stop() {
	fmt.Println("c.stop()...ConnId=", zc.ConnID)

	//调用销毁链接之前用户自定义的Hook函数
	zc.server.CallOnConnStop(zc)

	//回收工作
	if zc.IsClosed == true {
		return
	}
	zc.IsClosed = true

	//在链接关闭后告诉writer
	zc.writerExitChan <- true

	//关闭原生的套接字
	zc.Conn.Close()

	//将当前链接从链接管理模块中删除
	zc.server.GetConnMgr().Remove(zc.ConnID)

	//释放channel资源
	close(zc.messageChan)
	close(zc.writerExitChan)
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
	//将封包好的二进制数据发送给channel,让writer去写给客户端
	zc.messageChan <- binaryMessage
	return nil
}
