package ziface

// IRouter 路由的抽象接口
// 路由中的数据都 IRequest
type IRouter interface {
	// 处理 Connection 业务之前的钩子方法 Hook
	PreHandle(request IRequest)
	// 处理 Connection 主业务的钩子方法 Hook
	Handle(request IRequest)
	// 处理 Connection 业务之后的钩子方法 Hook
	PostHandle(request IRequest)
}
