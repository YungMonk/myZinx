package znet

// Message 实现IMessage
type Message struct {
	ID      uint32 // 消息的Id
	Data    []byte // 消息的内容
	DataLen uint32 // 消息的长度
}

// NewMessage 创建一个 Message 消息包
func NewMessage(msgid uint32, data []byte) *Message {
	return &Message{
		ID:      msgid,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgID 获取消息的ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// GetDataLen 获取消息的长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// SetMsgID 设置消息的ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
