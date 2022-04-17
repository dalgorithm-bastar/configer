package cmd

import (
    "context"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/util"
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
    //检测是否为文件夹
    s, err := os.Stat(object.PathIn)
    if err != nil || !s.IsDir() {
        fmt.Println("please specify input path to directory")
        return
    }
    fileMap := make(map[string][]byte)
    object.PathIn = filepath.Clean(object.PathIn)
    pathSli := strings.Split(object.PathIn, "/")
    lenVersion := len(pathSli[len(pathSli)-1])
    idx := len(object.PathIn) - lenVersion + 1
    err = filepath.Walk(object.PathIn, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            pathKey := path[idx-1:]
            data, err := ioutil.ReadFile(path)
            if err != nil {
                return err
            }
            fileMap[pathKey] = data
        }
        return nil
    })
    compressedFileData, err := util.CompressToStream("cfgpkg.tar.gz", fileMap)
    if err != nil {
        panic(err)
    }
    //新建客户端
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
    resp, err := client.PUT(context.Background(), &pb.CfgReq{
        UserName: object.UserName,
        File: &pb.AnyFile{
            FileName: "cfgpkg.tar.gz",
            FileData: compressedFileData,
        },
    })
    if err != nil {
        panic(err)
    }
    if resp.Status != "ok" {
        panic(resp.Status)
    }
    fmt.Println(fmt.Sprintf("Put cfgpkg to remote succeed"))
}
