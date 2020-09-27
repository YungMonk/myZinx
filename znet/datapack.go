package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/YungMonk/zinx/utils"
	"github.com/YungMonk/zinx/ziface"
)

// DataPack 具体的封包，拆包方法
type DataPack struct{}

// NewDataPack 拆包、封包的具体实例化
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头的长度
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4个字节) + ID (4个字节)
	return 8
}

// Pack  封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 buffer 的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将 DataLen 写入 databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 将 MsgID 写入 databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 将 Data 数据，写入 databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// UnPack 拆包方法
// （将包的 Head 信息都读取出来，之后再根据 Head 信息中 data 的长度，再进行一次读取）
func (dp *DataPack) UnPack(binarySteam []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的 ioReader
	dataBuff := bytes.NewReader(binarySteam)

	// 只解压 Head 信息，得到 DataLen 和 MsgID
	msg := &Message{}

	// 读取 DataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取 MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	// 判断 DataLen 是否已经超出了 MaxPackageSize 的限制
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, fmt.Errorf("too large msg data recieve")
	}

	// 这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
