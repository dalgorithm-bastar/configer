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
    "regexp"
    "strconv"
    "strings"

    "github.com/configcenter/internal/log"
    manage "github.com/configcenter/pkg/manager"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/template"
    "github.com/configcenter/pkg/util"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

var mode string

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
    Use:   "cluster",
    Short: "Create all configfiles of specified config scheme",
    Run:   Cluster,
}

func init() {
    rootCmd.AddCommand(clusterCmd)
    clusterCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path")
    clusterCmd.Flags().StringVarP(&object.Version, "version", "v", "", "assign a config version(required)")
    clusterCmd.Flags().StringVarP(&object.Env, "env", "e", "", "assign an env num of 2 bits(required)")
    clusterCmd.Flags().StringVarP(&object.Config, "config", "s", "", "assign config scheme")
    clusterCmd.Flags().StringVarP(&object.Cluster, "cluster", "c", "", "assign a cluster name")
    clusterCmd.Flags().StringVarP(&object.TemplateName, "template", "t", "", "assign template")
    clusterCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    clusterCmd.Flags().StringVarP(&mode, "mode", "m", "", "input \"remote\" or \"local\" to choose creating from remote or local(required)")
    clusterCmd.MarkFlagRequired("version")
    clusterCmd.MarkFlagRequired("env")
    clusterCmd.MarkFlagRequired("config")
    //clusterCmd.MarkFlagRequired("template")
    clusterCmd.MarkFlagRequired("pathout")
    clusterCmd.MarkFlagRequired("mode")
}

func Cluster(cmd *cobra.Command, args []string) {
    if mode != "local" && mode != "remote" {
        fmt.Println("please input correct arg mode, within \"remote\" or \"local\"")
        return
    }
    //校验环境号是否合法
    envNumFormat, err := regexp.Compile(manage.EnvNumString)
    if err != nil {
        panic(err)
    }
    if !envNumFormat.MatchString(object.Env) {
        panic(fmt.Sprintf("illegal envNum of %s, please input num of 2 bits", object.Env))
    }
    //从远端生成
    if mode == "remote" {
        var deploymentInfo, reflectedDeploymentInfo interface{}
        //检测文件夹路径是否合法，提前返回
        err := os.MkdirAll(object.PathOut, os.ModePerm)
        if err != nil {
            fmt.Println(err)
            return
        }
        //获取集群列表
        //新建请求体
        configReq := pb.CfgReq{
            UserName: object.UserName,
            EnvNum:   object.Env,
            Target:   []string{template.Clusters},
            CfgVersions: []*pb.CfgVersion{
                {
                    Version: object.Version,
                    Confs: []*pb.ConfigScheme{
                        {
                            ConfigName: object.Config,
                            Clusters: []*pb.Cluster{
                                {
                                    ClusterName: "",
                                    Nodes: []*pb.Node{
                                        {
                                            Template: "",
                                            LocalId:  "",
                                            GlobalId: "",
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
        clusterList := resp.SliceData
        for _, cluster := range clusterList {
            configReq.Target[0] = template.Templates
            configReq.CfgVersions[0].Confs[0].Clusters[0].ClusterName = cluster
            //获取该集群下的所有模板名称
            resp, err := client.GET(context.Background(), &configReq)
            if err != nil {
                fmt.Println(err)
                return
            }
            if resp.Status != "ok" {
                fmt.Println(resp.Status)
                return
            }
            templateFileMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            //获取部署信息
            configReq.Target[0] = template.DeploymentInfo
            resp, err = client.GET(context.Background(), &configReq)
            if err != nil {
                fmt.Println(err)
                return
            }
            if resp.Status != "ok" {
                fmt.Println(resp.Status)
                return
            }
            //解析得到的部署信息,拿到实例数目
            serviceListFile := resp.File.FileData
            err = json.Unmarshal(resp.File.FileData, &deploymentInfo)
            if err != nil {
                log.Sugar().Infof("json unmarshal deploymentinfo err of %v, data:%s", err, string(resp.File.FileData))
                return
            }
            dataMap := make(map[string]string)
            template.ConstructMap(dataMap, deploymentInfo, "")
            if _, ok := dataMap[template.ReplicatorNumKey]; !ok {
                fmt.Println("lack of replicator_number, please checkout servicelist on remote")
                return
            }
            replicatorNum, err := strconv.Atoi(dataMap[template.ReplicatorNumKey])
            if err != nil {
                fmt.Println("err replicator_number of: " + dataMap[template.ReplicatorNumKey])
                return
            }
            //取基础设施信息，以便使用主机名命名文件夹
            configReq.Target[0] = template.Infrastructure
            resp, err = client.GET(context.Background(), &configReq)
            if err != nil {
                fmt.Println(err)
                return
            }
            if resp.Status != "ok" {
                fmt.Println(resp.Status)
                return
            }
            infraDataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            infrastructureFile := infraDataMap[util.Join("/", object.Version, repository.Infrastructure)]
            if infrastructureFile == nil {
                panic("no infrastructure file on remote!")
            }
            deploymentDataReflected, err := template.GetDeploymentInfo(object.Env, serviceListFile, infrastructureFile)
            err = json.Unmarshal([]byte(deploymentDataReflected), &reflectedDeploymentInfo)
            if err != nil {
                log.Sugar().Infof("json unmarshal reflecteddeploymentinfo err of %v, data:%s", err, deploymentDataReflected)
                return
            }
            dataMapReflected := make(map[string]string)
            template.ConstructMap(dataMapReflected, reflectedDeploymentInfo, "")
            //循环生成配置文件
            for tmplPath, tmplFile := range templateFileMap {
                if tmplFile == nil || len(tmplFile) == 0 {
                    continue
                }
                tmplNameSlice := strings.Split(tmplPath, "/")
                tmplName := tmplNameSlice[len(tmplNameSlice)-1]
                for i := 0; i < replicatorNum; i++ {
                    configReq.Target[0] = template.NodeConfig
                    configReq.CfgVersions[0].Confs[0].Clusters[0].Nodes[0].Template = tmplName
                    configReq.CfgVersions[0].Confs[0].Clusters[0].Nodes[0].GlobalId = "0"
                    configReq.CfgVersions[0].Confs[0].Clusters[0].Nodes[0].LocalId = strconv.Itoa(i)
                    cfgResp, err := client.GET(context.Background(), &configReq)
                    if err != nil {
                        fmt.Println(err)
                        return
                    }
                    if cfgResp.Status != "ok" {
                        fmt.Println(cfgResp.Status)
                        return
                    }
                    //创建文件夹和文件,以集群名和节点号区分,为便于运维人员识别，节点号映射到主机名
                    hostnameKey := util.Join(".", template.DeploymentInfoKey, strconv.Itoa(i), template.HostNameKey)
                    hostname, ok := dataMapReflected[hostnameKey]
                    if !ok {
                        log.Sugar().Infof("no such hostname under key: %s", hostnameKey)
                        return
                    }
                    dirPath := util.Join("/", object.PathOut, cluster, hostname)
                    err = os.MkdirAll(dirPath, os.ModePerm)
                    if err != nil {
                        fmt.Println(err)
                        return
                    }
                    f, err := os.OpenFile(dirPath+"/"+tmplName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
                    n, err := f.Write(cfgResp.File.FileData)
                    if err == nil && n < len(cfgResp.File.FileData) {
                        err = io.ErrShortWrite
                        fmt.Println(err)
                        return
                    }
                    f.Close()
                }
            }
        }
        return
    }

    //从本地生成
    if object.PathIn == "" {
        fmt.Println("please specify input path to compressed file")
        return
    }
    //获取文件存储路径和内容
    fileName := filepath.Base(object.PathIn)
    if filepath.Ext(fileName) != ".zip" && !strings.Contains(fileName, ".tar.gz") {
        fmt.Println(fmt.Sprintf("unsupport compressed file of \"%s\"", fileName))
    }
    //TODO 文件夹数据
    //初始化repository，客户端为压缩包模式
    err = repository.NewStorage(context.Background(), repository.CompressedFileType, object.PathIn)
    if err != nil {
        fmt.Println(err)
        return
    }
    //生成该包内的所有配置文件
    //获取文件数据
    fileData, ok := repository.Src.GetSourceDataorOperator().(map[string][]byte)
    if !ok {
        panic("source data in compreesed file err")
    }
    //用于记录文件包内所有集群及每个集群对应的所有模板，只有有模板的集群才会被记录
    clusterMap := make(map[string][]string)
    //筛选集群和对应的模板名称
    for keyofPath := range fileData {
        if strings.Contains(keyofPath, "/"+repository.Templates+"/") {
            keySlice := strings.Split(keyofPath, "/")
            if len(keySlice) < 5 {
                continue
            }
            clusterName := keySlice[2]
            templateName := keySlice[4]
            if _, ok := clusterMap[clusterName]; !ok {
                tmplSlice := []string{templateName}
                clusterMap[clusterName] = tmplSlice
            } else {
                clusterMap[clusterName] = append(clusterMap[clusterName], templateName)
            }
        }
    }
    //根据筛选结果逐个集群进行填充
    for clusterName, tmplSlice := range clusterMap {
        //解析得到的部署信息,拿到实例数目
        var deploymentInfo, reflectedDeploymentInfo interface{}
        servicelistFile, err := repository.Src.Get(util.Join("/", object.Version, object.Config, clusterName, repository.ServiceList))
        err = json.Unmarshal(servicelistFile, &deploymentInfo)
        if err != nil {
            log.Sugar().Infof("json unmarshal deploymentinfo err of %v, data:%s", err, string(servicelistFile))
            return
        }
        dataMap := make(map[string]string)
        template.ConstructMap(dataMap, deploymentInfo, "")
        if _, ok := dataMap[template.ReplicatorNumKey]; !ok {
            fmt.Println("lack of replicator_number, please checkout servicelist on remote")
            return
        }
        replicatorNum, err := strconv.Atoi(dataMap[template.ReplicatorNumKey])
        if err != nil {
            fmt.Println("err replicator_number of: " + dataMap[template.ReplicatorNumKey])
            return
        }
        //取部署信息，以便使用主机名命名文件夹
        infrastructureFile, err := repository.Src.Get(util.Join("/", object.Version, repository.Infrastructure))
        deploymentDataReflected, err := template.GetDeploymentInfo(object.Env, servicelistFile, infrastructureFile)
        err = json.Unmarshal([]byte(deploymentDataReflected), &reflectedDeploymentInfo)
        if err != nil {
            log.Sugar().Infof("json unmarshal reflectedDeploymentinfo err of %v, data:%s", err, deploymentDataReflected)
            return
        }
        dataMapReflected := make(map[string]string)
        template.ConstructMap(dataMapReflected, reflectedDeploymentInfo, "")
        //循环取模板
        for _, tmplName := range tmplSlice {
            //对单个模板循环生成配置文件
            for i := 0; i < replicatorNum; i++ {
                path := util.Join("/", object.Version, object.Config, clusterName, repository.Templates, tmplName)
                tmplContent, err := repository.Src.Get(path)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                templateIns, err := template.NewTemplateImpl(repository.Src, object.Env, "0", strconv.Itoa(i), "tmplIns", object.Version, object.Config)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                data, err := templateIns.Fill(tmplContent, tmplName, servicelistFile)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                //创建文件夹和文件,以集群名和节点号区分,为便于运维人员识别，节点号映射到主机名
                hostnameKey := util.Join(".", template.DeploymentInfoKey, strconv.Itoa(i), template.HostNameKey)
                hostname, ok := dataMapReflected[hostnameKey]
                if !ok {
                    log.Sugar().Infof("no such hostname under key: %s", hostnameKey)
                    return
                }
                dirPath := util.Join("/", object.PathOut, clusterName, hostname)
                err = os.MkdirAll(dirPath, os.ModePerm)
                if err != nil {
                    fmt.Println(err)
                    return
                }
                f, err := os.OpenFile(dirPath+"/"+tmplName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
                n, err := f.Write(data)
                if err == nil && n < len(data) {
                    err = io.ErrShortWrite
                    fmt.Println(err)
                    return
                }
                f.Close()
            }
        }
    }
}
