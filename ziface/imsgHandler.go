package ziface

// IMsgHandle 消息管理抽象层
type IMsgHandle interface {
	// 调度/执行对应的 Router 消息处理方法
	DoMsgHandler(request IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(MsgID uint32, router IRouter)

	// 启动 worker 工作池
	StartWorkerPool()

	// 将消息发送给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
