package anet

import (
	"ainx/ainterface"
	"ainx/utils"
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
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler ainterface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ainterface.IConnManager

	// =======================
	//新增两个hook函数原型
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ainterface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ainterface.IConnection)

	// =======================
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
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Ainx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalSetting.Version,
		utils.GlobalSetting.MaxConn,
		utils.GlobalSetting.MaxPacketSize)

	// 开启一个go去做服务端的Listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
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
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接,从而控制系统的负载能力
			if s.ConnMgr.Len() >= utils.GlobalSetting.MaxConn {
				conn.Close()
				continue
			}
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}
func (s *Server) Stop() {
	fmt.Println("[STOP] Ainx server , name ", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}
func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}
func (s *Server) AddRouter(msgId uint32, router ainterface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router succ! ")
}

// 得到链接管理
func (s *Server) GetConnMgr() ainterface.IConnManager {
	return s.ConnMgr
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ainterface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ainterface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ainterface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ainterface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

/*
创建一个服务器句柄
*/
func NewServer() ainterface.IServer {
	//先初始化全局配置文件
	utils.GlobalSetting.Reload()

	s := &Server{
		Name:       utils.GlobalSetting.Name, //从全局参数获取
		IPVersion:  "tcp4",
		IP:         utils.GlobalSetting.Host,    //从全局参数获取
		Port:       utils.GlobalSetting.TcpPort, //从全局参数获取
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), //创建ConnManager
	}
	return s
}
