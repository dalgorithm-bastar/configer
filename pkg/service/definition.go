// Package service启动时构建Manage类并注册至grpc服务端
package service

import (
    "context"
    "regexp"

    "github.com/configcenter/config"
    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/pb"
    "github.com/spf13/viper"
)

//请求体中字段取值范围
const (
    VersionString  = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$` //初始化用于版本号校验的正则表达式
    TemplateString = `templates/`                                          //初始化用于模板文件筛选的正则表达式
    EnvNumString   = `^[0-9][0-9]$`                                        //初始化环境号校验
    LockName       = "/lock_configcenter_ordered"                          //etcd分布式锁名称

    //req中有限制范围的关键字
    ReqTarget = "target"
    ReqType   = "type"
    ReqAction = "action"

    //修改请求体中关键字时，必须同步修改CheckFlag函数
    //Target关键字取值范围
    TargetConfig         = "config"
    TargetInfrastructure = "infra"
    TargetVersion        = "version"
    TargetRaw            = "raw"
    TargetCache          = "cache"

    //Action关键字取值范围
    ActionResetAtRoot           = "resetAtRoot"
    ActionResetAtLeaf           = "resetAtLeaf"
    ActionAddAndReplaceAtLeaf   = "addAndReplaceAtLeaf"
    ActionAddButNoReplaceAtLeaf = "addButNoReplaceAtLeaf"

    //Type关键字取值范围
    TypeDeployment = "deployment"
    TypeService    = "service"
    TypeTemplate   = "template"
)

var (
    manager *Manager
)

type regExpStruct struct {
    RegExpOfVersion  *regexp.Regexp
    RegExpOfTemplate *regexp.Regexp
    RegExpOfEnvNum   *regexp.Regexp
}

//调度各模块实现请求
type Manager struct {
    ctx      context.Context
    grpcInfo GrpcInfoStruct
    regExp   regExpStruct
}

type GrpcInfoStruct struct {
    Socket      string `json:"socket"`
    LockTimeout int    `json:"locktimeout"`
}

// NewManager 读取grpc配置并初始化manager实例
func NewManager(ctxIn context.Context) error {
    manager = new(Manager)
    manager.SetGrpcInfo()
    versionFormat, err := regexp.Compile(VersionString)
    if err != nil {
        return err
    }
    templateFormat, err := regexp.Compile(TemplateString)
    if err != nil {
        return err
    }
    envNumFormat, err := regexp.Compile(EnvNumString)
    if err != nil {
        return err
    }
    manager.regExp = regExpStruct{
        RegExpOfVersion:  versionFormat,
        RegExpOfTemplate: templateFormat,
        RegExpOfEnvNum:   envNumFormat,
    }
    manager.ctx = ctxIn
    return nil
}

// GetManager 获取manager对象，一般用于测试
func GetManager() *Manager {
    return manager
}

// GetGrpcInfo 获取manager对象的grpc信息
func GetGrpcInfo() *GrpcInfoStruct {
    return &manager.grpcInfo
}

func CheckFlag(flagName, value string) bool {
    switch flagName {
    case ReqTarget:
        target := []string{TargetRaw, TargetCache, TargetConfig, TargetVersion, TargetInfrastructure}
        for _, v := range target {
            if v == value {
                return true
            }
        }
    case ReqAction:
        target := []string{ActionAddAndReplaceAtLeaf, ActionAddButNoReplaceAtLeaf, ActionResetAtLeaf, ActionResetAtRoot}
        for _, v := range target {
            if v == value {
                return true
            }
        }
    case ReqType:
        target := []string{TypeTemplate, TypeDeployment, TypeService}
        for _, v := range target {
            if v == value {
                return true
            }
        }
    }
    return false
}

func (m *Manager) SetGrpcInfo() {
    m.grpcInfo.Socket = viper.GetString(config.GrpcSocket)
    m.grpcInfo.LockTimeout = viper.GetInt(config.GrpcLockTimeout)
}

func (m *Manager) GET(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
    go log.Sugar().Infof("GET CfgReq Recieved: %+v", CfgReq)
    if CfgReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, versionList, file := Get(m.ctx, CfgReq)
    var status string
    if err != nil {
        status = err.Error()
    } else {
        status = "ok"
    }
    CfgResp := &pb.CfgResp{
        Status:      status,
        VersionList: versionList,
        File:        file,
    }
    log.Sugar().Infof("GET CfgResp Created: %+v", CfgResp)
    return CfgResp, nil
}

func (m *Manager) COMMIT(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
    go log.Sugar().Infof("COMMIT CfgReq Recieved: %+v", CfgReq)
    if CfgReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, versionList, file := commit(m.ctx, CfgReq)
    var status string
    if err != nil {
        status = err.Error()
    } else {
        status = "ok"
    }
    CfgResp := &pb.CfgResp{
        Status:      status,
        VersionList: versionList,
        File:        file,
    }
    log.Sugar().Infof("COMMIT CfgResp Created: %+v", CfgResp)
    return CfgResp, nil
}

func (m *Manager) DELETE(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
    go log.Sugar().Infof("DELETE CfgReq Recieved: %+v", CfgReq)
    if CfgReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, versionList, file := DeleteInManager(m.ctx, CfgReq)
    var status string
    if err != nil {
        status = err.Error()
    } else {
        status = "ok"
    }
    CfgResp := &pb.CfgResp{
        Status:      status,
        VersionList: versionList,
        File:        file,
    }
    log.Sugar().Infof("DELETE CfgResp Created: %+v", CfgResp)
    return CfgResp, nil
}

func (m *Manager) PUT(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
    //go log.Sugar().Infof("PUT CfgReq Recieved: %+v", CfgReq)
    if CfgReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, versionList, file := put(m.ctx, CfgReq)
    var status string
    if err != nil {
        status = err.Error()
    } else {
        status = "ok"
    }
    CfgResp := &pb.CfgResp{
        Status:      status,
        VersionList: versionList,
        File:        file,
    }
    log.Sugar().Infof("PUT CfgResp Created: %+v", CfgResp)
    return CfgResp, nil
}
