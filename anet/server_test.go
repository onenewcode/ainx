package config

import (
	"fmt"
	"net"
	"testing"
	"time"
)

/*
模拟客户端
*/
func ClientTest() {
	fmt.Println("Client Test ... start")
	// 3秒之后发起调用，让服务端有时间启动
	time.Sleep(3 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client start err,exit")
		return
	}
	for {
		_, err := conn.Write([]byte("hello word"))
		if err != nil {
			fmt.Println("client start err,exit")
			return
		}
		buf := make([]byte, 520)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf error")
			return
		}
		fmt.Printf("Server call back : %s,cnt =%d \n ", buf[:cnt], cnt)
		time.Sleep(1 * time.Second)

	}
}

// Server 模块测试函数
func TestServer(t *testing.T) {
	/*
		服务端测试
	*/
	//
	s := NewServer("first")
	go ClientTest()
	s.Serve()
}
