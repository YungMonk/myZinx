package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/YungMonk/zinx/utils"
	"github.com/YungMonk/zinx/ziface"
)

// Connection 链接模块
type Connection struct {
	// 当前 Conn 隶属于哪个 Server
	TCPServer ziface.IServer

	// 当前链接的 socket TCP 套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经停止/退出的 channel（由 Reader 告知 Writer 退出）
	ExitChan chan bool

	// 告知当前链接已经停止/退出的 channel（由 Reader 告知 Writer 退出）
	ExitBuffChan chan bool

	// 无缓冲通道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	// 有缓冲通道，用于读、写Goroutine之间的消息通信
	msgBuffChan chan []byte

	// 当前 Server 的消息管理模块，用来处理绑定的 MsgID 和对应处理业务API的关系
	MsgHandler ziface.IMsgHandle

	// 链接属性集合
	property map[string]interface{}

	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msghandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		MsgHandler:  msghandler,
		isClosed:    false,
		ExitChan:    make(chan bool, 1),
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    make(map[string]interface{}),
	}

	// 将 connection 加入到 connectionManager中
	c.TCPServer.GetConnMgr().AddConn(c)

	return c
}

// StartReader 链接的读取业务数据方法
func (c *Connection) StartReader() {
	fmt.Println("[StartReader Goroutine is running]")
	defer fmt.Println("[Reader is exit], connID=", c.ConnID, "remote addr is:", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包，解包对象
		dp := NewDataPack()

		// 读取客户端的msg head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			return
		}

		// 拆包，得到 MsgID 和 msgDataLen 放在msg.Data中
		// 第二次拆包得到 MsgData （根据msgDataLen从二进制流中读取MsgData，放入msg.Data中）
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			return
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				return
			}
		}
		msg.SetData(data)

		// 得到当前链接 Conn 的 Request 请求的数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中找到注册绑定的 Conn 对应的 Router 调用
		// 根据绑定好的 MsgID 找到对应的处理 api 业务进行执行
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给工作池处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// StartWriter 写消息Goroutine，专门给客户端发送消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[StartWriter Goroutine is running]")
	defer fmt.Println("[conn Writer exit]", c.RemoteAddr().String())

	// 不断的阻塞等待channel消息，往客户端客发送消息
	for {
		select {
		case data := <-c.msgChan:
			// 有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出，此时 Writer 也要退出
			return
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			// 代表 Reader 已经退出，此时 Writer 也要退出
			return
		}
	}
}

// Start 启动链接，让当前链接开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start ... \n ConnID :", c.ConnID)

	// 启动 从当前的链接读取数据的业务
	go c.StartReader()

	// 启动 往当前的链接写入数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建连接之后要调用的处理业务的函数（对应的 Hook 方法）
	c.TCPServer.CallOnConnStart(c)
}

// Stop 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Closed ... \n ConnID :", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 调用开发者注册的 销毁连接之前 需要执行的业务 hook 方法
	c.TCPServer.CallOnConnStop(c)

	// 关闭 socket 链接
	c.Conn.Close()

	// 告知 Writer 关闭
	c.ExitChan <- true
	// 告知 Writer 关闭
	c.ExitBuffChan <- true

	// 将当前 connection 从 connectionManager 中移除
	c.TCPServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.ExitBuffChan)
	close(c.msgChan)
	close(c.msgBuffChan)
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

// SendMsg 发送数据（将我们要发送给客户端的数据，先进行封包，再发送），将数据发送给远程客户端
func (c *Connection) SendMsg(msgid uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将 data 进行封包
	dp := NewDataPack()
	// binaryMsg格式 MsgDataLen/MsgID/MsgData
	binaryMsg, err := dp.Pack(NewMessage(msgid, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgid)
		return errors.New("Pack error msg")
	}

	// 将数据写回客户端
	c.msgChan <- binaryMsg

	return nil
}

// SendBuffMsg 发送数据（将我们要发送给客户端的数据，先进行封包，再发送），将数据发送给远程客户端
func (c *Connection) SendBuffMsg(msgid uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将 data 进行封包
	dp := NewDataPack()
	// binaryMsg格式 MsgDataLen/MsgID/MsgData
	binaryMsg, err := dp.Pack(NewMessage(msgid, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgid)
		return errors.New("Pack error msg")
	}

	// 将数据写回客户端
	c.msgBuffChan <- binaryMsg

	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	value, ok := c.property[key]
	if !ok {
		return nil, fmt.Errorf("the key=%s not found", key)
	}

	return value, nil
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
