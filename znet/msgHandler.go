package znet

import (
	"fmt"

	"github.com/YungMonk/zinx/utils"
	"github.com/YungMonk/zinx/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	// 存放每个msg对应的Router集合
	Apis map[uint32]ziface.IRouter

	// 负责Worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化/创建 MsgHandler 方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// DoMsgHandler 调度/执行对应的 Router 消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1.从 Request 中取出 MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api=", request.GetMsgID(), " is NOT FOUND, need register!")
		return
	}

	// 2.根据对应的 MsgID 调用相关的 router 业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandler) AddRouter(MsgID uint32, router ziface.IRouter) {
	// 1.判断当前的 msgID 是否已经存在于 Apis 中
	if _, ok := mh.Apis[MsgID]; ok {
		panic(fmt.Sprintf("repeat api, msgID = %d", MsgID))
	}

	// 2.添加 msgID 与 router 的绑定关系
	mh.Apis[MsgID] = router

	fmt.Println("Add api MsgID = ", MsgID, " success.")
}

// StartWorkerPool 启动一个 worker 工作池（开起工作池的动作只能发生一次，框架只能有一个工作池）
func (mh *MsgHandler) StartWorkerPool() {
	// 根据 WorkerPoolSize 分别开启 Worker，每个 Worker 用一个 goroutine 来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {

		/**
		 * 一个 Worker 被启动
		 */

		// 1.当前的 Worker 对应的 channel 消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 2.启动当前的 Worker，阻塞消息从 channel 传递过来
		go mh.starOneWorker(i, mh.TaskQueue[i])
	}

}

// 启动一个 worker 工作流程
func (mh *MsgHandler) starOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID=", workerID, " is started.")
	// 不断的阻塞等待对应的消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的 Request，执行当前 Request 所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息发送给消息任务队列(TaskQueue)由Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {

	// 1.将消息平均分配给不同的 worker
	// 根据客户端建立的 connID 来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(),
		",request MsgID=", request.GetMsgID(),
		" to WorkerID=", workerID,
	)

	// 2.将消息发送给 worker 中对应的 TaskQueue 进行处理
	mh.TaskQueue[workerID] <- request
}
