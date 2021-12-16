// Package manage 启动时构建Manage类并注册至grpc服务端
package manage

import (
    "context"
    "regexp"

    "github.com/configcenter/config"
    "github.com/spf13/viper"
)

//请求体Target字段取值范围
const (
    VersionString  = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$` //初始化用于版本号校验的正则表达式
    TemplateString = `templates/`                                          //初始化用于模板文件筛选的正则表达式
    EnvNumString   = `^[0-9][0-9]$`                                        //初始化环境号校验
    LockName       = "/lock_configcenter_ordered"                          //etcd分布式锁名称
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
    //file, err := os.Open(grpcConfigLocation)
    //if err != nil {
    //    return err
    //}
    //binaryFlie, err := ioutil.ReadAll(file)
    //if err != nil {
    //    return err
    //}
    //err = json.Unmarshal(binaryFlie, &manager.grpcInfo)
    //if err != nil {
    //    return err
    //}
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
    //err = file.Close()
    //if err != nil {
    //    return err
    //}
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

func (m *Manager) SetGrpcInfo() {
    m.grpcInfo.Socket = viper.GetString(config.GrpcSocket)
    m.grpcInfo.LockTimeout = viper.GetInt(config.GrpcLockTimeout)
}
