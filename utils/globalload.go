package utils

import (
	"ainx/ainterface"
	"fmt"
	"github.com/spf13/viper"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalSet struct {
	TcpServer ainterface.IServer //当前Zinx的全局Server对象
	Host      string             //当前服务器主机IP
	TcpPort   string             //当前服务器主机监听端口号
	Name      string             //当前服务器名称
	Version   string             //当前Zinx版本号

	MaxPacketSize uint32 //都需数据包的最大值
	MaxConn       uint32 //当前服务器主机允许的最大链接个数
}

// todo 未来支持多种配置文件格式
// 读取用户的配置文件
func (g *GlobalSet) Reload() {
	vp := viper.New()          //创建viper对象
	vp.SetConfigName("config") //配置文件的名称
	vp.AddConfigPath("./")
	vp.SetConfigType("yaml") //配置文件的拓展名
	err := vp.ReadInConfig() //读取配置文件的内容
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中

	err = vp.Unmarshal(&GlobalSetting)
	if err != nil {
		fmt.Println(GlobalSetting.Host)
		return
	}
}

/*
提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalSetting = &GlobalSet{
		Name:          "AinxServerApp",
		Version:       "V0.4",
		TcpPort:       "8080",
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}
}

/*
定义一个全局的对象
*/
var GlobalSetting *GlobalSet
