// Package generation thirdPart文件用于按照第三方所需的输出格式生成配置文件
package generation

import (
	"encoding/json"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/configcenter/pkg/util"
	"gopkg.in/yaml.v3"
)

const (
	_huaruiDir = "huarui_files"
	_hostCfg   = "host.cfg"
	_nodeCfg   = "node.cfg"
	_routeCfg  = "route.cfg"
	_bizNet    = "biznet"
)

func huaRuiMain(infrastructure []byte, deployment []ChartDeployMain, topicList []byte) (map[string][]byte, error) {
	resMap := make(map[string][]byte)
	//获取系统信息，确定换行符
	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}
	hostCfgData := huaRuiGetHostCfg(infrastructure, newLine)
	resMap[_huaruiDir+"/"+_hostCfg] = hostCfgData
	nodeCfgData := huaRuiGetNodeCfg(deployment, newLine)
	resMap[_huaruiDir+"/"+_nodeCfg] = nodeCfgData
	routeCfgData := huaRuiGetRouteCfg(topicList, newLine)
	resMap[_huaruiDir+"/"+_routeCfg] = routeCfgData
	return resMap, nil
}

func huaRuiGetHostCfg(infrastructure []byte, newLine string) []byte {
	var resData = ""
	var infraStruct InfraMain
	_ = yaml.Unmarshal(infrastructure, &infraStruct)
	for _, host := range infraStruct.Host {
		for i, _ := range host.Network {
			if host.Network[i].Name == _bizNet {
				if resData == "" {
					resData = host.HostName + "|" + host.Network[i].Ipv4 + "|" + " "
				} else {
					resData = resData + newLine + host.HostName + "|" + host.Network[i].Ipv4 + "|" + " "
				}
			}
		}
	}
	return []byte(resData)
}

func huaRuiGetNodeCfg(deployment []ChartDeployMain, newLine string) []byte {
	var resData = ""
	for _, platIns := range deployment {
		for _, nodeTypeIns := range platIns.NodeTypeList {
			nodeTypeName := nodeTypeIns.NodeType
			for _, setIns := range nodeTypeIns.SetList {
				setId := setIns.SetID
				for _, nodeIns := range setIns.Deployment.Node {
					if resData == "" {
						resData = nodeIns.HostName + "|" + nodeTypeName + "|" + strconv.Itoa(int(setId)) +
							"|" + strconv.Itoa(int(nodeIns.NodeId))
					} else {
						resData = resData + newLine + nodeIns.HostName + "|" + nodeTypeName + "|" + strconv.Itoa(int(setId)) +
							"|" + strconv.Itoa(int(nodeIns.NodeId))
					}
				}
			}
		}
	}
	return []byte(resData)
}

func huaRuiGetRouteCfg(topicList []byte, newLine string) []byte {
	resData, msgSlice := "", make([]string, 0)
	var topicListStruct ChartTpc
	msgMap := make(map[string]int)
	_ = json.Unmarshal(topicList, &topicListStruct)
	//暴力法
	for _, platIns := range topicListStruct.Platforms {
		for _, nodeTypeIns := range platIns.NodeTypeList {
			nodeTypeSub := nodeTypeIns.NodeType
			for _, setIns := range nodeTypeIns.SetList {
				setIdSub := setIns.SetID
				//处理bizTopic
				for _, bizTopicIns := range setIns.BroadInfo.SubExtern.BizTopic {
					pubSetSlice := strings.Split(bizTopicIns.PubCluster, ".")
					setIdPub := bizTopicIns.PubSetId
					isRmb := "N"
					if *bizTopicIns.IsRMB == 1 {
						isRmb = "Y"
					}
					msg := util.Join("|", "1", pubSetSlice[1], strconv.Itoa(int(setIdPub)), nodeTypeSub,
						strconv.Itoa(int(setIdSub)), bizTopicIns.EndPoint, strconv.Itoa(int(bizTopicIns.TopicId)), isRmb)
					if _, ok := msgMap[msg]; !ok {
						msgMap[msg] = 0
					}
				}
				//处理setTopic
				for _, setTopicIns := range setIns.BroadInfo.SubExtern.SetTopic {
					pubSetSlice := strings.Split(setTopicIns.PubCluster, ".")
					setIdPub := setTopicIns.PubSetId
					isRmb := "N"
					if *setTopicIns.IsRMB == 1 {
						isRmb = "Y"
					}
					msg := util.Join("|", "1", pubSetSlice[1], strconv.Itoa(int(setIdPub)), nodeTypeSub,
						strconv.Itoa(int(setIdSub)), setTopicIns.EndPoint, strconv.Itoa(int(setTopicIns.TopicId)), isRmb)
					if _, ok := msgMap[msg]; !ok {
						msgMap[msg] = 0
					}
				}
			}
		}
	}
	for k, _ := range msgMap {
		msgSlice = append(msgSlice, k)
	}
	sort.Strings(msgSlice)
	for i, _ := range msgSlice {
		if resData == "" {
			resData = msgSlice[i]
		} else {
			resData = resData + newLine + msgSlice[i]
		}
	}
	return []byte(resData)
}
