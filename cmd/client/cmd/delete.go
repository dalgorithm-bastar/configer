/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/configcenter/pkg/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete configfile under the selected username",
	Long:  `using this command to delete ALL configfiles under the selected username`,
	Run:   Delete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func Delete(cmd *cobra.Command, args []string) {
	//构建请求结构体
	configReq := pb.CfgReq{
		UserName:    object.UserName,
		Target:      nil,
		File:        nil,
		CfgVersions: nil,
	}
	//读取grpc配置信息
	err := GetGrpcClient()
	if err != nil {
		panic(err)
	}
	//新建grpc客户端
	conn, err := grpc.Dial(GrpcInfo.Socket, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := pb.NewConfigCenterClient(conn)
	resp, err := client.DELETE(context.Background(), &configReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.Status != "ok" {
		fmt.Println(resp.Status)
		return
	}
	fmt.Println(fmt.Sprintf("Delete cache of user %s succeed", object.UserName))
}
