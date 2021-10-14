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

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "create a new config version on remote",
	Long: `post command is used for submit the configfile under the username selected.
you can either assign version number in flag --version or leave it as default
attention that you have to put configfile first`,
	Run: Post,
}

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().StringVarP(&object.Version, "version", "v", "", "put version number here if you want")
}

func Post(cmd *cobra.Command, args []string) {
	//构建请求结构体
	configReq := pb.CfgReq{
		UserName:    object.UserName,
		Target:      []string{object.Version},
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
	resp, err := client.POST(context.Background(), &configReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.Status != "ok" {
		fmt.Println(resp.Status)
		return
	}
	fmt.Println(fmt.Sprintf("Post(Commit) cache of user %s succeed, new version num %s", object.UserName, resp.SliceData[0]))
}
