package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/configcenter/pkg/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete cache under the selected username",
	Long:  `using this command to delete ALL configfiles under the selected username`,
	Run:   Delete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func Delete(cmd *cobra.Command, args []string) {
	if object.UserName == "" {
		fmt.Println("nil username detected, please input username with option -u")
		os.Exit(1)
	}
	//构建请求结构体
	configReq := pb.CfgReq{
		UserName: object.UserName,
	}
	//读取grpc配置信息
	err := GetGrpcClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	//新建grpc客户端
	conn, err := grpc.Dial(GrpcInfo.Socket, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	defer conn.Close()
	client := pb.NewConfigCenterClient(conn)
	resp, err := client.DELETE(context.Background(), &configReq)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	if resp.Status != "ok" {
		fmt.Println(resp.Status)
		os.Exit(1)
		return
	}
	fmt.Println(fmt.Sprintf("Delete cache of user %s succeed", object.UserName))
}
