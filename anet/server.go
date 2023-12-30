package anet

import (
	"ainx/ainterface"
	"errors"
	"fmt"
	"net"
	"time"
)

type Server struct {
	// 设置服务器名称
	Name string
	// 设置网络协议版本
	IPVersion string
	// 设置服务器绑定IP
	IP string
	// 设置端口号
	Port string
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	Router ainterface.IRouter
	//todo 未来目标提供更多option字段来控制server实例化
}

// ============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// ============== 实现 ainterface.IServer 里的全部接口方法 ========
// 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %s, is starting\n", s.IP, s.Port)

	// 开启一个go去做服务端的Listener业务
	// todo 未来目标是提供更多协议，可以利用if或者switch对IPVersion进行判断而选择采取哪种协议，下面整个方法要重写
	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, s.IP+":"+s.Port)
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}
		// 2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		//	  已经成功监听
		fmt.Println("start Ainx server  ", s.Name, " success, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)
	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}
func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}
func (s *Server) AddRouter(router ainterface.IRouter) {
	s.Router = router
	fmt.Println("Add Router succ! ")
}

/*
创建一个服务器句柄
*/
func NewServer(name string) ainterface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      "8080",
		Router:    nil,
	}
	return s
}
