package zinxInterface

/*
	抽象InterfaceRequest,对一次性请求的数据进行封装
*/

type InterfaceRequest interface {
	//得到当前请求的链接
	GetConnection() InterfaceConnection

	//得到链接的数据
	GetData() []byte

	//得到数据的长度
	GetDataLen() int
}