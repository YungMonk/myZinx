package znet

import "github.com/YungMonk/zinx/ziface"

// BaseRouter 实现 IRouter 接口时，先嵌入 BaseRouter 基类，然后根据这个基类进行重写就可以了
type BaseRouter struct{}

/**
 * 这里之所以 BaseRouter 方法都为空
 * 是因为有的 Router 不希望有 PreHandle、PostHandle 这个业务
 * 所以 Router 全部继承 BaseRouter 的好处就是不需要实现 PreHandle、PostHandle
 */

// PreHandle 处理 Connection 业务之前的钩子方法 Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 处理 Connection 业务之后的钩子方法 Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
