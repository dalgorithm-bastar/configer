package cmd

import (
    "context"
    "fmt"
    "path/filepath"
    "strings"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/util"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

// getbyenvCmd represents the getbyenv command
var getbyenvCmd = &cobra.Command{
    Use:   "getbyenv",
    Short: "get latest config file by envnum",
    Long:  `this function is used for any application or user who wants to get newest config of specific envnum on remote`,
    Run:   getConfigByEnvNum,
}

func init() {
    rootCmd.AddCommand(getbyenvCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // getbyenvCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // getbyenvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    getbyenvCmd.Flags().StringVarP(&object.Env, "env", "e", "", "assign an env num of 2 bits(required)")
    getbyenvCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    getbyenvCmd.MarkFlagRequired("env")
    getbyenvCmd.MarkFlagRequired("pathout")
}

func getConfigByEnvNum(cmd *cobra.Command, args []string) {
    object.PathOut = filepath.Clean(object.PathOut)
    //新建客户端
    //读取grpc配置信息
    err := GetGrpcClient()
    if err != nil {
        fmt.Println(err)
        return
    }
    //新建grpc客户端
    conn, err := grpc.Dial(GrpcInfo.Socket, grpc.WithInsecure())
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()
    client := pb.NewConfigCenterClient(conn)
    resp, err := client.GetLatestConfigByEnvNum(context.Background(), &pb.EnvNumReq{
        EnvNum: object.Env,
    })
    if err != nil {
        fmt.Println(err)
        return
    }
    if resp.Status != "ok" {
        fmt.Println(resp.Status)
        return
    }

    dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
    if err != nil {
        fmt.Println(err)
        return
    }
    //从环境号获取时，要先从结果获取版本和方案
    for path, _ := range dataMap {
        pathSlice := strings.Split(path, "/")
        if len(pathSlice) < 3 {
            continue
        }
        versionSchemeSlice := strings.Split(pathSlice[0], "_")
        if len(versionSchemeSlice) != 2 {
            fmt.Printf("got err file name:%s", path)
            return
        }
        object.Version, object.Scheme = versionSchemeSlice[0], versionSchemeSlice[1]
        break
    }
    if _, ok := dataMap[object.Version+"/"+repository.Perms]; !ok {
        fmt.Printf("err: get no permission file from remote, file path:%s", object.Version+"/"+repository.Perms)
        return
    }
    permStruct, err := generatePermStruct(dataMap[object.Version+"/"+repository.Perms])
    if err != nil {
        fmt.Println(err)
        return
    }
    delete(dataMap, object.Version+"/"+repository.Perms)
    err = WriteFilesToLocal(dataMap, permStruct, object.Version+"/"+object.Scheme)
    if err != nil {
        fmt.Println(err)
        return
    }
}
