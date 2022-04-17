package generation

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/configcenter/pkg/repository"
)

func GenerateDeploymentInfo(infrastrcture []byte, rawData map[string][]byte) ([]ChartDeployMain, error) {
    //返回值
    var chartDeployMain []ChartDeployMain
    //解析基础设施信息
    var infraStruct InfraMain
    err := json.Unmarshal(infrastrcture, &infraStruct)
    if err != nil {
        return nil, err
    }
    //分配和构建部署信息
    var setId, nodeId uint16 = 1, 1
    for path, data := range rawData {
        //找出deploment.json文件
        if strings.Contains(path, repository.Deployment) {
            //校验长度
            pathSlice := strings.SplitN(path, "/", 7)
            if len(pathSlice) != 7 {
                return nil, errors.New(fmt.Sprintf("deployment.json in wrong place ,with path %s", path))
            }
            //校验id范围
            if uint32(nodeId)+1 > 65535 {
                return nil, errors.New("nodeId used up, current num over 65535")
            }
            //填充和扩展部署信息
            var deployStruct DeployMain
            err := json.Unmarshal(data, &deployStruct)
            if err != nil {
                return nil, err
            }
            for i1, node := range deployStruct.Node {
                isHostExist := false
                for _, host := range infraStruct.Host {
                    if host.HostName == node.HostName {
                        deployStruct.Node[i1] = host
                        deployStruct.Node[i1].NodeId = nodeId
                        nodeId++
                        deployStruct.Node[i1].NodeIndex = uint16(i1) + 1
                        isHostExist = true
                        break
                    }
                }
                if !isHostExist {
                    return nil, errors.New(fmt.Sprintf("cannot find target host for node: %s, request for host: %s", path+",node:"+strconv.Itoa(i1+1), node.HostName))
                }
            }
            //向结构体插值，同时分配相关ID和索引
            isPlatExist, isNodeTypeExist := false, false
            platIndex, nodeTypeIndex := 0, 0
            for i1, chartPlatform := range chartDeployMain {
                if chartPlatform.Platform == pathSlice[2] {
                    isPlatExist = true
                    platIndex = i1
                    for i2, nodeType := range chartPlatform.NodeTypeList {
                        if nodeType.NodeType == pathSlice[3] {
                            isNodeTypeExist = true
                            nodeTypeIndex = i2
                            break
                        }
                    }
                    break
                }
            }
            if !isPlatExist && !isNodeTypeExist {
                //添加一个新的平台类型
                //向结构体中添加
                var setIndex uint16 = 1
                chartDeployMain = append(chartDeployMain, ChartDeployMain{
                    Platform: pathSlice[2],
                    NodeTypeList: []ChartDeployNodeType{{
                        NodeType: pathSlice[3],
                        SetList: []ChartDeploySet{{
                            SetID:      setId,
                            SetIndex:   setIndex,
                            SetName:    pathSlice[5],
                            Deployment: deployStruct,
                        }}},
                    },
                })
                setId++
            }
            if isPlatExist && !isNodeTypeExist {
                //添加一个新的节点类型
                //向结构体中添加
                var setIndex uint16 = 1
                chartDeployMain[platIndex].NodeTypeList = append(chartDeployMain[platIndex].NodeTypeList, ChartDeployNodeType{
                    NodeType: pathSlice[3],
                    SetList: []ChartDeploySet{{
                        SetID:      setId,
                        SetIndex:   setIndex,
                        SetName:    pathSlice[5],
                        Deployment: deployStruct,
                    }},
                })
                setId++
            }
            if isPlatExist && isNodeTypeExist {
                //添加一个新集群
                setIndex := len(chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList) + 1
                chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList = append(chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList, ChartDeploySet{
                    SetID:      setId,
                    SetIndex:   uint16(setIndex),
                    SetName:    pathSlice[5],
                    Deployment: deployStruct,
                })
                setId++
            }
        }
    }
    return chartDeployMain, nil
}
