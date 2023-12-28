package ainterface

type IConnection interface {
	// 启动连接，让当前连接开始工作
	Start()
	// 停止链接，结束当前连接状态
	Stop()
}
