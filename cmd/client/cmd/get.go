package cmd

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "syscall"

    "github.com/configcenter/pkg/generation"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/service"
    "github.com/configcenter/pkg/util"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
    "gopkg.in/yaml.v3"
)

const _separator = string(os.PathSeparator)

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
    clusterCmd.Flags().StringVarP(&object.Set, "cluster(set)", "c", "", "assign a cluster(set) name")
    clusterCmd.Flags().StringVarP(&object.PathOut, "pathout", "o", "", "assign output path(required)")
    clusterCmd.Flags().StringVarP(&mode, "mode", "m", "", "input \"remote\" or \"local\" to choose creating from remote or local(required)")
    clusterCmd.Flags().StringVarP(&object.TopicIpRange, "topicIp", "", "", "assign Topic IP range(required)")
    clusterCmd.Flags().StringVarP(&object.TopicPortRange, "topicPort", "", "", "assign Topic Port range(required)")
    clusterCmd.Flags().StringVarP(&object.TcpPortRange, "tcpPort", "", "", "assign TCP Port range(required)")
    clusterCmd.MarkFlagRequired("mode")
    /*clusterCmd.MarkFlagRequired("version")
      clusterCmd.MarkFlagRequired("env")
      clusterCmd.MarkFlagRequired("scheme")
      clusterCmd.MarkFlagRequired("pathout")
      clusterCmd.MarkFlagRequired("ip")
      clusterCmd.MarkFlagRequired("port")*/
}

func Get(cmd *cobra.Command, args []string) {
    if mode != "local" && mode != "remote" {
        fmt.Println("please input correct arg mode, within \"remote\" or \"local\"")
        return
    }
    object.PathIn = filepath.Clean(object.PathIn)
    object.PathOut = filepath.Clean(object.PathOut)
    if mode == "local" {
        switch object.Target {
        case service.TargetConfig:
            if object.Target == "" || object.Version == "" || object.Env == "" || object.PathIn == "" || object.TcpPortRange == "" ||
                object.Scheme == "" || object.PathOut == "" || object.TopicIpRange == "" || object.TopicPortRange == "" {
                fmt.Println("lack of args within target,version,env,scheme,pathout,pathin,topicIp,topicPort,tcpPort")
                return
            }
            //校验环境号是否合法
            envNumFormat, err := regexp.Compile(service.EnvNumString)
            if err != nil {
                fmt.Println(err)
                return
            }
            if !envNumFormat.MatchString(object.Env) {
                fmt.Printf("illegal envNum of %s, please input num of 2 bits like 00 or 01 .etc", object.Env)
                return
            }
            mask := syscall.Umask(0)
            defer syscall.Umask(mask)
            //从本地生成
            s, err := os.Stat(object.PathIn)
            if err != nil || !s.IsDir() {
                fmt.Printf("open input path err:%s", err.Error())
                return
            }
            permData, err := generatePermFile()
            if err != nil {
                fmt.Println(err)
                return
            }
            permStruct, err := generatePermStruct(permData)
            if err != nil {
                fmt.Println(err)
                return
            }
            //读基础设施文件
            infraPath := filepath.Dir(object.PathIn) + _separator + repository.Infrastructure
            infraData, err := ioutil.ReadFile(infraPath)
            if err != nil {
                fmt.Println(err)
                return
            }
            err = repository.NewStorage(context.Background(), repository.DirType, object.PathIn)
            if err != nil {
                fmt.Println(err)
                return
            }
            rawData, err := repository.Src.GetbyPrefix(filepath.Base(object.PathIn) + "/" + object.Scheme)
            if err != nil {
                fmt.Println(err)
                return
            }
            if rawData == nil || len(rawData) == 0 {
                fmt.Printf("get nil rawdata from input args: %s",
                    filepath.FromSlash(filepath.Base(object.PathIn)+"/"+object.Scheme))
                return
            }
            //校验配置包目录结构
            err = checkInputPackage(rawData)
            if err != nil {
                fmt.Println(err)
                return
            }
            tpcIpSlice, tpcPortSlice, tcpPortSlice := strings.Split(object.TopicIpRange, ","),
                strings.Split(object.TopicPortRange, ","), strings.Split(object.TcpPortRange, ",")
            configData, err := generation.Generate(infraData, rawData, object.Env, tpcIpSlice, tpcPortSlice, tcpPortSlice)
            if err != nil {
                fmt.Println(err)
                return
            }
            err = WriteFilesToLocal(configData, permStruct, filepath.Base(object.PathIn)+"/"+object.Scheme)
            if err != nil {
                fmt.Println(err)
                return
            }
            fmt.Printf("get %s success", object.Target)
            fmt.Println()
        default:
            fmt.Println("on loacl mode there is only target config supported, please checkout arg target")
        }
        return
    }
    if mode == "remote" {
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
        //选择获取对象
        switch object.Target {
        case service.TargetVersion:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
            }
            resp, err := GetResp(configReq, conn)
            if err != nil {
                fmt.Println(err)
                return
            }
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
            resp, err := GetResp(configReq, conn)
            if err != nil {
                fmt.Println(err)
                return
            }
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                fmt.Println(err)
                return
            }
            if _, ok := dataMap[object.Version+"/"+repository.Perms]; !ok {
                fmt.Println("err: get no permission file from remote")
                return
            }
            permStruct, err := generatePermStruct(dataMap[object.Version+"/"+repository.Perms])
            if err != nil {
                fmt.Println(err)
                return
            }
            delete(dataMap, object.Version+"/"+repository.Perms)
            err = WriteFilesToLocal(dataMap, permStruct, "")
            if err != nil {
                fmt.Println(err)
                return
            }
        case service.TargetCache:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
                Version:  object.UserName,
            }
            resp, err := GetResp(configReq, conn)
            if err != nil {
                fmt.Println(err)
                return
            }
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                fmt.Println(err)
                return
            }
            if _, ok := dataMap[object.UserName+"/"+repository.Perms]; !ok {
                fmt.Println("err: get no permission file from remote")
                return
            }
            permStruct, err := generatePermStruct(dataMap[object.UserName+"/"+repository.Perms])
            if err != nil {
                fmt.Println(err)
                return
            }
            delete(dataMap, object.UserName+"/"+repository.Perms)
            err = WriteFilesToLocal(dataMap, permStruct, "")
            if err != nil {
                fmt.Println(err)
                return
            }
        case service.TargetInfrastructure:
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
            }
            resp, err := GetResp(configReq, conn)
            if err != nil {
                fmt.Println(err)
                return
            }
            dataMap := map[string][]byte{resp.File.FileName: resp.File.FileData}
            err = WriteFilesToLocal(dataMap, generation.PermFile{}, "")
            if err != nil {
                fmt.Println(err)
                return
            }
        case service.TargetConfig:
            if object.Target == "" || object.Version == "" || object.Env == "" || object.PathIn == "" || object.TcpPortRange == "" ||
                object.Scheme == "" || object.PathOut == "" || object.TopicIpRange == "" || object.TopicPortRange == "" {
                fmt.Println("lack of args within target,version,env,scheme,pathout,pathin,topicIp,topicPort,tcpPort")
                return
            }
            configReq := pb.CfgReq{
                UserName: object.UserName,
                Target:   object.Target,
                Version:  object.Version,
                Scheme:   object.Scheme,
                EnvNum:   object.Env,
                ArgRange: &pb.ArgRange{
                    TopicIp:   strings.Split(object.TopicIpRange, ","),
                    TopicPort: strings.Split(object.TopicPortRange, ","),
                    TcpPort:   strings.Split(object.TcpPortRange, ","),
                },
            }
            resp, err := GetResp(configReq, conn)
            if err != nil {
                fmt.Println(err)
                return
            }
            dataMap, err := util.DecompressFromStream(resp.File.FileName, resp.File.FileData)
            if err != nil {
                fmt.Println(err)
                return
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
    }
    fmt.Printf("get %s success", object.Target)
}

func GetResp(req pb.CfgReq, conn *grpc.ClientConn) (*pb.CfgResp, error) {
    client := pb.NewConfigCenterClient(conn)
    resp, err := client.GET(context.Background(), &req)
    if err != nil {
        return nil, err
    }
    if resp.Status != "ok" {
        return nil, errors.New(resp.Status)
    }
    return resp, nil
}

// WriteFilesToLocal 根据权限要求写入本地,versionScheme指定生成的方案
func WriteFilesToLocal(fileMap map[string][]byte, permStruct generation.PermFile, versionScheme string) error {
    //消除系统掩码的影响
    umaskNum := syscall.Umask(0)
    defer syscall.Umask(umaskNum)
    //清空同一版本同一方案下的文件夹
    for path, _ := range fileMap {
        pathSli := strings.Split(path, "/")
        if len(pathSli) <= 1 {
            if path == repository.Infrastructure {
                break
            }
            return fmt.Errorf("unexpected filepath: %s", path)
        }
        pathSli = pathSli[:len(pathSli)-1]
        err := os.RemoveAll(object.PathOut + _separator + pathSli[0])
        if err != nil {
            return fmt.Errorf("rmv old dir err: %s,path:%s", err.Error(), path)
        }
        break
    }
    //写文件
    //先写有权限要求的文件
    for _, permUnit := range permStruct.FilePerms {
        if permUnit.IsDir == "1" {
            //转成输出路径
            inputPathSli := strings.SplitN(permUnit.Path, "/", 6)
            if len(inputPathSli) < 6 {
                return fmt.Errorf("err template path:%s, please checkout source pkg dir struct", permUnit.Path)
            }
            //若指定方案生成，跳过其他方案
            if versionScheme != "" && !strings.Contains(permUnit.Path, versionScheme) {
                continue
            }
            for outPath, _ := range fileMap {
                outPathSli := strings.SplitN(outPath, "/", 7)
                if len(outPathSli) < 7 {
                    continue
                }
                //平台名与节点类型均一致，判定为某台主机的配置文件夹
                if inputPathSli[2] == outPathSli[2] && inputPathSli[3] == outPathSli[3] {
                    //组合目标文件夹路径
                    dirPath := util.Join("/", outPathSli[:6]...)
                    dirPath = dirPath + "/" + util.Join("/", inputPathSli[5:]...)
                    //检测目标文件夹是否已存在
                    _, err := os.Stat(filepath.FromSlash(object.PathOut + "/" + dirPath))
                    if err != nil {
                        if !os.IsNotExist(err) {
                            return fmt.Errorf("get an err when checking dirPath:%s, err: %v",
                                filepath.FromSlash(object.PathOut+"/"+dirPath), err)
                        }
                        //文件夹不存在
                        permNum, _ := strconv.ParseInt(permUnit.Perm, 8, 0)
                        err := os.MkdirAll(filepath.FromSlash(object.PathOut+"/"+dirPath), os.FileMode(permNum))
                        if err != nil {
                            return fmt.Errorf("make dir err: %s, dirPath: %s", err.Error(),
                                filepath.FromSlash(object.PathOut+"/"+dirPath))
                        }
                    } else {
                        //文件夹存在，修改权限
                        permNum, _ := strconv.ParseInt(permUnit.Perm, 8, 0)
                        err := os.Chmod(filepath.FromSlash(object.PathOut+"/"+dirPath), os.FileMode(permNum))
                        if err != nil {
                            return fmt.Errorf("chmod dir err: %s, dirPath: %s", err.Error(),
                                filepath.FromSlash(object.PathOut+"/"+dirPath))
                        }
                    }
                }
            }
        } else {
            //处理文件，此时在不同节点的配置文件夹下每份文件一定是不存在的
            //寻找目标文件
            inputPathSli := strings.SplitN(permUnit.Path, "/", 6)
            if len(inputPathSli) < 6 {
                fmt.Printf("err template path:%s, please checkout source pkg dir struct", permUnit.Path)
            }
            for outPath, data := range fileMap {
                outPathSli := strings.SplitN(outPath, "/", 7)
                if len(outPathSli) < 7 {
                    continue
                }
                //平台名与节点类型均一致，且从template目录开始后缀名一致，判定为目标文件
                if inputPathSli[2] == outPathSli[2] && inputPathSli[3] == outPathSli[3] && inputPathSli[5] == outPathSli[6] {
                    err := os.MkdirAll(filepath.Dir(filepath.FromSlash(object.PathOut+"/"+outPath)), os.FileMode(0755))
                    if err != nil {
                        return fmt.Errorf("make dir err: %s, dirPath: %s", err.Error(),
                            filepath.Dir(filepath.FromSlash(object.PathOut+"/"+outPath)))
                    }
                    permNum, _ := strconv.ParseInt(permUnit.Perm, 8, 0)
                    f, err := os.OpenFile(filepath.FromSlash(object.PathOut+"/"+outPath),
                        os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(permNum))
                    if err != nil {
                        return fmt.Errorf("creating file err:%s, filepath:%s", err.Error(), filepath.FromSlash(object.PathOut+"/"+outPath))
                    }
                    n, err := f.Write(data)
                    if err == nil && n < len(data) {
                        err = io.ErrShortWrite
                        return fmt.Errorf("write filedata to disk err of short write")
                    }
                    f.Close()
                }
            }
        }
    }
    for path, data := range fileMap {
        pathSli := strings.Split(path, "/")
        if len(pathSli) > 1 {
            pathSli = pathSli[:len(pathSli)-1]
            dirPath := strings.Join(pathSli, "/")
            err := os.MkdirAll(filepath.FromSlash(object.PathOut+"/"+dirPath), 0755)
            if err != nil {
                return fmt.Errorf("make dir err: %s, current dir: %s", err.Error(), object.PathOut+"/"+dirPath)
            }
        }
        filePath := filepath.FromSlash(object.PathOut + "/" + path)
        f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
        if err != nil {
            return fmt.Errorf("creating file err:%s, filepath:%s", err.Error(), filePath)
        }
        //将json文件格式化输出
        if filepath.Ext(path) == ".json" {
            var outBuf bytes.Buffer
            _ = json.Indent(&outBuf, data, "", "    ")
            data = outBuf.Bytes()
        }
        n, err := f.Write(data)
        if err == nil && n < len(data) {
            err = io.ErrShortWrite
            return fmt.Errorf("write filedata to disk err of short write")
        }
        f.Close()
    }
    return nil
}

func generatePermFile() ([]byte, error) {
    permStruct := generation.PermFile{}
    pathSli := strings.Split(object.PathIn, _separator)
    lenVersion := len(pathSli[len(pathSli)-1])
    idx := len(object.PathIn) - lenVersion + 1
    err := filepath.Walk(object.PathIn, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
            return err
        }
        if strings.Contains(path, _separator+_templateFlag+_separator) {
            pathKey := path[idx-1:]
            fileInfo, statErr := os.Stat(path)
            if statErr != nil {
                return statErr
            }
            permStr, isDir := strconv.FormatUint(uint64(fileInfo.Mode()), 8), "0"
            if fileInfo.IsDir() {
                permStr = permStr[len(permStr)-3:]
                isDir = "1"
            }
            //fmt.Println(fileInfo.Mode())
            permStruct.FilePerms = append(permStruct.FilePerms, generation.PermUnit{
                Path:  filepath.ToSlash(pathKey),
                IsDir: isDir,
                Perm:  permStr,
            })
        }
        return err
    })
    if err != nil {
        fmt.Printf("error when getting permission value walking the path %q: %v\n", object.PathIn, err)
        return nil, err
    }
    permData, err := yaml.Marshal(permStruct)
    return permData, err
}

func generatePermStruct(permData []byte) (generation.PermFile, error) {
    permStruct := generation.PermFile{}
    err := yaml.Unmarshal(permData, &permStruct)
    if err != nil {
        return permStruct, fmt.Errorf("unmarshal perm file err: %v", err)
    }
    return permStruct, nil
}
