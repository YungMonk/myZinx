package znet

import (
	"fmt"
	"net"

	"github.com/YungMonk/zinx/utils"
	"github.com/YungMonk/zinx/ziface"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

// Server IServer接口的实现
type Server struct {
	// 服务器名称
	Name string

	// 绑定的IP版本
	IPVersion string

	// 绑定的IP
	IP string

	// 绑定的端口
	Port int

	// 当前 Server 的消息管理模块，用来处理绑定的 MsgID 和对应处理业务API的关系
	MsgHandler ziface.IMsgHandle

	// 当前 Server 的连接管理器
	ConnMgr ziface.IConnManager

	// 当前 Server 创建连接之后自动调用 Hook 函数----OnConnStart
	OnConnStart func(conn ziface.IConnection)

	// 当前 Server 销毁连接之前自动调用 Hook 函数----OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// NewServer 初始化 Server 模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  utils.GlobalObject.IPVersion,
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TCPPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}

// Start 实现IServer接口中的，Start()方法
func (s *Server) Start() {
	fmt.Printf("[Zinx] ServerName:%s, Linstener at IP:%s, Port is %d, is staring\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TCPPort,
	)

	go func() {

		// 0.开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1.获取服务器的 TCP 的 Adress
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 2.监听服务器的地址
		linsten, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("linsten %s, err %+v", s.IPVersion, err)
			return
		}
		fmt.Printf("Start Zinx Server success, %s success, Linstenning...\n", s.Name)

		var cid uint32
		cid = 0

		// 3.阻塞的等待客户端的链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			clinetConn, err := linsten.AcceptTCP()
			if err != nil {
				fmt.Println("accept err:$s\n", err)
				continue
			}

			// 检查是否超过最大连接值
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// @todo 给客户端响应一个超出最大连接数的错误包
				fmt.Println("Too many connections MaxConn=", utils.GlobalObject.MaxConn)
				clinetConn.Close()
				continue
			}

			// 将处理新链接的方法和conn进行绑定得到我们的链接模块
			dealConn := NewConnection(s, clinetConn, cid, s.MsgHandler)
			cid++

			// 启动 当前链接业务处理
			go dealConn.Start()
		}

	}()
}

// Stop 实现IServer接口中的，Start()方法
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name ", s.Name)

	// 将服务器中的资源，状态或者已经开辟的链接信息进行停止或者回收
	s.ConnMgr.ClearConn()
}

// Serve 实现IServer接口中的，Start()方法
func (s *Server) Serve() {
	// 启动Server服务
	s.Start()

	//@todo 做一些启动服务器之后的额外功能或业务

	// 阻塞uddy
	select {}

}

// AddRouter 路由方法，给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(MsgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(MsgID, router)
	fmt.Println("Add Router success")
}

// GetConnMgr 获取当前 Server 的连接管理器
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 注册 OnConnStart 钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册 OnConnStop 钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用 OnConnStart 钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart == nil {
		return
	}

	fmt.Println("[Call] OnStart() had been callback.")
	s.OnConnStart(conn)
}

// CallOnConnStop 调用 OnConnStop 钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop == nil {
		return
	}

	fmt.Println("[Call] OnStart() had been callback.")
	s.OnConnStop(conn)
}

func init() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	fmt.Println(bottomLine)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
}
