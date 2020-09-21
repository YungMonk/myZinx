package ziface

import "net"

// IConnection 定义链接模块的抽象层
type IConnection interface {
	// 启动链接，让当前链接开始工作
	Start()

	// 停止链接，结束当前链接的工作
	Stop()

	// 获取当前链接绑定的 socket conn
	GetTCPConnection() *net.TCPConn

	// 获取当前链接模块的链接ID
	GetConnID() uint32

	// 获取远程客户端的TCP状态，IP，Port
	RemoteAddr() net.Addr

	// 发送数据，将数据发送给远程客户端
	SendMsg(msgid uint32, data []byte) error
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
