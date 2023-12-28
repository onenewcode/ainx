package anet

import (
	"ainx/ainterface"
	"net"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID也可以称作SessionID，ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 处理该链接方法的API
	handleAPI ainterface.HandFunc
	// 告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
}

// 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ainterface.HandFunc) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		handleAPI: callback_api,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}
// 处理conn读数据的Goroutine
func (c *Connection) Start  {

}