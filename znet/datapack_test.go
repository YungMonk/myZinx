package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// TestDataPack 只负责 datapack 的拆包，封包的单元测试
func TestDataPack(T *testing.T) {

	/**
	 * 模拟的服务器
	 * 1.创建 SocketTCP
	 * 2.从客户端读取数据，拆包处理（创建一个 go 来承载客户端业务）
	 */
	listener, err := net.Listen("tcp4", "0.0.0.0:7777")
	if err != nil {
		fmt.Println("Server listen is error:", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error:", err)
				continue
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// 拆包的过程
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从 conn 中读取 Head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("Read head error", err)
						break
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error ", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						// msg是有数据的，需要第二次读取
						// 第二次从 conn 中读取出 data 内容（根据 DataLen 的长度）
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						// 根据 DataLen 长度，再次从 io 流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error ", err)
							return
						}

						// 完整的消息已经读取完毕
						fmt.Println("-----> Recv MsgId=", msg.ID, ",DataLen=", msg.DataLen, ",Data=", string(msg.Data))
					}

					time.Sleep(time.Second * 1)

				}

			}(conn)
		}
	}()

	/**
	 * 模拟的客户端
	 * 1.创建 SocketTCP
	 * 2.从客户端读取数据，拆包处理
	 */

	conn, err := net.Dial("tcp4", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Client dail error:", err)
		return
	}

	// 创建一个封包 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个 msg 包一起发送

	// 封装第一个 msg 包
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("Clinet pack msg_1 error: ", err)
	}

	// 封装第二个 msg 包
	msg2 := &Message{
		ID:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', ' ', 'h', 'a', 'o', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("Clinet pack msg_2 error: ", err)
	}

	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
