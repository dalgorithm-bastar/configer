package main

import (
	"fmt"
	"github.com/configcenter/internal/log"
	manage "github.com/configcenter/pkg/manager"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"google.golang.org/grpc"
	"net"
)

//配置中心服务端配置文件存放的路径，应与客户端命令行工具同步修改
const (
	logConfigLocation  = "config/log.json"
	grpcConfigLocation = "config/grpc.json"
	etcdConfigLocation = "config/etcdClientv3.json"
)

//
func main() {
	//main函数中出现错误直接panic
	//初始化日志文件
	err := log.NewLogger(logConfigLocation)
	if err != nil {
		panic(err)
	}
	//log.Sugar().Info("log init success")

	//初始化repository，服务端为etcd模式
	err = repository.NewStorage(repository.EtcdType, etcdConfigLocation)
	if err != nil {
		panic(err)
	}

	//初始化manager
	err = manage.NewManager(grpcConfigLocation)
	if err != nil {
		panic(err)
	}

	//监听grpc端口
	manager := manage.GetManager()
	listen, err := net.Listen("tcp", ":"+manage.GetGrpcInfo().Port)
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	pb.RegisterConfigCenterServer(srv, manager)
	fmt.Println("Server started")
	srv.Serve(listen)
}
