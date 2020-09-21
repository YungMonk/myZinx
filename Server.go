package main

import (
	"fmt"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/znet"
)

// PingRouter 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// PreHandle 处理 Connection 业务之前的钩子方法 Hook
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle ...")

	_, err := request.GetConnection().GetTCPConnection().Write([]byte("befor ping..."))
	if err != nil {
		fmt.Println("call back before ping error", err)
	}
}

// Handle 处理 Connection 主业务的钩子方法 Hook
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router handle ...")

	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping..."))
	if err != nil {
		fmt.Println("call back ping error", err)
	}
}

// PostHandle 处理 Connection 业务之后的钩子方法 Hook
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle ...")

	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error", err)
	}
}

func main() {
	// 1.创建Server句柄，使用zinx的api
	s := znet.NewServer("[zinx.v0.1]")

	// 2.给服务注册路由
	s.AddRouter(0, &PingRouter{})

	// 3.启动Server
	s.Serve()
}
