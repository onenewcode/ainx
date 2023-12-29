package config

import "ainx/ainterface"

type Request struct {
	conn ainterface.IConnection //已经和客户端建立好的链接
	date []byte                 //客户端请求数据
}

// 获取请求链接信息
func (r *Request) GetConnection() ainterface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.date
}
