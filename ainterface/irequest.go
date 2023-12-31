package ainterface

/*
IRequest 接口
实际是把客户端请求链接信息和请求数据包放在Request里
*/
type IRequest interface {
	GetConnection() IConnection //获取请求链接信息
	GetData() []byte            //获取请求消息的数据
	GetMsgID() uint32           //获取消息ID
}
