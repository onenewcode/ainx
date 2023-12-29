package ainterface

// 定义服务器接口
type IServer interface {
	//启动服务器方法
	Start()
	//停止服务器方法
	Stop()
	//开启业务服务方法
	Serve()
	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(router IRouter)
	// todo 路由分组 未来目标 添加类似hertz Group分组

}
