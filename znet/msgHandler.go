package znet

import (
	"fmt"

	"github.com/YungMonk/zinx/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	// 存放每个msg对应的Router集合
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandler 初始化/创建 MsgHandler 方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的 Router 消息处理方法
func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 1.从 Request 中取出 MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api=", request.GetMsgID(), " is NOT FOUND, need register!")
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
