package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/configcenter/internal/log"
	manage "github.com/configcenter/pkg/manager"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"google.golang.org/grpc"
	"xchg.ai/sse/gracefully"
)

//配置中心服务端配置文件存放的路径，应与客户端命令行工具同步修改
const (
	logConfigLocation  = "config/log.json"
	grpcConfigLocation = "config/grpc.json"
	etcdConfigLocation = "config/etcdClientv3.json"
)

func main() {
	//启动过程中出现错误直接panic
	//初始化日志文件
	err := log.NewLogger(logConfigLocation)
	if err != nil {
		panic(err)
	}
	gracefully.Log = log.Zap()
	processCtx, processCancel := context.WithCancel(gracefully.Background())
	log.Sugar().Debug("log init success")

	//初始化repository，服务端为etcd模式
	err = repository.NewStorage(processCtx, repository.EtcdType, etcdConfigLocation)
	if err != nil {
		panic(err)
	}

	//初始化manager
	err = manage.NewManager(processCtx, grpcConfigLocation)
	if err != nil {
		panic(err)
	}

	//监听grpc端口
	manager := manage.GetManager()
	listen, err := net.Listen("tcp", manage.GetGrpcInfo().Socket)
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	pb.RegisterConfigCenterServer(srv, manager)

	gracefully.Go(func() {
		srv.Serve(listen)
		fmt.Println("Server shut down")
	})
	gracefully.Go(func() {
		select {
		case <-processCtx.Done():
			srv.GracefulStop()
		}
	})
	gracefully.RegisterExitHandler(func() {
		fmt.Println("Quit with all goroutine cleaned")
	})
	fmt.Println("Server start")
	gracefully.Wait()
	processCancel()
	os.Exit(gracefully.ExitCode())
}
