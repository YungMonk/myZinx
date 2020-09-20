package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

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

	// 当前的Server增加一个Router，供 Server 注册的链接来处理对应业务
	Router ziface.IRouter
}

// Start 实现IServer接口中的，Start()方法
func (s *Server) Start() {
	fmt.Printf("[Zinx] ServerName:%s, Linstener at IP:%s, Port is %d, is staring",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TCPPort,
	)

	go func() {

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
				fmt.Printf("accept err:$s\n", err)
				continue
			}

			// 将处理新链接的方法和conn进行绑定得到我们的链接模块
			dealConn := NewConnection(clinetConn, cid, s.Router)
			cid++

			// 启动 当前链接业务处理
			go dealConn.Start()
		}

	}()
}

// Stop 实现IServer接口中的，Start()方法
func (s *Server) Stop() {
	// 将服务器中的资源，状态或者已经开辟的链接信息进行停止或者回收
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
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success")
}

// NewServer 初始化 Server 模块的方法
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: utils.GlobalObject.IPVersion,
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TCPPort,
		Router:    nil,
	}
}
