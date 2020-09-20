package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/YungMonk/zinx/ziface"
)

// GlobalObj 储存有关框架的所有配置，供其它模块使用
// 一些参数可以通过 zinx.json 由用户进行配置
type GlobalObj struct {
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
	Version        string //当前Zinx的版本号
	MaxConn        int    // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前服务器数据包的最大值
}

// GlobalObject 定义一个全局对外的 GlobalObj 对象
var GlobalObject *GlobalObj

// Reload 从 zinx/conf 中加载用户自定义参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// 将 json 文件中的数据解析到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// init 提供一个init方法，初始化当前GlobalObj
func init() {
	// 如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		Name:           "Zinx Server App",
		Host:           "0.0.0.0",
		TCPPort:        8999,
		IPVersion:      "tcp4",
		Version:        "v0.3",
		MaxConn:        1000,
		MaxPackageSize: 4092,
	}

	// 应该尝试从zinx/conf中去加载一些用户自定义的参数
	GlobalObject.Reload()
}
