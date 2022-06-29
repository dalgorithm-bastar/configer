package service

import (
    "context"
    "fmt"
    "strings"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "gopkg.in/yaml.v3"
)

func getLatestConfigByEnvNum(ctx context.Context, req *pb.EnvNumReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
    //校验请求体
    if !manager.regExp.RegExpOfEnvNum.MatchString(req.EnvNum) {
        return fmt.Errorf("envNum format err, input envNum:%s, please input envNum of 2 bit number", req.EnvNum), nil, nil
    }
    //获取环境号记录文件
    envFile, err := repository.Src.Get(repository.EnvFile)
    if err != nil || envFile == nil {
        return fmt.Errorf("get envfile err or envfile not exist yet, err:%s; envfile:%v", err.Error(), envFile), nil, nil
    }
    var genSrcStructs GenSrcGrp
    err = yaml.Unmarshal(envFile, &genSrcStructs)
    if err != nil {
        return fmt.Errorf("envFile on remote do not match format, please checkout"), nil, nil
    }
    //检测目标环境号是否存在
    for _, genSrc := range genSrcStructs.GenSrcs {
        if genSrc.EnvNum != req.EnvNum {
            continue
        }
        cfgReq := &pb.CfgReq{
            UserName: genSrc.UserName,
            Target:   TargetConfig,
            EnvNum:   genSrc.EnvNum,
            Version:  genSrc.Version,
            Scheme:   genSrc.Scheme,
            ArgRange: &pb.ArgRange{
                TopicIp:   strings.Split(genSrc.Ip, ","),
                TopicPort: strings.Split(genSrc.CastPort, ","),
                TcpPort:   strings.Split(genSrc.TcpPort, ","),
            },
        }
        return getConfig(ctx, cfgReq, []byte(genSrc.Infrastructure))
    }
    //目标环境号不存在，返回提示信息
    return fmt.Errorf("envNum:%s has not been generated and recorded on remote yet, please generate first", req.EnvNum), nil, nil
}
