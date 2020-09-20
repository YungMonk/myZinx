package znet

import (
	"fmt"
	"net"

	"github.com/YungMonk/zinx/ziface"

	"github.com/YungMonk/zinx/utils"
)

// Connection 链接模块
type Connection struct {
	// 当前链接的 socket TCP 套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经停止/退出的 channel
	ExitChan chan bool

	// 该链接处理的方法 Router
	Router ziface.IRouter
}

// NewConnection 初始化链接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// StartReader 链接的读取业务数据方法
func (c *Connection) StartReader() {
	fmt.Println("StartReader Goroutine is running ...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote addr is:", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 将远程客户端的数据写入 buffer 中，当前最大长度是 512
		buffer := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Reader recieve data error:", err)
			continue
		}

		// 得到当前链接 Conn 的 Request 请求的数据
		req := Request{
			conn: c,
			data: buffer,
		}

		// 从路由中找到注册绑定的 Conn 对应的 Router 调用
		// 执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// Start 启动链接，让当前链接开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start ... \n ConnID :", c.ConnID)

	// 启动 从当前的链接读取数据的业务
	go c.StartReader()

	// 启动 往当前的链接写入数据的业务

}

// Stop 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Closed ... \n ConnID :", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 关闭 socket 链接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前链接绑定的 socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态，IP，Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
