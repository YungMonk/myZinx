package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/YungMonk/zinx/ziface"
	"github.com/YungMonk/zinx/zlog"
)

// GlobalObj 储存有关框架的所有配置，供其它模块使用
// 一些参数可以通过 zinx.json 由用户进行配置
type GlobalObj struct {

	/**
	 * Config file path
	 */
	ConfFilePath string // 配置文件的路径

	/**
	 * Server
	 */
	TCPServer ziface.IServer // 当前zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TCPPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器名称
	IPVersion string         // IP类型 tcp4,tcp6

	/**
	 * Zinx
	 */
	Version          string // 当前Zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前服务器数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作 worker 池的数量
	MaxWorkerTaskLen uint32 // 每个 Worker 对应的消息队列的任务的数量最大值（限定条件）(每个消息队列中所存请求数量)
	MaxMsgChanLen    uint32 // SendBuffMsg发送消息的缓冲最大长度

	/**
	 * Logger
	 */
	LogDir        string // 日志所在文件夹              默认"./log"
	LogFile       string // 日志文件名称                默认""        --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   // 是否关闭Debug日志级别调试信息 默认false     -- 默认打开debug信息
	LogLevel      int    // 日志级别调试信息             默认LogDebug     -- 默认打开debug信息
}

// GlobalObject 定义一个全局对外的 GlobalObj 对象
var GlobalObject *GlobalObj

// PathExists 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Reload 从 zinx/conf 中加载用户自定义参数
func (g *GlobalObj) Reload() {

	if configFileExists, _ := PathExists(g.ConfFilePath); !configFileExists {
		fmt.Println("Config file is not exists")
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}

	// 将 json 文件中的数据解析到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

	// logger 设置
	if g.LogDir != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose {
		zlog.CloseDebug()
	}
}

// init 提供一个init方法，初始化当前GlobalObj
func init() {
	// 如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		Name:             "Zinx Server App",
		Host:             "0.0.0.0",
		TCPPort:          8999,
		IPVersion:        "tcp4",
		Version:          "v0.11",
		MaxConn:          1000,
		MaxPackageSize:   4092,
		WorkerPoolSize:   uint32(runtime.NumCPU()),
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/zinx.json",
		MaxMsgChanLen:    1024,
		LogDir:           "./log",
		LogFile:          "",
		LogDebugClose:    false,
		LogLevel:         zlog.LogDebug,
	}

	// 应该尝试从zinx/conf中去加载一些用户自定义的参数
	GlobalObject.Reload()
}
