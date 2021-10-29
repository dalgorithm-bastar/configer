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
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/template"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

// putCmd represents the put command
var putCmd = &cobra.Command{
    Use:   "put",
    Short: "put compressed configfile to remote",
    Long: `put command enables you to save config file under your username on remote
attention that the configfile would not be submit.
learn more about that on command "post"`,
    Run: Put,
}

func init() {
    rootCmd.AddCommand(putCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // putCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // putCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    putCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path of compressedfile(required)")
    putCmd.MarkFlagRequired("pathin")
}

func Put(cmd *cobra.Command, args []string) {
    baseName := filepath.Base(object.PathIn)
    //检测是否为压缩包，提前返回
    if !strings.Contains(baseName, ".tar.gz") && !strings.Contains(baseName, ".zip") {
        fmt.Println("Please input file with format targz or zip")
        return
    }
    file, err := os.Open(object.PathIn)
    if err != nil {
        fmt.Println(fmt.Sprintf("open file at %s err", object.PathIn))
        panic(err)
    }
    binaryFile, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(fmt.Sprintf("read file at %s err", object.PathIn))
        panic(err)
    }
    //构建请求结构体
    configReq := pb.CfgReq{
        UserName: object.UserName,
        Target:   []string{template.CtlFindFlag},
        File: &pb.CompressedFile{
            FileName: baseName,
            FileData: binaryFile,
        },
    }
    //读取grpc配置信息
    err = GetGrpcClient()
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
    resp, err := client.PUT(context.Background(), &configReq)
    if err != nil {
        fmt.Println(err)
        return
    }
    if resp.Status != "ok" {
        fmt.Println(resp.Status)
        return
    }
    fmt.Println(fmt.Sprintf("Put compressedFile %s to remote succeed", baseName))
}
