package znet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YungMonk/zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	// 连接集合
	connections map[uint32]ziface.IConnection

	// 连接锁（保护连接集合的读写锁）
	connLock sync.RWMutex
}

// NewConnManager 创建当前连接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// AddConn 添加连接
func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 connection 加入到 connmanager 中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID=", conn.GetConnID(), " connection add to connmanager successfuly. conn num=", cm.Len())
}

// Remove 移除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将 connection 从 connmanager 中移除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID=", conn.GetConnID(), " connection delete from connmanager successfuly. conn num=", cm.Len())
}

// Get 读取连接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map，加写锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	conn, ok := cm.connections[connID]
	if !ok {
		return nil, errors.New("connection not found")
	}

	return conn, nil
}

// Len 统计连接
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清空连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源 map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除 conn 并停止 conn 工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()

		// 删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear all connections successfuly. conn num=", cm.Len())
}
