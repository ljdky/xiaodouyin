package rpc

import (
	"xiaodouyin/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/client"
)

// 服务
var UserService userservice.Client

// 初始化服务
func init() {
	// 声明错误
	var err error

	// 初始化到 UserService 的 RPC 客户端
	// 暂时本地测试 rpc
	UserService, err = userservice.NewClient("UserService", client.WithHostPorts("localhost:8899"))

	if err != nil {
		panic(err.Error())
	}
}
