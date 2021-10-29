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
    "io"
    "os"
    "path/filepath"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/template"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
    Use:   "get",
    Short: "get file from remote",
    Long:  `get file from remote, the target can be templates, servicelist, pubilcinfo, maniputation`,
    Example: `$ ./cfgtool get --target servicelist --version 1.0.0 --env 00 --cluster EzEI.set0 --pathout /home/someuser
$ ./cfgtool get -t infrastructure -v 1.0.0 -e 00 -c EzEI.set0 -o /home/someuser`,
    Run: Get,
}

func init() {
    rootCmd.AddCommand(getCmd)
    getCmd.Flags().StringVarP(&object.Target, "target", "t", "", "assign file type needed(required)")
    getCmd.Flags().StringVarP(&object.Version, "version", "v", "", "assign a config version(required)")
    getCmd.Flags().StringVarP(&object.Env, "env", "e", "", "assign an environment number")
    getCmd.Flags().StringVarP(&object.Cluster, "cluster", "c", "", "assign a cluster name")
    getCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path")
}

func Get(cmd *cobra.Command, args []string) {
    //获取Target关键字
    switch object.Target {
    case template.Templates, template.Services, template.Manipulations, template.Infrastructure, template.DeploymentInfo:
        if object.PathOut == "" {
            fmt.Println(fmt.Sprintf("Path required when target is %s", object.Target))
            return
        }
    case template.Versions, template.Environments, template.Clusters:
    //默认返回错误
    default:
        fmt.Println(fmt.Sprintf("Target of %s can not be recognized", object.Target))
        return
    }
    //构建请求结构体
    configReq := pb.CfgReq{
        UserName: object.UserName,
        Target:   []string{object.Target},
        File:     nil,
        CfgVersions: []*pb.CfgVersion{
            {
                Version: object.Version,
                Envs: []*pb.Environment{
                    {
                        Num: object.Env,
                        Clusters: []*pb.Cluster{
                            {
                                ClusterName: object.Cluster,
                                Nodes: []*pb.Node{
                                    {
                                        GlobalId: object.GlobalId,
                                        LocalId:  object.LocalId,
                                        Template: object.TemplateName, //可传空
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
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
    resp, err := client.GET(context.Background(), &configReq)
    if err != nil {
        fmt.Println(err)
        return
    }
    if resp.Status != "ok" {
        fmt.Println(resp.Status)
        return
    }
    switch object.Target {
    case template.Templates, template.Services, template.Manipulations, template.Infrastructure, template.DeploymentInfo:
        //创建文件夹和文件
        dirPath := filepath.Dir(object.PathIn)
        if object.PathOut != "" {
            dirPath = object.PathOut
        }
        err = os.MkdirAll(dirPath, os.ModePerm)
        if err != nil {
            fmt.Println(err)
            return
        }
        f, err := os.OpenFile(dirPath+"/"+resp.File.FileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
        defer f.Close()
        if err != nil {
            fmt.Println(err)
            return
        }
        n, err := f.Write(resp.File.FileData)
        if err == nil && n < len(resp.File.FileData) {
            err = io.ErrShortWrite
            fmt.Println(err)
            return
        }
    case template.Versions, template.Environments, template.Clusters:
        fmt.Println(resp.SliceData)
        return
    }
}
