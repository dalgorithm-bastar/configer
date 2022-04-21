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
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "syscall"

    "github.com/configcenter/pkg/generation"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/service"
    "github.com/configcenter/pkg/util"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

var mode string

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
    Use:   "get",
    Short: "Get file from remote or local",
    Run:   Get,
}

func init() {
    rootCmd.AddCommand(clusterCmd)
    clusterCmd.Flags().StringVarP(&object.Target, "target", "t", "", "assign target file Type within raw,cache,infra,version,config")
    clusterCmd.Flags().StringVarP(&object.Type, "type", "y", "", "assign data type within deployment,service,template")
    clusterCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path")
    clusterCmd.Flags().StringVarP(&object.Version, "version", "v", "", "assign a config version")
    clusterCmd.Flags().StringVarP(&object.Env, "env", "e", "", "assign an env num of 2 bits(required)")
    clusterCmd.Flags().StringVarP(&object.Scheme, "scheme", "s", "", "assign config scheme")
    clusterCmd.Flags().StringVarP(&object.Platform, "platform", "l", "", "assign platform")
    clusterCmd.Flags().StringVarP(&object.NodeType, "nodetype", "n", "", "assign config nodetype")
    clusterCmd.Flags().StringVarP(&object.Set, "cluster", "c", "", "assign a cluster name")
    clusterCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    clusterCmd.Flags().StringVarP(&mode, "mode", "m", "", "input \"remote\" or \"local\" to choose creating from remote or local(required)")
    clusterCmd.Flags().StringVarP(&object.IpRange, "ip", "", "", "assign IP range(required)")
    clusterCmd.Flags().StringVarP(&object.PortRange, "port", "p", "", "assign Port range(required)")
    //clusterCmd.MarkFlagRequired("version")
    clusterCmd.MarkFlagRequired("env")
    clusterCmd.MarkFlagRequired("scheme")
    clusterCmd.MarkFlagRequired("pathout")
    clusterCmd.MarkFlagRequired("mode")
    clusterCmd.MarkFlagRequired("ip")
    clusterCmd.MarkFlagRequired("port")
}

func Get(cmd *cobra.Command, args []string) {
    if mode != "local" && mode != "remote" {
        fmt.Println("please input correct arg mode, within \"remote\" or \"local\"")
        return
    }
    object.PathIn = filepath.Clean(object.PathIn)
    if mode == "local" {
        switch object.Target {
        case service.TargetConfig:
            //校验环境号是否合法
            envNumFormat, err := regexp.Compile(service.EnvNumString)
            if err != nil {
                panic(err)
            }
            if !envNumFormat.MatchString(object.Env) {
                panic(fmt.Sprintf("illegal envNum of %s, please input num of 2 bits", object.Env))
            }
            mask := syscall.Umask(0)
            defer syscall.Umask(mask)
            //从本地生成
            s, err := os.Stat(object.PathIn)
            if err != nil || !s.IsDir() {
                fmt.Println("please specify input path to directory")
                return
            }
            //读基础设施文件
            infraPath := filepath.Dir(object.PathIn) + "/" + repository.Infrastructure
            infraData, err := ioutil.ReadFile(infraPath)
            if err != nil {
                panic(err)
            }
            err = repository.NewStorage(context.Background(), repository.DirType, object.PathIn)
            if err != nil {
                panic(err)
            }
            rawData, err := repository.Src.GetbyPrefix(filepath.Base(object.PathIn) + "/" + object.Set)
            if err != nil {
                panic(err)
            }
            IpSlice := strings.Split(object.IpRange, ",")
            portSlice := strings.Split(object.PortRange, ",")
            configData, err := generation.Generate(infraData, rawData, object.Env, IpSlice, portSlice)
            WriteFilesToLocal(configData)
            fmt.Printf("get %s success", object.Target)
        default:
            fmt.Println("on loacl mode there is only target config supported, please checkout target")
        }
        return
    }
    if mode == "remote" {
        //新建客户端
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
        //选择获取对象
        switch object.Target {
        case service.TargetVersion:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
            }
            resp := GetResp(configReq, conn)
            fmt.Printf("%+v", &resp.VersionList)
        case service.TargetRaw:
            if object.Version == "" {
                fmt.Println("lack of required arg: version, on remote mode of target raw")
            }
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
                Version:  object.Version,
                Scheme:   object.Scheme,
                Type:     object.Type,
                Platform: object.Platform,
                NodeType: object.NodeType,
            }
            resp := GetResp(configReq, conn)
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                panic(err)
            }
            WriteFilesToLocal(dataMap)
        case service.TargetCache:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
                Version:  object.UserName,
            }
            resp := GetResp(configReq, conn)
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                panic(err)
            }
            WriteFilesToLocal(dataMap)
        case service.TargetInfrastructure:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
            }
            resp := GetResp(configReq, conn)
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                panic(err)
            }
            WriteFilesToLocal(dataMap)
        case service.TargetConfig:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
                Version:  object.Version,
                Scheme:   object.Scheme,
                EnvNum:   object.Env,
                ArgRange: &pb.ArgRange{
                    TopicIp:   strings.Split(object.IpRange, ","),
                    TopicPort: strings.Split(object.PortRange, ","),
                },
            }
            resp := GetResp(configReq, conn)
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                panic(err)
            }
            WriteFilesToLocal(dataMap)
        }
    }
    fmt.Printf("get %s success", object.Target)
}

func GetResp(req pb.CfgReq, conn *grpc.ClientConn) *pb.CfgResp {
    client := pb.NewConfigCenterClient(conn)
    resp, err := client.GET(context.Background(), &req)
    if err != nil {
        panic(err)
    }
    if resp.Status != "ok" {
        panic(resp.Status)
    }
    return resp
}

func WriteFilesToLocal(fileMap map[string][]byte) {
    //写文件
    //清空同一版本同一方案下的文件夹
    for path, _ := range fileMap {
        pathSli := strings.Split(path, "/")
        if len(pathSli) > 1 {
            pathSli = pathSli[:len(pathSli)-1]
            err := os.RemoveAll(object.PathOut + "/" + pathSli[0])
            if err != nil {
                panic(errors.New(fmt.Sprintf("rmv old dir err: %s,path:%s", err.Error())))
            }
        } else {
            if path == "deployList.json" || path == "topicList.json" {
                err := os.RemoveAll(object.PathOut + "/" + path)
                if err != nil {
                    panic(errors.New(fmt.Sprintf("rmv old dir err: %s,path:%s", err.Error())))
                }
            }
        }
    }
    for path, data := range fileMap {
        pathSli := strings.Split(path, "/")
        if len(pathSli) > 1 {
            pathSli = pathSli[:len(pathSli)-1]
            dirPath := strings.Join(pathSli, "/")
            err := os.MkdirAll(object.PathOut+"/"+dirPath, os.ModePerm)
            if err != nil {
                fmt.Println(err)
                return
            }
        }
        filePath := object.PathOut + "/" + path
        f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
        n, err := f.Write(data)
        if err == nil && n < len(data) {
            err = io.ErrShortWrite
            fmt.Println(err)
            return
        }
        f.Close()
    }
}
