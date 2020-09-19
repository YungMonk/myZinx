package znet

import (
	"zinx/ziface"
)

// Server iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string

	// 绑定的ip版本
	IPVersion string

	// 绑定的ip
	IP string

	// 监听的端口
	Port int
}

// Start 实现iServer中的Start()方法
func (s *Server) Start() {

}

// Stop 实现接口iServer中的Stop()方法
func (s *Server) Stop() {

}

// Serve 实现接口iServer中的Serve()方法
func (s *Server) Serve() {

}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return s
}
