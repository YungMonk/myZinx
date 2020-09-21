package ziface

// IRequest 把客户端的请求链接和请求数据包装到了一个Request中
type IRequest interface {
	// 获取请求的链接
	GetConnection() IConnection

	// 获取请求的数据
	GetData() []byte

	// 获取消息ID
	GetMsgID() uint32
}
