package ziface

// IDataPack 拆包、封包模块，直接面向 TCP 的数据流，解决 TCP 的粘包问题
type IDataPack interface {
	// 获取包头的长度
	GetHeadLen() uint32

	// 封包方法
	Pack(msg IMessage) ([]byte, error)

	// 拆包方法
	UnPack([]byte) (IMessage, error)
}
