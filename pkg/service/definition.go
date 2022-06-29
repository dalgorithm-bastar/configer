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
    _versionString  = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$` //初始化用于版本号校验的正则表达式
    _templateString = `templates/`                                          //初始化用于模板文件筛选的正则表达式
    EnvNumString    = `^[0-9][0-9]$`                                        //初始化环境号校验
    _lockName       = "/lock_configcenter_ordered"                          //etcd分布式锁名称

    //req中有限制范围的关键字
    _reqTarget = "target"
    _reqType   = "type"
    _reqAction = "action"

    //修改请求体中关键字时，必须同步修改CheckFlag函数
    //Target关键字取值范围
    TargetConfig         = "config"
    TargetInfrastructure = "infra"
    TargetVersion        = "version"
    TargetRaw            = "raw"
    TargetCache          = "cache"

    //Action关键字取值范围
    _actionResetAtRoot           = "resetAtRoot"
    _actionResetAtLeaf           = "resetAtLeaf"
    _actionAddAndReplaceAtLeaf   = "addAndReplaceAtLeaf"
    _actionAddButNoReplaceAtLeaf = "addButNoReplaceAtLeaf"

    //Type关键字取值范围
    _typeDeployment = "deployment"
    _typeService    = "service"
    _typeTemplate   = "template"
)

var (
    manager *Manager
)

type GenSrc struct {
    UserName       string `yaml:"userName"`
    Version        string `yaml:"version"`
    Scheme         string `yaml:"scheme"`
    EnvNum         string `yaml:"envNum"`
    Ip             string `yaml:"ip"`
    CastPort       string `yaml:"castPort"`
    TcpPort        string `yaml:"tcpPort"`
    Infrastructure string `yaml:"infrastructure"`
}

type GenSrcGrp struct {
    GenSrcs []GenSrc `yaml:"genSrcs"`
}

type regExpStruct struct {
    RegExpOfVersion  *regexp.Regexp
    RegExpOfTemplate *regexp.Regexp
    RegExpOfEnvNum   *regexp.Regexp
}

// Manager 调度各模块实现请求
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
    versionFormat, err := regexp.Compile(_versionString)
    if err != nil {
        return err
    }
    templateFormat, err := regexp.Compile(_templateString)
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
    case _reqTarget:
        target := []string{TargetRaw, TargetCache, TargetConfig, TargetVersion, TargetInfrastructure}
        for _, v := range target {
            if v == value {
                return true
            }
        }
    case _reqAction:
        target := []string{_actionAddAndReplaceAtLeaf, _actionAddButNoReplaceAtLeaf, _actionResetAtLeaf, _actionResetAtRoot}
        for _, v := range target {
            if v == value {
                return true
            }
        }
    case _reqType:
        target := []string{_typeTemplate, _typeDeployment, _typeService}
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

func (m *Manager) GetLatestConfigByEnvNum(ctx context.Context, envNumReq *pb.EnvNumReq) (*pb.CfgResp, error) {
    if envNumReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, versionList, file := getLatestConfigByEnvNum(m.ctx, envNumReq)
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
