//template.go中声明和实现了一组函数
//按照manage的要求注册和填充模板文件、服务信息、公共信息
//此文件内的函数不直接填充模板内容，但做好一切准备工作
//执行填充的相关函数见tmplFuncs.go
package template

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/util"
    "github.com/tidwall/sjson"
)

//target 关键字放在template模块中，从而避免循环引用
const (
    NodeConfig        = "nodeConfig"        //获取单个节点配置
    Templates         = "templates"         //获取某集群下的所有模板
    Services          = "services"          //获取某个集群下的服务清单
    Manipulations     = "manipulations"     //获取某集群下工作流程
    DeploymentInfo    = "deploymentInfo"    //获取某集群下全部部署信息
    Infrastructure    = "infrastructure"    //获取某版本下所有公共信息文件
    Versions          = "versions"          //获取所有版本号
    ConfigScheme      = "configScheme"      //获取某版本下所有配置方案名
    Clusters          = "clusters"          //获取某环境下所有集群名称
    PartlyOnline      = "partlyOnline"      //上传一份本地模板，利用服务端信息生成后返回
    CtlFindFlag       = "ctlFind"           //用于单条查询
    ReplicatorNumKey  = "replicator_number" //在服务清单中指明实例数目的键
    DeploymentInfoKey = "deployment_info"   //在服务清单中指明部署信息的键
    HostNameKey       = "hostname"          //在服务清单中指明主机名称的键
    NodeIdKey         = "NODE_ID"           //在服务清单中指明集群节点号的键
    EnvNumKey         = "envNum"            //在模板中指定环境号的键
)

type TemplateImpl struct {
    funcMap      map[string]interface{} //公共信息表，服务信息表和模板函数表
    allTemplates *template.Template     //公有模板，在该模板上注册所有函数，并根据需要注册待填充模板
}

// NewCtlFindTemplate 用于新建单条查找请求对应的填充模板
func NewCtlFindTemplate(tmplInstanceName string) (*TemplateImpl, error) {
    templateIns := new(TemplateImpl)
    templateIns.allTemplates = template.New(tmplInstanceName)
    //模板中的函数名是map的key即可，添加新函数务必同时添加注释！！！
    templateIns.funcMap = map[string]interface{}{
        //用于命令行查询，涉及公共信息的不执行映射
        "CtlFind": CtlFind,
    }
    //注册函数列表
    templateIns.allTemplates.Funcs(templateIns.funcMap)
    //解析并生效
    tmp, err := templateIns.allTemplates.Parse(tmplInstanceName)
    if err != nil {
        return nil, err
    }
    templateIns.allTemplates = tmp
    return templateIns, nil
}

// NewTemplateImpl 创建对象，除模板实例名称外，其他均为填充模板时需要用到的快照信息
func NewTemplateImpl(source repository.Storage, envNum, globalId, localId, tmplInstanceName, version, conf string) (*TemplateImpl, error) {
    templateIns := new(TemplateImpl)
    templateIns.allTemplates = template.New(tmplInstanceName)
    //获取公共信息文件，仅对GetInfo函数有意义
    prefixPublic := util.Join("/", version, repository.Infrastructure)
    binaryData, err := source.Get(prefixPublic) //从数据源取值
    if err != nil {
        log.Sugar().Errorf("get Infrastructure from repository err of %v when init tmpl, under path %s", err, prefixPublic)
        return nil, err
    }
    if binaryData == nil {
        log.Sugar().Infof("get nil Infrastructure under path %s when init tmpl", prefixPublic)
        return nil, errors.New(fmt.Sprintf("no infrastructureData under path %s, please checkout in etcd or compressedfile", prefixPublic))
    }
    //模板中的函数名是map的key即可，添加新函数务必同时添加注释！！！
    //此处使用闭包声明模板函数，在模板执行时使用的信息为传参时的快照信息
    templateIns.funcMap = map[string]interface{}{
        //用于查询当前环境下的服务信息，涉及公共信息的执行映射
        "GetInfo": func(mode, clusterObject, service string) (string, error) {
            if mode != "normal" && mode != "slice" {
                return "", errors.New("mode out of range, please assign \"normal\" or \"slice\"")
            }
            defaultIndex := true
            if mode == "normal" {
                defaultIndex = false
            }
            return func(src repository.Storage, en string, binData []byte, dftIdx bool, glbId, lcId, vr, cf, cl, sr string) (string, error) {
                return baseGet(src, binData, dftIdx, en, glbId, lcId, vr, cf, cl, sr)
            }(source, envNum, binaryData, defaultIndex, globalId, localId, version, conf, clusterObject, service)
        },
        //用于命令行查询，涉及公共信息的不执行映射
        "CtlFind": CtlFind,
        //用于获取全局服务信息，存在潜在的风险，不推荐使用
        "UnsafeGetInfo": func(mode, anyVersion, anyConf, clusterObject, service string) (string, error) {
            if mode != "normal" && mode != "slice" {
                return "", errors.New("mode out of range, please assign \"normal\" or \"slice\"")
            }
            defaultIndex := true
            if mode == "normal" {
                defaultIndex = false
            }
            return func(src repository.Storage, dftIdx bool, en, glbId, lcId, vr, acf, cl, sr string) (string, error) {
                return baseGet(src, nil, dftIdx, en, glbId, lcId, vr, acf, cl, sr)
            }(source, defaultIndex, envNum, globalId, localId, anyVersion, anyConf, clusterObject, service)
        },
        "GetNodeIdInfo": func(NodeId, clusterObject, service string) (string, error) {
            return func(src repository.Storage, binData []byte, dftIdx bool, en, glbId, lcId, vr, cf, cl, sr string) (string, error) {
                return baseGet(src, binData, dftIdx, en, glbId, lcId, vr, cf, cl, sr)
            }(source, binaryData, true, envNum, globalId, NodeId, version, conf, clusterObject, service)
        },
        "ParseFloat": ParseFloat,
        "FmtFloat":   FmtFloat64,
        "Itoa":       Itoa,
        "Atoi":       Atoi,
        "Add":        Add,
        "Mine":       Mine,
    }
    //注册函数列表
    templateIns.allTemplates.Funcs(templateIns.funcMap)
    //解析并生效
    tmp, err := templateIns.allTemplates.Parse(tmplInstanceName)
    if err != nil {
        return nil, err
    }
    templateIns.allTemplates = tmp
    return templateIns, nil
}

// AddTmpl 按照输入的模板名称和内容向公共模板中添加新模板
func (t *TemplateImpl) addTmpl(tmplContent []byte, tmplName string) error {
    tmpl := t.allTemplates.New(tmplName)
    tmpl, err := tmpl.Parse(string(tmplContent))
    if err != nil {
        return err
    }
    return nil
}

// Fill 按照输入的模板名称和内容进行填充，返回填充后的结果
func (t *TemplateImpl) Fill(tmplContent []byte, tmplName string, srcContent []byte) ([]byte, error) {
    if tmplContent == nil || srcContent == nil {
        msg := fmt.Sprintf("nil tmplContent or srcContent input when tmpl filling, tmplName is %s", tmplName)
        log.Sugar().Info(msg)
        return nil, errors.New(msg)
    }
    var srcData interface{}
    err := json.Unmarshal(srcContent, &srcData)
    if err != nil {
        return nil, err
    }
    err = t.addTmpl(tmplContent, tmplName)
    if err != nil {
        return nil, err
    }
    var data bytes.Buffer
    err = t.allTemplates.ExecuteTemplate(&data, tmplName, srcData)
    if err != nil {
        return nil, err
    }
    return data.Bytes(), nil
}

// GetDeploymentInfo 接收两份文件，返回一份新的部署信息文件
func GetDeploymentInfo(envNum string, serviceData, infrastructureData []byte) (string, error) {
    if serviceData == nil || infrastructureData == nil {
        return "", errors.New("nil input when filling service data")
    }
    infrastructureData, err := sjson.SetBytes(infrastructureData, EnvNumKey, envNum)
    if err != nil {
        return "", errors.New("err when set envNum into infra: " + err.Error())
    }
    res, err := fillSrvbyInfra(serviceData, infrastructureData)
    if err != nil {
        return "", err
    }
    return string(res), nil
}
