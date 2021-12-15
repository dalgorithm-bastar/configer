// tmplFuncs.go文件用于定义和实现模板中可以使用的函数，以及定义部署信息生成函数
package template

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
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

func Itoa(i int) string {
    return strconv.Itoa(i)
}

func Atoi(s string) (int, error) {
    return strconv.Atoi(s)
}

// CtlFind 命令行单条信息查询函数，不会对service查询结果进行映射
func CtlFind(tar, ver, conf, clusterObject, service string) (string, error) {
    var filePath string
    var data interface{}
    switch tar {
    case Services:
        filePath = util.Join("/", ver, conf, clusterObject, repository.ServiceList)
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

func baseGet(src repository.Storage, infrastructureData []byte, defaultIndex bool, envNum, globalId, localId, ver, conf, clusterObject, service string) (string, error) {
    //获取serviceList文件
    prefixService := util.Join("/", ver, conf, clusterObject, repository.ServiceList)
    var data interface{}
    srvData, err := src.Get(prefixService) //从数据源取值
    if err != nil {
        log.Sugar().Errorf("get serviceList from repository err of %v in baseGet, under path %s", err, prefixService)
        return "", err
    }
    if srvData == nil {
        errInfo := fmt.Sprintf("no data under path %s in baseGet, please checkout in etcd or compressedfile", prefixService)
        log.Sugar().Info(errInfo)
        return "", errors.New(errInfo)
    }
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
    //向基础设施信息中插入环境号信息
    infrastructureData, err = sjson.SetBytes(infrastructureData, EnvNumKey, envNum)
    if err != nil {
        return "", errors.New(err.Error() + "; quit when insert envNum")
    }
    //先行使用基础设施信息填充服务清单，再使用可用的服务清单继续查询
    srvData, err = fillSrvbyInfra(srvData, infrastructureData)
    if err != nil {
        log.Sugar().Infof("fill serviceList err of %v, data:%s", err, srvData)
        return "", err
    }
    err = json.Unmarshal(srvData, &data)
    if err != nil {
        log.Sugar().Infof("json unmarshal serviceList err of %v, data:%s", err, srvData)
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
    return serviceValue, nil
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

// 在预设函数中用于填充服务清单
func fillSrvbyInfra(servicelistInput, infrastructure []byte) ([]byte, error) {
    tmpl := template.New("forSrv")
    tmpl, err := tmpl.Parse(string(servicelistInput))
    if err != nil {
        return nil, errors.New(err.Error() + "; init tmpl for srv failed")
    }
    var infraData interface{}
    err = json.Unmarshal(infrastructure, &infraData)
    if err != nil {
        return nil, err
    }

    var data bytes.Buffer
    err = tmpl.ExecuteTemplate(&data, "forSrv", infraData)
    if err != nil {
        return nil, err
    }
    return data.Bytes(), nil
}
