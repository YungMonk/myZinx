package ziface

// IConnManager 连接管理抽象模块
type IConnManager interface {

	// 添加连接
	AddConn(conn IConnection)

	// 移除连接
	Remove(conn IConnection)

	// 读取连接
	Get(connID uint32) (IConnection, error)

	// 统计连接
	Len() int

	// 清空连接
	ClearConn()
}
