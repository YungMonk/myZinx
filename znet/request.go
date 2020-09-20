package znet

import "github.com/YungMonk/zinx/ziface"

// Request 实现 IRequest 接口
type Request struct {
	// 已经和客户端建立好的链接
	conn ziface.IConnection

	// 客户端请求的数据
	data []byte
}

// GetConnection 获取请求的链接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 获取请求的数据
func (r *Request) GetData() []byte {
	return r.data
}
