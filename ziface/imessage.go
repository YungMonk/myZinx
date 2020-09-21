package ziface

// IMessage 是将请求的消息封闭到 Message 中所定义的抽象接口
type IMessage interface {

	// 获取消息的ID
	GetMsgID() uint32

	// 获取消息的内容
	GetData() []byte

	// 获取消息的长度
	GetDataLen() uint32

	// 设置消息的ID
	SetMsgID(uint32)

	// 设置消息的内容
	SetData([]byte)

	// 设置消息的长度
	SetDataLen(uint32)
}
