package zinxNet

import (
	"testing"
	"fmt"
	"net"
	"io"
)

/*
	定义数据封装类的单元测试
*/

func TestZinxDataPack_Pack(t *testing.T) {
	fmt.Println("test datapack...")
	// 模拟一个server,收到二进制流数据进行解包
	// 1.创建listenner
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listenner err:", err)
		return
	}
	// 2.开启监听
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server Accept err:", err)
				return
			}
			//3.读取客户端的请求,拆包过程
			go func(conn *net.Conn) {
				//获取一个数据封装类对象
				zdp := NewZinxDataPack()
				for {
					//第一次从conn读,把headData读出来,获取二进制流数据
					headData := make([]byte, zdp.GetHeadLen()) //读取8个字节的数据
					_, err := io.ReadFull(*conn, headData)      //直到headData填充满,才会返回,否则阻塞
					if err != nil {
						fmt.Println("server ReadFull headData err:", err)
						break
					}
					//将headData拆包,得到dataLen和dataID
					headMessage, err := zdp.UnPack(headData)
					if err != nil {
						fmt.Println("server UnPack err:", err)
						return
					}
					//判断数据是否为空
					if headMessage.GetMsgLen() > 0 {
						//如果数据去有内容,则进行第二次读取,获取data
						//注意,headMessage类型时InterfaceMessage,所以需要先进行类型断言,才能进一步操作
						message := headMessage.(*ZinxMessage)
						//给message的data属性开辟空间,长度就是拆包时得到的长度
						message.Data = make([]byte, message.GetMsgLen())
						//根据dataLen的长度进行第二次数据读取
						_, err := io.ReadFull(*conn, message.Data) //直到data填充满,才会返回,否则阻塞
						if err != nil {
							fmt.Println("server ReadFull data err:", err)
							break
						}
						fmt.Println("---> Receive MsgId =", message.Id, "dataLen=", message.DataLen, "data=", string(message.Data))
					}
				}
			}(&conn)
		}
	}()

	/*
		模拟一个client,进行封包发包操作
	*/
	conn,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dail err:",err)
		return
	}
	//封包
	//获取一个数据封装类对象
	zdp := NewZinxDataPack()
	//给数据包赋值
	msg1 := &ZinxMessage{
		Id:1,
		DataLen:4,
		Data:[]byte{'z','i','n','x'},
	}
	//封装第一个数据包
	sendData1,err := zdp.Pack(msg1)
	if err != nil {
		fmt.Println("client Pack err:",err)
		return
	}
	//给数据包赋值
	msg2 := &ZinxMessage{
		Id:2,
		DataLen:5,
		Data:[]byte{'h','e','l','l','o'},
	}
	//封装第一个数据包
	sendData2,err := zdp.Pack(msg2)
	if err != nil {
		fmt.Println("client Pack err:",err)
		return
	}

	//将两个包黏在一起发送
	sendData1 = append(sendData1,sendData2...)//打散再添加
	//发送
	_,err = conn.Write(sendData1)
	if err != nil {
		fmt.Println("client Write err:",err)
		return
	}

	//让test不结束
	select {

	}
}
