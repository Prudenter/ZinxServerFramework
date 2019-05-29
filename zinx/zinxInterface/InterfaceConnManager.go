/**
* @Author: ASlowPerson  
* @Date: 19-5-28 下午9:43
* @mail:18706733725@163.com
*/

package zinxInterface
/*
	定义抽象的链接管理模块,存放链接的集合
*/
type InterfaceConnManager interface {
	//添加链接
	Add(conn InterfaceConnection)
	//删除链接
	Remove(connId uint32)
	//根据链接ID得到链接
	Get(connId uint32)(InterfaceConnection,error)
	//得到目前服务器的链接总数
	GetLen() int
	//清空全部链接的方法
	ClearConn()
}
