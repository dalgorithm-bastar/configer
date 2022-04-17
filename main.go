package main

import (
    "context"
    "fmt"
    "net"
    "os"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    service "github.com/configcenter/pkg/service"
    "google.golang.org/grpc"
    "xchg.ai/sse/gracefully"
)

func main() {
    //启动过程中出现错误直接panic
    //初始化日志文件
    err := log.NewLogger()
    if err != nil {
        panic(err)
    }
    gracefully.Log = log.Zap()
    processCtx, processCancel := context.WithCancel(gracefully.Background())
    log.Sugar().Debug("log init success")

    //初始化repository，服务端为etcd模式
    err = repository.NewStorage(processCtx, repository.EtcdType, "")
    if err != nil {
        panic(err)
    }

    //初始化manager
    err = service.NewManager(processCtx)
    if err != nil {
        panic(err)
    }

    //监听grpc端口
    manager := service.GetManager()
    listen, err := net.Listen("tcp", service.GetGrpcInfo().Socket)
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
