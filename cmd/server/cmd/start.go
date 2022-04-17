package cmd

import (
    "context"
    "fmt"
    "net"
    "os"

    "github.com/configcenter/config"
    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/service"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "xchg.ai/sse/gracefully"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
    Use:   "start",
    Short: "start configCenter server",
    Long:  `start configCenter server`,
    Run:   Start,
}

func init() {
    rootCmd.AddCommand(startCmd)
    // Here you will define your flags and configuration settings.
    startCmd.Flags().String(config.EtcdUserName, "", "set etcd username")
    _ = viper.BindPFlag(config.EtcdUserName, startCmd.Flag(config.EtcdUserName))

    startCmd.Flags().String(config.EtcdPassWord, "", "set password to specified etcd user")
    _ = viper.BindPFlag(config.EtcdPassWord, startCmd.Flag(config.EtcdPassWord))

    startCmd.Flags().Int(config.EtcdOperationTimeout, 3, "set etcd timeout by second")
    _ = viper.BindPFlag(config.EtcdOperationTimeout, startCmd.Flag(config.EtcdOperationTimeout))

    startCmd.Flags().String(config.GrpcSocket, "", "set grpc socket")
    _ = viper.BindPFlag(config.GrpcSocket, startCmd.Flag(config.GrpcSocket))

    startCmd.Flags().Int(config.GrpcLockTimeout, 30, "set etcd lock timeout by second when post config")
    _ = viper.BindPFlag(config.GrpcLockTimeout, startCmd.Flag(config.GrpcLockTimeout))

    startCmd.Flags().String(config.LogLogPath, "log/", "set dir for log files")
    _ = viper.BindPFlag(config.LogLogPath, startCmd.Flag(config.LogLogPath))

    startCmd.Flags().String(config.LogFileName, "", "set name of log file")
    _ = viper.BindPFlag(config.LogFileName, startCmd.Flag(config.LogFileName))

    startCmd.Flags().String(config.LogEncodingType, "", "set encoding type for log info, \"json\" for json, \"normal\" for normal")
    _ = viper.BindPFlag(config.LogEncodingType, startCmd.Flag(config.LogEncodingType))

    startCmd.Flags().String(config.LogRecordLevel, "info", "set log level within info and debug")
    _ = viper.BindPFlag(config.LogRecordLevel, startCmd.Flag(config.LogRecordLevel))
    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // startCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Start(cmd *cobra.Command, args []string) {
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
