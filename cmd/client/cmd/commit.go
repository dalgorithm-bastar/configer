package cmd

import (
    "context"
    "fmt"
    "io/ioutil"
    "path/filepath"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/service"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
    Use:   "commit",
    Short: "commit a new config version on remote from cache under username,or commit infrastructure",
    Long: `commit command is used for submit the configfile under the username selected.
you can either assign version number in flag --version or leave it as default
attention that you have to put configfile first;otherwise,commit infrastructure is also supported`,
    Run: Post,
}

func init() {
    rootCmd.AddCommand(postCmd)
    postCmd.Flags().StringVarP(&object.Target, "target", "t", "", "select raw or infrastructure to commit")
    postCmd.Flags().StringVarP(&object.PathIn, "pathin", "i", "", "assign input path, only for infrastructure")
    postCmd.Flags().StringVarP(&object.Version, "version", "v", "", "put version number here if you want")
    clusterCmd.MarkFlagRequired("target")
}

func Post(cmd *cobra.Command, args []string) {
    if object.Target != service.TargetRaw && object.Target != service.TargetInfrastructure {
        fmt.Printf("err commit type:%s, please input target within raw or infra", object.Target)
        return
    }
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
    //构建请求结构体
    configReq := pb.CfgReq{
        UserName: object.UserName,
        Target:   object.Target,
        File:     nil,
    }
    //根据提交类型补充请求结构体
    switch object.Target {
    case service.TargetRaw:
        if object.Version != "" {
            configReq.Version = object.Version
        }
    case service.TargetInfrastructure:
        filePath := filepath.Clean(object.PathIn)
        f, err := ioutil.ReadFile(filePath)
        if err != nil {
            panic(err)
        }
        configReq.File = &pb.AnyFile{
            FileName: repository.Infrastructure,
            FileData: f,
        }
    }
    client := pb.NewConfigCenterClient(conn)
    resp, err := client.COMMIT(context.Background(), &configReq)
    if err != nil {
        fmt.Println(err)
        return
    }
    if resp.Status != "ok" {
        fmt.Println(resp.Status)
        return
    }
    fmt.Println(fmt.Sprintf("Commit succeed, possibly new version for raw num %+v", resp.VersionList))
}
