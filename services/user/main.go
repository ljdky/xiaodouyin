package main

import (
	"net"
	user "xiaodouyin/kitex_gen/user/userservice"
	"xiaodouyin/services/user/db"

	"github.com/cloudwego/kitex/server"
)

func Init() {
	db.Init()
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8899")
	if err != nil {
		panic(err)
	}
	Init()

	svr := user.NewServer(
		new(UserServiceImpl),
		server.WithServiceAddr(addr),
	)

	if err = svr.Run(); err != nil {
		panic(err.Error())
	}
}
