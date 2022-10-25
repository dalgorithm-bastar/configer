package generation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/util"
	"gopkg.in/yaml.v3"
)

func GenerateDeploymentInfo(infrastructure []byte, rawSlice []RawFile, envNum string) ([]ChartDeployMain, error) {
	//返回值
	var chartDeployMain []ChartDeployMain
	//解析基础设施信息
	var infraStruct InfraMain
	//已经校验过，不存在错误
	_ = yaml.Unmarshal(infrastructure, &infraStruct)
	/*if err != nil {
	      return nil, fmt.Errorf("load infrastructure err: %s", err.Error())
	  }
	  //校验基础设施信息的键值是否符合预期
	  infraDecode, _ := yaml.Marshal(infraStruct)
	  keysOk := util.CheckYaml(infrastructure, infraDecode)
	  if !keysOk {
	      return nil, errors.New("unexpected keys on infrastructure")
	  }*/
	//分配和构建部署信息
	var nodeIdDispatched, setIdDispatched uint16 = 1, 1
	setIdMap, nodeIdMap, err := getIdMap(rawSlice)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range rawSlice {
		path := fileInfo.Path
		data := fileInfo.Data
		//找出deploment.yaml文件
		if strings.Contains(path, define.Deployment) {
			//校验长度
			pathSlice := strings.SplitN(path, "/", 7)
			if len(pathSlice) != 7 {
				return nil, errors.New(fmt.Sprintf("deployment.yaml in wrong place ,with path %s", path))
			}
			/*//校验id范围
			  if uint32(nodeId)+1 > 65535 {
			      return nil, errors.New("nodeId used up, current num over 65535")
			  }*/
			//填充和扩展部署信息
			var deployStruct DeployMain
			err := yaml.Unmarshal(data, &deployStruct)
			if err != nil {
				return nil, err
			}
			//校验部署信息的键值是否符合预期
			deployDecode, _ := yaml.Marshal(deployStruct)
			dpOk := util.CheckYaml(data, deployDecode)
			if !dpOk {
				return nil, fmt.Errorf("unexpected keys on %s, please checkout carefully", path)
			}
			if deployStruct.UserName != "" {
				deployStruct.UserName = strings.ReplaceAll(deployStruct.UserName, "@@", envNum)
			}
			if deployStruct.SetID == 0 {
				for ; setIdDispatched <= 65535; setIdDispatched++ {
					if _, ok := setIdMap[setIdDispatched]; !ok {
						deployStruct.SetID = setIdDispatched
						setIdMap[setIdDispatched] = 1
						setIdDispatched++
						break
					}
					if setIdDispatched == 65535 {
						return nil, fmt.Errorf("uint16 setId overflow, please checkout set quantity")
					}
				}
			}
			if deployStruct.SetID == 0 {
				return nil, fmt.Errorf("lack of arg: setID, filepath:%s", path)
			}
			for i1, node := range deployStruct.Node {
				isHostExist := false
				for _, host := range infraStruct.Host {
					if host.HostName == node.HostName {
						nodeIdTmp := deployStruct.Node[i1].NodeId
						deployStruct.Node[i1] = host
						deployStruct.Node[i1].NodeId = nodeIdTmp
						if deployStruct.Node[i1].NodeId == 0 {
							for ; nodeIdDispatched <= 65535; nodeIdDispatched++ {
								if _, ok := nodeIdMap[nodeIdDispatched]; !ok {
									deployStruct.Node[i1].NodeId = nodeIdDispatched
									nodeIdMap[nodeIdDispatched] = 1
									nodeIdDispatched++
									break
								}
								if nodeIdDispatched == 65535 {
									return nil, fmt.Errorf("uint16 nodeId overflow, please checkout node quantity")
								}
							}
						}
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
							SetID:      deployStruct.SetID,
							SetIndex:   setIndex,
							SetName:    pathSlice[5],
							Deployment: deployStruct,
						}}},
					},
				})
			}
			if isPlatExist && !isNodeTypeExist {
				//添加一个新的节点类型
				//向结构体中添加
				var setIndex uint16 = 1
				chartDeployMain[platIndex].NodeTypeList = append(chartDeployMain[platIndex].NodeTypeList, ChartDeployNodeType{
					NodeType: pathSlice[3],
					SetList: []ChartDeploySet{{
						SetID:      deployStruct.SetID,
						SetIndex:   setIndex,
						SetName:    pathSlice[5],
						Deployment: deployStruct,
					}},
				})
			}
			if isPlatExist && isNodeTypeExist {
				//添加一个新集群
				setIndex := len(chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList) + 1
				chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList = append(chartDeployMain[platIndex].NodeTypeList[nodeTypeIndex].SetList, ChartDeploySet{
					SetID:      deployStruct.SetID,
					SetIndex:   uint16(setIndex),
					SetName:    pathSlice[5],
					Deployment: deployStruct,
				})
			}
		}
	}
	return chartDeployMain, nil
}

func getIdMap(rawSlice []RawFile) (map[uint16]uint16, map[uint16]uint16, error) {
	setIdMap, nodeIdMap := make(map[uint16]uint16), make(map[uint16]uint16)
	for _, fileInfo := range rawSlice {
		path := fileInfo.Path
		data := fileInfo.Data
		//找出deploment.yaml文件
		if strings.Contains(path, define.Deployment) {
			//校验长度
			pathSlice := strings.SplitN(path, "/", 7)
			if len(pathSlice) != 7 {
				return nil, nil, errors.New(fmt.Sprintf("deployment.yaml in wrong place ,with path %s", path))
			}
			//填充和扩展部署信息
			var deployStruct DeployMain
			err := yaml.Unmarshal(data, &deployStruct)
			if err != nil {
				return nil, nil, err
			}
			//校验部署信息的键值是否符合预期
			deployDecode, _ := yaml.Marshal(deployStruct)
			dpOk := util.CheckYaml(data, deployDecode)
			if !dpOk {
				return nil, nil, fmt.Errorf("unexpected keys on %s, please checkout carefully", path)
			}
			if deployStruct.SetID != 0 {
				if _, ok := setIdMap[deployStruct.SetID]; ok {
					return nil, nil, fmt.Errorf("duplicated setId assigned: %d, file:%s", deployStruct.SetID, path)
				}
				setIdMap[deployStruct.SetID] = 1
			}
			for i, _ := range deployStruct.Node {
				if deployStruct.Node[i].NodeId != 0 {
					if _, ok := nodeIdMap[deployStruct.Node[i].NodeId]; ok {
						return nil, nil, fmt.Errorf("duplicated nodeId assigned: %d, file:%s", deployStruct.Node[i].NodeId, path)
					}
					nodeIdMap[deployStruct.Node[i].NodeId] = 1
				}
			}
		}
	}
	return setIdMap, nodeIdMap, nil
}
