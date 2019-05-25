package utils

import (
	"io/ioutil"
	"encoding/json"
)

/*
	定义全局配置文件类
*/
type GlobalObj struct {
	//关于server的配置
	Host       string //当前监听的IP
	Port       int    //当前监听的Port
	Name       string //当前zinxServer的名称
	Version    string //当前监框架的版本号
	MaxPackage uint32 //服务器每Read一次的最大长度
}

//定义一个全局的对外的配置对象
var Globj *GlobalObj

//定义init方法,只要导入当前包,就会执行init方法,加载配置文件
func init() {
	//初始化全局配置对象,设置默认值
	Globj = &GlobalObj{
		Name:       "ZinxServerFramework",
		Host:       "0.0.0.0",
		Port:       8888,
		Version:    "V_0.4",
		MaxPackage: 512,
	}
	//加载开发者自定义的配置文件
	//Globj.LoadConfig()
}

//定义一个加载配置文件的方法
func (globj *GlobalObj) LoadConfig() {
	//读取配置文件路径下的json配置文件,这个路径是我们规定给开发者的,是其main函数所在路径的相对路径
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//解析zinx.json文件,将其中的数据赋值给Globj
	err = json.Unmarshal(data, &Globj)	//这里为什么还要取地址
	if err != nil {
		panic(err)
	}
}
