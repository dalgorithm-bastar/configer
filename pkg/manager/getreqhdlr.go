package manage

import (
    "context"
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/template"
    "github.com/configcenter/pkg/util"

    "github.com/configcenter/pkg/pb"
)

// Get 用于从服务端获取信息
func Get(ctx context.Context, req *pb.CfgReq) (error, []string, string, []byte) {
    //无用户名时返回错误
    if req.UserName == "" {
        return errors.New("no username when getting"), nil, "", nil
    }
    //target字段为空时返回错误
    if req.Target == nil || req.Target[0] == "" {
        return errors.New("no target specified when getting"), nil, "", nil
    }
    //确定请求结构体内容是否与target匹配
    err := checkData(req.CfgVersions, req.Target[0])
    if err != nil {
        log.Sugar().Infof("checkData failed in GET req, get err %v, with target %s and data %+v", err, req.Target[0], req.CfgVersions)
        return err, nil, "", nil
    }
    //根据约定常量确定要获取的资源类型
    switch req.Target[0] {
    case template.NodeConfig:
        return getNodeConfig(req)
    case template.Templates:
        return getTemplates(req)
    case template.Services:
        return getServices(req)
    case template.Manipulations:
        return getManipulations(req)
    case template.Infrastructure:
        return getInfrastructure(req)
    case template.DeploymentInfo:
        return getDeploymentInfo(req)
    case template.Versions:
        return getVersions(req)
    case template.Environments:
        return getEnvironments(req)
    case template.Clusters:
        return getClusters(req)
    case template.PartlyOnline:
        return getNodeConfigPartlyOnline(req)
    case template.CtlFindFlag:
        return getKeyValueResult(req)
    //默认返回错误
    default:
        err := fmt.Sprintf("Target of %s can not be recognized in GET req", req.Target[0])
        log.Sugar().Infof(err)
        return errors.New(err), nil, "", nil
    }
}

//判断请求体是否符合约定，不符合则提前返回错误，目前不支持网状结构查询
func checkData(data []*pb.CfgVersion, tag string) error {
    switch tag {
    case template.NodeConfig, template.PartlyOnline:
        if len(data) != 1 {
            return errors.New("args error: too many cfgVersions or missing cfgVersions")
        }
        if len(data[0].Envs) != 1 {
            return errors.New("args error: too many envs or missing envs")
        }
        if len(data[0].Envs[0].Clusters) != 1 {
            return errors.New("args error: too many clusters or missing clusters")
        }
        if len(data[0].Envs[0].Clusters[0].Nodes) != 1 {
            return errors.New("args error: too many Nodes or missing Nodes")
        }
        //检查Id是否合法
        _, err := strconv.Atoi(data[0].Envs[0].Clusters[0].Nodes[0].GlobalId)
        if err != nil {
            return err
        }
        _, err = strconv.Atoi(data[0].Envs[0].Clusters[0].Nodes[0].LocalId)
        if err != nil {
            return err
        }
    case template.Templates, template.Services, template.Manipulations, template.DeploymentInfo:
        if len(data) != 1 {
            return errors.New("args error: too many cfgVersions or missing cfgVersions")
        }
        if len(data[0].Envs) != 1 {
            return errors.New("args error: too many envs or missing envs")
        }
        if len(data[0].Envs[0].Clusters) != 1 {
            return errors.New("args error: too many clusters or missing clusters")
        }
    case template.Clusters:
        if len(data) != 1 {
            return errors.New("args error: too many cfgVersions or missing cfgVersions")
        }
        if len(data[0].Envs) != 1 {
            return errors.New("args error: too many envs or missing envs")
        }
    case template.Infrastructure, template.Environments:
        if len(data) != 1 {
            return errors.New("args error: too many cfgVersions or missing cfgVersions")
        }
    case template.Versions, template.CtlFindFlag:
        //无要求
    default:
        //默认不做处理
    }
    return nil
}

func getNodeConfig(req *pb.CfgReq) (error, []string, string, []byte) {
    //获取基本参数
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    tmplobj := req.CfgVersions[0].Envs[0].Clusters[0].Nodes[0].Template
    path := util.Join("/", version, env, cluster, repository.Templates, tmplobj)
    tmplContent, err := repository.Src.Get(path)
    if err != nil {
        log.Sugar().Errorf("get tmpl from repository err of %v, under path %s", err, path)
        return err, nil, "", nil
    }
    if tmplContent == nil {
        log.Sugar().Infof("get nil template under path %s", path)
        return errors.New(fmt.Sprintf("No Template under save path %s", path)), nil, "", nil
    }
    return createNodeConfig(req, tmplobj, tmplContent)
}

func getTemplates(req *pb.CfgReq) (error, []string, string, []byte) {
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    path := util.GetPrefix(util.Join("/", version, env, cluster, repository.Templates))
    fileMap, err := repository.Src.GetbyPrefix(path)
    if err != nil {
        log.Sugar().Errorf("get tmpl by prefix from repository err of %v, under prefix %s", err, path)
        return err, nil, "", nil
    }
    if fileMap == nil {
        log.Sugar().Infof("get nil template under prefix %s", path)
        return errors.New(fmt.Sprintf("No Templates under save path %s", path)), nil, "", nil
    }
    compressedFileName := util.Join("_", version, env, cluster, template.Templates) + ".tar.gz"
    if req.File != nil && req.File.FileName != "" {
        compressedFileName = req.File.FileName
    }
    compressedFileData, err := util.CompressToStream(compressedFileName, fileMap)
    if err != nil {
        log.Sugar().Errorf("compress templates err of %v, filename %s, filemap %+v", err, compressedFileName, fileMap)
        return err, nil, "", nil
    }
    return nil, nil, compressedFileName, compressedFileData
}

func getServices(req *pb.CfgReq) (error, []string, string, []byte) {
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    path := util.Join("/", version, env, cluster, repository.ServiceList)
    //单文件用Get接口
    fileData, err := repository.Src.Get(path)
    if err != nil {
        log.Sugar().Errorf("get serviceList from repository err of %v, under path %s", err, path)
        return err, nil, "", nil
    }
    if fileData == nil {
        log.Sugar().Infof("get nil serviceList under path %s", path)
        return errors.New(fmt.Sprintf("No servicelist under save path %s", path)), nil, "", nil
    }
    compressedFileName := util.Join("_", version, env, cluster, repository.ServiceList) + ".tar.gz"
    if req.File != nil && req.File.FileName != "" {
        compressedFileName = req.File.FileName
    }
    compressedFileData, err := util.CompressToStream(compressedFileName,
        map[string][]byte{
            path: fileData,
        })
    if err != nil {
        log.Sugar().Errorf("compress serviceList err of %v, filename %s, filemap %s:%s", err, compressedFileName, path, fileData)
        return err, nil, "", nil
    }
    return nil, nil, compressedFileName, compressedFileData
}

func getManipulations(req *pb.CfgReq) (error, []string, string, []byte) {
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    path := util.GetPrefix(util.Join("/", version, env, cluster, repository.Manipulations))
    fileMap, err := repository.Src.GetbyPrefix(path)
    if err != nil {
        log.Sugar().Errorf("get manipulations by prefix from repository err of %v, under prefix %s", err, path)
        return err, nil, "", nil
    }
    if fileMap == nil {
        log.Sugar().Infof("get nil manipulations under prefix %s", path)
        return errors.New(fmt.Sprintf("No manipulations under save path %s", path)), nil, "", nil
    }
    compressedFileName := util.Join("_", version, env, cluster, repository.Manipulations) + ".tar.gz"
    if req.File != nil && req.File.FileName != "" {
        compressedFileName = req.File.FileName
    }
    compressedFileData, err := util.CompressToStream(compressedFileName, fileMap)
    if err != nil {
        log.Sugar().Errorf("compress manipulations err of %v, filename %s, filemap %+v", err, compressedFileName, fileMap)
        return err, nil, "", nil
    }
    return nil, nil, compressedFileName, compressedFileData
}

func getInfrastructure(req *pb.CfgReq) (error, []string, string, []byte) {
    version := req.CfgVersions[0].Version
    path := util.Join("/", version, repository.Infrastructure)
    fileData, err := repository.Src.Get(path)
    if err != nil {
        log.Sugar().Errorf("get Infrastructure from repository err of %v, under path %s", err, path)
        return err, nil, "", nil
    }
    if fileData == nil {
        log.Sugar().Infof("get nil Infrastructure under path %s", path)
        return errors.New(fmt.Sprintf("No infrastructure under save path %s", path)), nil, "", nil
    }
    compressedFileName := util.Join("_", version, repository.Infrastructure) + ".tar.gz"
    if req.File != nil && req.File.FileName != "" {
        compressedFileName = req.File.FileName
    }
    compressedFileData, err := util.CompressToStream(compressedFileName,
        map[string][]byte{
            path: fileData,
        })
    if err != nil {
        log.Sugar().Errorf("compress serviceList err of %v, filename %s, filemap %s:%s", err, compressedFileName, path, fileData)
        return err, nil, "", nil
    }
    return nil, nil, compressedFileName, compressedFileData
}

func getVersions(req *pb.CfgReq) (error, []string, string, []byte) {
    versions, err := repository.Src.Get(repository.Versions)
    if err != nil {
        log.Sugar().Errorf("get versions from repository err of %v, under path %s", err, repository.Versions)
        return err, nil, "", nil
    }
    if versions == nil {
        log.Sugar().Infof("get nil version under path %s", repository.Versions)
        return errors.New("No version saved in configcenter yet"), nil, "", nil
    }
    return nil, strings.Split(string(versions), ","), "", nil
}

func getEnvironments(req *pb.CfgReq) (error, []string, string, []byte) {
    key := util.Join("/", req.CfgVersions[0].Version, repository.Envs)
    envs, err := repository.Src.Get(key)
    if err != nil {
        log.Sugar().Errorf("get envs from repository err of %v, under path %s", err, key)
        return err, nil, "", nil
    }
    if envs == nil {
        log.Sugar().Infof("get nil envs under path %s", key)
        return errors.New(fmt.Sprintf("No envs under save path %s", key)), nil, "", nil
    }
    return nil, strings.Split(string(envs), ","), "", nil
}

func getClusters(req *pb.CfgReq) (error, []string, string, []byte) {
    key := util.Join("/", req.CfgVersions[0].Version, req.CfgVersions[0].Envs[0].Num, repository.Clusters)
    clusters, err := repository.Src.Get(key)
    if err != nil {
        log.Sugar().Errorf("get clusters from repository err of %v, under path %s", err, key)
        return err, nil, "", nil
    }
    if clusters == nil {
        log.Sugar().Infof("get nil clusters under path %s", key)
        return errors.New(fmt.Sprintf("No clusters under save path %s", key)), nil, "", nil
    }
    return nil, strings.Split(string(clusters), ","), "", nil
}

func getKeyValueResult(req *pb.CfgReq) (error, []string, string, []byte) {
    //获取基本参数
    if req.File == nil || req.File.FileData == nil {
        log.Sugar().Infof("Get nil file or fileContent in CtlFind, file:%+v", req.File)
        return errors.New("nil File or FileContent transformed, please check your request"), nil, "", nil
    }
    tmplObj := req.File.FileName
    if tmplObj == "" {
        tmplObj = "ctlFindFile.txt"
    }
    tmplContent := req.File.FileData
    templateIns, err := template.NewCtlFindTemplate("tmplIns")
    if err != nil {
        log.Sugar().Infof("init tmpl err in CtlFind of %v", err)
        return err, nil, "", nil
    }
    data, err := templateIns.Fill(tmplContent, tmplObj)
    if err != nil {
        log.Sugar().Infof("Fill tmpl err in CtlFind of %v, tmplname %s, tmplContent %s", err, tmplObj, tmplContent)
        return err, nil, "", nil
    }
    return nil, nil, tmplObj, data
}

func getNodeConfigPartlyOnline(req *pb.CfgReq) (error, []string, string, []byte) {
    //完全从服务端获取
    if req.File == nil || req.File.FileData == nil || len(req.File.FileData) == 0 {
        return getNodeConfig(req)
    }
    //使用客户端模板获取
    tmplobj := req.File.FileName
    tmplContent := req.File.FileData
    return createNodeConfig(req, tmplobj, tmplContent)
}

//传入待填充模板的名称和内容
func createNodeConfig(req *pb.CfgReq, tmplObj string, tmplContent []byte) (error, []string, string, []byte) {
    //获取基本参数
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    //tmplobj := req.File.FileName
    //创建模板实例
    templateIns, err := template.NewTemplateImpl(repository.Src,
        req.CfgVersions[0].Envs[0].Clusters[0].Nodes[0].GlobalId,
        req.CfgVersions[0].Envs[0].Clusters[0].Nodes[0].LocalId,
        "tmplIns",
        req.CfgVersions[0].Version,
        req.CfgVersions[0].Envs[0].Num)
    if err != nil {
        log.Sugar().Errorf("Init tmpl err in creating nodeconfig, err:%v", err)
        return err, nil, "", nil
    }
    //将要填充的模板注册到模板实例中
    data, err := templateIns.Fill(tmplContent, tmplObj)
    if err != nil {
        log.Sugar().Infof("Fill tmpl err in creating nodeconfig of %v, tmplname %s, tmplContent %s", err, tmplObj, tmplContent)
        return err, nil, "", nil
    }
    fileName := util.Join("_", version, env, cluster, req.CfgVersions[0].Envs[0].Clusters[0].Nodes[0].LocalId, tmplObj)
    if req.File != nil && req.File.FileName != "" {
        fileName = req.File.FileName
    }
    return nil, nil, fileName, data
}

//该函数在json之间处理数据，不需要用到模板填充
func getDeploymentInfo(req *pb.CfgReq) (error, []string, string, []byte) {
    version := req.CfgVersions[0].Version
    env := req.CfgVersions[0].Envs[0].Num
    cluster := req.CfgVersions[0].Envs[0].Clusters[0].ClusterName
    servicePath := util.Join("/", version, env, cluster, repository.ServiceList)
    serviceData, err := repository.Src.Get(servicePath)
    if err != nil {
        log.Sugar().Errorf("get serviceList from repository err of %v, under path %s", err, servicePath)
        return err, nil, "", nil
    }
    if serviceData == nil {
        log.Sugar().Infof("get nil serviceList under path %s", servicePath)
        return errors.New(fmt.Sprintf("No servicelist under path %s", servicePath)), nil, "", nil
    }
    infrastructurePath := util.Join("/", version, repository.Infrastructure)
    infrastructureData, err := repository.Src.Get(infrastructurePath)
    if err != nil {
        log.Sugar().Errorf("get infrastructure from repository err of %v, under path %s", err, infrastructurePath)
        return err, nil, "", nil
    }
    if infrastructureData == nil {
        log.Sugar().Infof("get nil infrastructure under path %s", infrastructurePath)
        return errors.New(fmt.Sprintf("No infrastructure under path %s", infrastructurePath)), nil, "", nil
    }
    deploymentInfo, err := template.GetDeploymentInfo(serviceData, infrastructureData)
    if err != nil {
        log.Sugar().Errorf("get deploymentInfo err of %v, servicepath:%s, infrastructurepath:%s", err, servicePath, infrastructurePath)
        return err, nil, "", nil
    }
    return nil, nil, util.Join("_", version, env, cluster, template.DeploymentInfo) + ".json", []byte(deploymentInfo)
}
