// tmplFuncs.go文件用于定义和实现模板中可以使用的函数，以及定义部署信息生成函数
package template

import (
    "encoding/json"
    "errors"
    "fmt"
    "math"
    "strconv"
    "strings"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/util"
    "github.com/tidwall/sjson"
)

// Atoi 在模板中进行类型转换string->int
func ParseFloat(s string) (float64, error) {
    return strconv.ParseFloat(s, 64)
}

// Itoa 在模板中进行类型转换int->string
func FmtFloat64(i float64) string {
    return strconv.FormatFloat(i, 'f', 0, 64)
}

func Add(i float64) float64 {
    return i + 1
}

func Mine(i float64) float64 {
    return i - 1
}

// CtlFind 命令行单条信息查询函数，不会对service查询结果进行映射
func CtlFind(tar, ver, env, clusterObject, service string) (string, error) {
    var filePath string
    var data interface{}
    switch tar {
    case Services:
        filePath = util.Join("/", ver, env, clusterObject, repository.ServiceList)
    case Infrastructure:
        filePath = util.Join("/", ver, repository.Infrastructure)
    }
    binaryData, err := repository.Src.Get(filePath) //从数据源取值
    if err != nil {
        log.Sugar().Errorf("get filedata from repository err of %v in CtlFind, under path %s", err, filePath)
        return "", err
    }
    if binaryData == nil {
        errInfo := fmt.Sprintf("no data under path %s in CtlFind", filePath)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
    err = json.Unmarshal(binaryData, &data)
    if err != nil {
        log.Sugar().Infof("json unmarshal fileData err of %v, data:%s", err, binaryData)
        return "", err
    }
    dataMap := make(map[string]string)
    ConstructMap(dataMap, data, "")
    if v, ok := dataMap[service]; ok {
        return v, nil
    }
    errInfo := fmt.Sprintf("can not find \"%s\" in \"%s\"", service, filePath)
    log.Sugar().Info(errInfo)
    return "", errors.New(errInfo)
}

// GetInfobyNodeId 将给出的nodeid翻译成servicelist上实例链表的序号，再调用底层函数实现替换
func GetInfobyNodeId(src repository.Storage, infrastructureData []byte, defaultIndex bool, globalId, nodeId, ver, env, clusterObject, service string) (string, error) {
    //构建serviceList map
    prefixService := util.Join("/", ver, env, clusterObject, repository.ServiceList)
    var data interface{}
    binaryData, err := src.Get(prefixService) //从数据源取值
    if err != nil {
        log.Sugar().Errorf("get serviceList from repository err of %v in baseGet, under path %s", err, prefixService)
        return "", err
    }
    if binaryData == nil {
        errInfo := fmt.Sprintf("no data under path %s in baseGet, please checkout in etcd or compressedfile", prefixService)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
    err = json.Unmarshal(binaryData, &data)
    if err != nil {
        log.Sugar().Infof("json unmarshal serviceList err of %v, data:%s", err, binaryData)
        return "", err
    }
    serviceMap := make(map[string]string)
    ConstructMap(serviceMap, data, "")
    //查找是否存在该节点号对应的nodeid，找不到则返回
    var localId string
    //是否使用隐式索引
    for k, v := range serviceMap {
        if strings.Contains(k, DeploymentInfoKey) && strings.Contains(k, NodeIdKey) && v == nodeId {
            keySlice := strings.SplitN(k, ".", 3)
            if len(keySlice) < 3 {
                errInfo := fmt.Sprintf("err key with nodeid, key is %s", k)
                log.Sugar().Info(errInfo)
                return "", errors.New(errInfo)
            }
            localId = keySlice[1]
        }
    }
    return baseGet(src, infrastructureData, defaultIndex, globalId, localId, ver, env, clusterObject, service)
}

func baseGet(src repository.Storage, infrastructureData []byte, defaultIndex bool, globalId, localId, ver, env, clusterObject, service string) (string, error) {
    //构建serviceList map
    prefixService := util.Join("/", ver, env, clusterObject, repository.ServiceList)
    var data interface{}
    binaryData, err := src.Get(prefixService) //从数据源取值
    if err != nil {
        log.Sugar().Errorf("get serviceList from repository err of %v in baseGet, under path %s", err, prefixService)
        return "", err
    }
    if binaryData == nil {
        errInfo := fmt.Sprintf("no data under path %s in baseGet, please checkout in etcd or compressedfile", prefixService)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
    err = json.Unmarshal(binaryData, &data)
    if err != nil {
        log.Sugar().Infof("json unmarshal serviceList err of %v, data:%s", err, binaryData)
        return "", err
    }
    serviceMap := make(map[string]string)
    ConstructMap(serviceMap, data, "")
    //查找是否存在该项服务，找不到则返回
    serviceValue := ""
    //是否使用隐式索引
    if !defaultIndex {
        if v, ok := serviceMap[service]; ok {
            serviceValue = v
        }
    } else {
        for k, v := range serviceMap {
            //加快循环速度
            if len(k) != len(service)+len(localId)+1 {
                continue
            }
            targetSlice := strings.Split(service, ".")
            if len(targetSlice) < 2 {
                return "", errors.New(fmt.Sprintf("orgnization of servicelist illigal:%s", service))
            }
            //插入localid
            targetSlice = append(targetSlice[:1], append([]string{localId}, targetSlice[1:]...)...)
            targetKey := strings.Join(targetSlice, ".")
            if k == targetKey {
                serviceValue = v
                break
            }
        }
    }
    // 替换项为空值，报异常
    if serviceValue == "" {
        errInfo := fmt.Sprintf("No such service or index in serviceList: %s, %s", service, prefixService)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
    //服务值不包含"{}",直接返回，否则需要替换
    if len(serviceValue) < 2 || serviceValue[0] != '{' || serviceValue[len(serviceValue)-1] != '}' {
        return serviceValue, nil
    }
    if len(serviceValue) == 2 {
        log.Sugar().Infof("get nil value in {}, key:%s, value:%s", service, serviceValue)
        return "", errors.New("cannot set nil value within {}")
    }
    //去除括号
    serviceValue = serviceValue[1 : len(serviceValue)-1]
    //未传入公共信息文件时，获取文件
    if infrastructureData == nil {
        prefixPublic := util.Join("/", ver, repository.Infrastructure)
        infrastructureData, err = src.Get(prefixPublic) //从数据源取值
        if err != nil {
            log.Sugar().Errorf("get infrastructure from repository err of %v in baseGet, under path %s", err, prefixPublic)
            return "", err
        }
        if infrastructureData == nil {
            log.Sugar().Infof("get nil infrastructure under path %s in baseGet", prefixPublic)
            return "", errors.New(fmt.Sprintf("no data under path %s, please checkout in etcd or compressedfile", prefixPublic))
        }
    }
    err = json.Unmarshal(infrastructureData, &data)
    if err != nil {
        log.Sugar().Infof("json unmarshal infrastructure err of %v, data:%s", err, infrastructureData)
        return "", err
    }
    publicInfoMap := make(map[string]string)
    ConstructMap(publicInfoMap, data, "")
    //查找是否存在映射,替换或返回
    //按照localid和约定的主机名称键值获取hostname
    hostname := ""
    for srvKey, srvVal := range serviceMap {
        targetKey := util.Join(".", DeploymentInfoKey, localId, HostNameKey)
        if srvKey == targetKey {
            hostname = srvVal
            break
        }
    }
    if hostname == "" {
        errInfo := fmt.Sprintf("No hostname under localId %s on servicelist %s", localId, prefixService)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
    //映射
    targetKey := hostname + "." + serviceValue
    for k, v := range publicInfoMap {
        if strings.HasSuffix(k, targetKey) && v != "" {
            return v, nil
        }
    }
    errInfo := fmt.Sprintf("no such infrastructure info of %s on %s", targetKey, util.Join("/", ver, env, repository.Infrastructure))
    log.Sugar().Info(errInfo)
    return "", errors.New(errInfo)
}

// ConstructMap 从任意json构建key为完整路径的map
func ConstructMap(resMap map[string]string, data interface{}, currentPath string) {
    switch data.(type) {
    case string:
        resMap[currentPath[0:len(currentPath)-1]] = data.(string)
    case float64:
        resMap[currentPath[0:len(currentPath)-1]] = strconv.FormatFloat(data.(float64), 'f', 0, 64)
    case []interface{}:
        //内置类型直接返回csv
        if interfaceSlice, ok := data.([]interface{}); ok {
            isBaseData := true
            var stringSlice []string
            for _, v := range interfaceSlice {
                _, isString := v.(string)
                _, isFloat := v.(float64)
                if !isString && !isFloat {
                    isBaseData = false
                    break
                }
                if isString {
                    stringSlice = append(stringSlice, v.(string))
                }
                if isFloat {
                    stringSlice = append(stringSlice, strconv.Itoa(int(math.Floor(v.(float64)+0.5)))) //四舍五入取整
                }
            }
            if isBaseData {
                resMap[currentPath[0:len(currentPath)-1]] = strings.Join(stringSlice, ",")
                return
            }
        }
        for i, v := range data.([]interface{}) {
            ConstructMap(resMap, v, currentPath+strconv.Itoa(i)+".")
        }
    case map[string]interface{}:
        for k, v := range data.(map[string]interface{}) {
            ConstructMap(resMap, v, currentPath+k+".")
        }
    }
}

// GetDeploymentInfo 接收两份文件，返回一份新的部署信息文件
func GetDeploymentInfo(serviceData, infrastructureData []byte) (string, error) {
    if serviceData == nil || infrastructureData == nil {
        return "", errors.New("nil input when filling service data")
    }
    var serviceInterface, infrastructureInterface interface{}
    err := json.Unmarshal(serviceData, &serviceInterface)
    if err != nil {
        log.Sugar().Infof("json unmarshal serviceList err of %v, servicedata:%s", err, serviceData)
        return "", err
    }
    err = json.Unmarshal(infrastructureData, &infrastructureInterface)
    if err != nil {
        log.Sugar().Infof("json unmarshal infrastructure err of %v, infrastructureData:%s", err, infrastructureData)
        return "", err
    }
    serviceMap := make(map[string]string)
    ConstructMap(serviceMap, serviceInterface, "")
    if _, ok := serviceMap[ReplicatorNumKey]; !ok {
        log.Sugar().Infof("Missing replicator_number in serviceList, servicedata %s", serviceData)
        return "", errors.New("Missing replicator_number in serviceList")
    }
    deploymentInfo, err := sjson.Set("", ReplicatorNumKey, serviceMap[ReplicatorNumKey])
    if err != nil {
        log.Sugar().Errorf("insert deploymentInfo failed, datatoInsert:%+v, deploymentInfo:%s", serviceMap, deploymentInfo)
        return "", errors.New("insert deploymentInfo failed")
    }
    infrastructureMap := make(map[string]string)
    ConstructMap(infrastructureMap, infrastructureInterface, "")
    //筛选部署信息
    for k, v := range serviceMap {
        if strings.HasPrefix(k, DeploymentInfoKey) {
            valueToWrite := v
            if len(v) > 1 && v[0] == '{' && v[len(v)-1] == '}' {
                if len(v) <= 2 {
                    log.Sugar().Infof("get nil value in {}, key:%s, value:%s", k, v)
                    return "", errors.New(fmt.Sprintf("can not set nil value in {}, key:%s, value:%s", k, v))
                }
                valueToWrite = v[1 : len(v)-1]
                //取localid,按照localid查询hostname
                serviceKeySlice := strings.SplitN(k, ".", 3)
                if len(serviceKeySlice) < 3 {
                    log.Sugar().Warnf("Get special service value:%s of key:%s", v, k)
                    continue
                }
                keyofHostName := util.Join(".", DeploymentInfoKey, serviceKeySlice[1], HostNameKey)
                if _, ok := serviceMap[keyofHostName]; !ok {
                    err := fmt.Sprintf("can not find hostname key of %s on servicelist", keyofHostName)
                    log.Sugar().Info(err)
                    return "", errors.New(err)
                }
                hostName := serviceMap[keyofHostName]
                suffixKey := hostName + "." + valueToWrite
                for infraKey, infraVal := range infrastructureMap {
                    if strings.HasSuffix(infraKey, suffixKey) {
                        valueToWrite = infraVal
                        break
                    }
                }
                if valueToWrite == v[1:len(v)-1] {
                    err := fmt.Sprintf("can not find infrastructure info of suffixkey %s", suffixKey)
                    log.Sugar().Info(err)
                    return "", errors.New(err)
                }
            }
            deploymentInfo, err = sjson.Set(deploymentInfo, k, valueToWrite)
            if err != nil {
                log.Sugar().Errorf("insert deploymentInfo failed, datatoInsert:%+v, deploymentInfo:%s", serviceMap, deploymentInfo)
                return "", errors.New("insert deploymentInfo failed")
            }
        }
    }
    return deploymentInfo, nil
}
