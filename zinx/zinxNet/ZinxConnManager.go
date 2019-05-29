/**
* @Author: ASlowPerson  
* @Date: 19-5-28 下午9:44
* @mail:18706733725@163.com
*/

package zinxNet

import (
	"ZinxServerFramework/zinx/zinxInterface"
	"sync"
	"fmt"
	"errors"
)

/*
	实现具体的链接管理类
*/
type ZinxConnManager struct {
	zinxInterface.InterfaceConnManager
	//存放所有链接的map,用于管理全部的链接
	connections map[uint32]zinxInterface.InterfaceConnection
	//保护所有链接集合map的锁
	connLock sync.RWMutex
}

//定义初始化方法
func NewZinxConnManager() zinxInterface.InterfaceConnManager {
	return &ZinxConnManager{
		//给map开辟空间
		connections: make(map[uint32]zinxInterface.InterfaceConnection),
	}
}

//添加链接
func (connMgr *ZinxConnManager) Add(conn zinxInterface.InterfaceConnection) {
	//添加写锁
	connMgr.connLock.Lock()
	//解锁
	defer connMgr.connLock.Unlock()
	//添加conn到map中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("Add connId = ", conn.GetConnID(), "to manager success!")
}

//删除链接
func (connMgr *ZinxConnManager) Remove(connId uint32) {
	//添加写锁
	connMgr.connLock.Lock()
	//解锁
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, connId)
	fmt.Println("Remove connId = ", connId, "from manager success!")
}

//根据链接ID得到链接
func (connMgr *ZinxConnManager) Get(connId uint32) (zinxInterface.InterfaceConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not Found!")
	}
}

//得到目前服务器的链接总数
func (connMgr *ZinxConnManager) GetLen() int {
	return len(connMgr.connections)
}

//清空全部链接的方法
func (connMgr *ZinxConnManager) ClearConn() {
	//添加写锁
	connMgr.connLock.Lock()
	//解锁
	defer connMgr.connLock.Unlock()
	//遍历删除后
	for connId, conn := range connMgr.connections {
		//将全部的conn关闭
		conn.Stop()
		//删除链接
		delete(connMgr.connections, connId)
	}
	fmt.Println("ClearConn all connecitons success!conn num = ",connMgr.GetLen())
}
