/**
* @Author: ASlowPerson  
* @Date: 19-5-29 上午10:24
* @mail:18706733725@163.com
*/

package main

import (
	"ZinxServerFramework/ZinxDemo/protobufDemo/pb"
	"github.com/golang/protobuf/proto"
	"fmt"
)

func main() {
	person := &pb.Person{
		Name:   "slowPerson",
		Age:    18,
		Emails: []string{"18706733725@163.com", "1040798027@qq.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "111111111111111",
			},
			&pb.PhoneNumber{
				Number: "222222222222",
			},
			&pb.PhoneNumber{
				Number: "111111111111111",
			},
		},
		//oneof赋值
		//Data:&pb.Person_Socre{
		//	Socre:777,
		//},
		Data:&pb.Person_School{
			School:"xiangongye",
		},
	}
	//将一个protobuf结构体对象 转化成二进制数据
	//任何proto message结构体 在go中他们都是基础Message接口的
	//编码
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	//data就是我们要刚给对端发送的二进制数据
	//解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		fmt.Println("Unmarshal err", err)
		return
	}

	fmt.Println("源数据", person)
	fmt.Println("解码后的数据", newPerson)
	fmt.Println("name = ", newPerson.GetName(), "age = ", newPerson.GetAge(), " emails: ", newPerson.GetEmails())
	fmt.Println("phones = ", newPerson.GetPhones())
	fmt.Println("score = ", newPerson.GetSocre())
	fmt.Println("school = ", newPerson.GetSchool())
}
