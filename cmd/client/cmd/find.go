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
    "io/ioutil"
    "os"

    manage "github.com/configcenter/pkg/manager"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/template"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

const (
    FindFileName = "CtlFindFile.txt"
)

var GrpcInfo manage.GrpcInfoStruct

// findCmd represents the find command
var findCmd = &cobra.Command{
    Use:   "find",
    Short: "find particular info from target",
    Long: `find particular info by making phrases as go template format, result will be presented on cmdline
the params of func CtlFind is (Target, version, env, cluster, service)
please input "" to the param not used`,
    Example: `$ ./cfgtool find --phrase "{{CtlFind(\"servicelist\" \"1.0.0\" \"00\" \"EzEI.set0\" \"MUDP_IP\")}}"" --pathout /home/someuser
$ ./cfgtool find -p "{{CtlFind(\"infrastructure\" "\"1.0.0\" \"\" \"\" \"hostName1_IP\")}}" -o /home/someuser`,
    Run: Find,
}

func init() {
    rootCmd.AddCommand(findCmd)
    findCmd.Flags().StringVarP(&object.Phrase, "phrase", "p", "", "phrase as go template format(required)")
    findCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    findCmd.MarkFlagRequired("phrase")
    findCmd.MarkFlagRequired("pathout")
}

func Find(cmd *cobra.Command, args []string) {
    //构建请求结构体
    configReq := pb.CfgReq{
        UserName: object.UserName,
        Target:   []string{template.CtlFindFlag},
        File: &pb.CompressedFile{
            FileName: FindFileName,
            FileData: []byte(object.Phrase),
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
    err = os.MkdirAll(object.PathOut, os.ModePerm)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(string(resp.File.FileData))
    f, err := os.OpenFile(object.PathOut+"/"+FindFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
}

func GetGrpcClient() error {
    //读取grpc配置文件
    var file *os.File
    file, err := os.Open(grpcConfigLocationInProject)
    if err != nil {
        file, err = os.Open(grpcConfigLocationInExe)
        if err != nil {
            fmt.Println(fmt.Sprintf("Can not open grpcConfigFile in %s or %s", grpcConfigLocationInProject, grpcConfigLocationInExe))
            return err
        }
    }
    binaryFlie, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println("Read grpcConfigFile err")
        return err
    }
    err = json.Unmarshal(binaryFlie, &GrpcInfo)
    if err != nil {
        fmt.Println("Json Unmarshal grpcConfigFile err")
        return err
    }
    return nil
}
