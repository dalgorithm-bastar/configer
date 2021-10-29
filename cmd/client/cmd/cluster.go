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
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/template"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
    Use:   "cluster",
    Short: "Create all configfiles of specified cluster",
    Run:   Cluster,
}

func init() {
    rootCmd.AddCommand(clusterCmd)
    clusterCmd.Flags().StringVarP(&object.Version, "version", "v", "", "assign a config version(required)")
    clusterCmd.Flags().StringVarP(&object.Env, "env", "e", "", "assign an environment number(required)")
    clusterCmd.Flags().StringVarP(&object.Cluster, "cluster", "c", "", "assign a cluster name(required)")
    clusterCmd.Flags().StringVarP(&object.TemplateName, "template", "t", "", "assign template(required)")
    clusterCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    clusterCmd.MarkFlagRequired("version")
    clusterCmd.MarkFlagRequired("env")
    clusterCmd.MarkFlagRequired("cluster")
    clusterCmd.MarkFlagRequired("template")
    clusterCmd.MarkFlagRequired("pathout")
}

func Cluster(cmd *cobra.Command, args []string) {
    var deploymentInfo interface{}
    //检测文件夹路径是否合法，提前返回
    err := os.MkdirAll(object.PathOut, os.ModePerm)
    if err != nil {
        fmt.Println(err)
        return
    }
    //新建请求体
    configReq := pb.CfgReq{
        UserName: object.UserName,
        Target:   []string{template.DeploymentInfo},
        CfgVersions: []*pb.CfgVersion{
            {
                Version: object.Version,
                Envs: []*pb.Environment{
                    {
                        Num: object.Env,
                        Clusters: []*pb.Cluster{
                            {
                                ClusterName: object.Cluster,
                            },
                        },
                    },
                },
            },
        },
    }
    //先获取部署信息，拿到实例数目
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
    resp, err := client.GET(context.Background(), &configReq)
    if err != nil {
        fmt.Println(err)
        return
    }
    if resp.Status != "ok" {
        fmt.Println(resp.Status)
        return
    }
    //解析得到的部署信息,拿到实例数目
    err = json.Unmarshal(resp.File.FileData, &deploymentInfo)
    if err != nil {
        log.Sugar().Infof("json unmarshal deploymentinfo err of %v, data:%s", err, string(resp.File.FileData))
        return
    }
    dataMap := make(map[string]string)
    template.ConstructMap(dataMap, deploymentInfo, "")
    if _, ok := dataMap["replicator_number"]; !ok {
        fmt.Println("lack of replicator_number, please checkout servicelist on remote")
        return
    }
    replicatorNum, err := strconv.Atoi(dataMap["replicator_number"])
    if err != nil {
        fmt.Println("err replicator_number of: " + dataMap["replicator_number"])
        return
    }
    //循环生成配置文件
    for i := 0; i < replicatorNum; i++ {
        //新建请求体
        cfgReq := pb.CfgReq{
            UserName: object.UserName,
            Target:   []string{template.NodeConfig},
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
                                            GlobalId: "0",
                                            LocalId:  strconv.Itoa(i),
                                            Template: object.TemplateName,
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
        }
        //发送请求
        cfgResp, err := client.GET(context.Background(), &cfgReq)
        if err != nil {
            fmt.Println(err)
            return
        }
        if cfgResp.Status != "ok" {
            fmt.Println(cfgResp.Status)
            return
        }
        f, err := os.OpenFile(object.PathOut+"/configfile_"+strconv.Itoa(i)+filepath.Ext(object.TemplateName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
        if err != nil {
            fmt.Println(err)
            return
        }
        n, err := f.Write(cfgResp.File.FileData)
        if err == nil && n < len(cfgResp.File.FileData) {
            err = io.ErrShortWrite
            fmt.Println(err)
            return
        }
        f.Close()
    }
}
