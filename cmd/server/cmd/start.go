package cmd

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/define"
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
	startCmd.Flags().String(define.EtcdUserName, "", "set etcd username")
	_ = viper.BindPFlag(define.EtcdUserName, startCmd.Flag(define.EtcdUserName))

	startCmd.Flags().String(define.EtcdPassWord, "", "set password to specified etcd user")
	_ = viper.BindPFlag(define.EtcdPassWord, startCmd.Flag(define.EtcdPassWord))

	startCmd.Flags().Int(define.EtcdOperationTimeout, 3, "set etcd timeout by second")
	_ = viper.BindPFlag(define.EtcdOperationTimeout, startCmd.Flag(define.EtcdOperationTimeout))

	startCmd.Flags().String(define.GrpcSocket, "", "set grpc socket")
	_ = viper.BindPFlag(define.GrpcSocket, startCmd.Flag(define.GrpcSocket))

	startCmd.Flags().Int(define.GrpcLockTimeout, 30, "set etcd lock timeout by second when post config")
	_ = viper.BindPFlag(define.GrpcLockTimeout, startCmd.Flag(define.GrpcLockTimeout))

	startCmd.Flags().String(define.LogLogPath, "log/", "set dir for log files")
	_ = viper.BindPFlag(define.LogLogPath, startCmd.Flag(define.LogLogPath))

	startCmd.Flags().String(define.LogFileName, "", "set name of log file")
	_ = viper.BindPFlag(define.LogFileName, startCmd.Flag(define.LogFileName))

	startCmd.Flags().String(define.LogEncodingType, "", "set encoding type for log info, \"json\" for json, \"normal\" for normal")
	_ = viper.BindPFlag(define.LogEncodingType, startCmd.Flag(define.LogEncodingType))

	startCmd.Flags().String(define.LogRecordLevel, "info", "set log level within info and debug")
	_ = viper.BindPFlag(define.LogRecordLevel, startCmd.Flag(define.LogRecordLevel))
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
	gracefully.Log = log.Logger
	processCtx, processCancel := context.WithCancel(gracefully.Background())

	//初始化repository，服务端为etcd模式
	err := repository.NewStorage(processCtx, define.EtcdType, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//初始化manager
	err = service.NewManager(processCtx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//监听grpc端口
	manager := service.GetManager()
	listen, err := net.Listen("tcp", service.GetGrpcInfo().Socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
