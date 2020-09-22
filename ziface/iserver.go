package ziface

// IServer 定义一个服务器
type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 运行服务器
	Serve()

	// 路由方法，给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(MsgID uint32, router IRouter)

	// 获取当前服务的连接管理器
	GetConnMgr() IConnManager

	// 注册 OnConnStart 钩子函数的方法
	SetOnConnStart(func(conn IConnection))

	// 注册 OnConnStop 钩子函数的方法
	SetOnConnStop(func(conn IConnection))

	// 调用 OnConnStart 钩子函数的方法
	CallOnConnStart(conn IConnection)

	// 调用 OnConnStop 钩子函数的方法
	CallOnConnStop(conn IConnection)
}
