package generation

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"
	"gopkg.in/yaml.v3"
)

const (
	DEPLOYLIST = "deployList.yaml"
	TOPICLIST  = "topicList.json"
	ORIGIN     = "origin_res"
)

//RawFile 用于定序生成
type RawFile struct {
	Path string
	Data []byte
}

// Generate 由标准输入生成标准输出
func Generate(infrastructure []byte, rawData map[string][]byte, envNum string, ipRange, portRange []string) (map[string][]byte, error) {
	//校验ip和port
	if len(ipRange)%2 != 0 || len(portRange)%2 != 0 || len(ipRange) == 0 || len(portRange) == 0 {
		return nil, errors.New("err iprange or portrange length")
	}
	for _, ip := range ipRange {
		legal := isIpv4Legal(ip)
		if !legal {
			return nil, errors.New(fmt.Sprintf("err ip input:%s", ip))
		}
	}
	for _, port := range portRange {
		_, legal := checkPort(port)
		if !legal {
			return nil, errors.New(fmt.Sprintf("err port input:%s", port))
		}
	}
	//删除空文件，配置模板中的空文件除外
	for k, _ := range rawData {
		if len(rawData[k]) == 0 && !strings.Contains(k, repository.Template) {
			delete(rawData, k)
		}
	}
	//对rawData的键排序，转换成rawSlice
	rawSlice := sortRawData(rawData)
	//构造返回结果文件包
	resMap := make(map[string][]byte)
	//扩充部署信息
	dplyStructList, err := GenerateDeploymentInfo(infrastructure, rawSlice)
	if err != nil {
		return nil, err
	}
	//该方案不部署，提前返回
	if len(dplyStructList) == 0 {
		return resMap, nil
	}
	//生成topicInfo总表
	topicInfoList, err := GenerateTopicInfo(dplyStructList, rawSlice, ipRange, portRange, envNum)
	if err != nil {
		return nil, err
	}
	//按照部署信息填充模板，同时获取文件夹前缀
	versionHead, err := FillTemplates(dplyStructList, rawData, resMap)
	if err != nil {
		return nil, err
	}
	//处理总表，补充resMap
	err = FinishResMap(resMap, dplyStructList, topicInfoList, versionHead+"/"+ORIGIN)
	if err != nil {
		return nil, err
	}
	//处理第三方文件生成
	err = addThirdPartFiles(resMap, infrastructure, dplyStructList, versionHead)
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

func sortRawData(rawData map[string][]byte) []RawFile {
	var keySlice []string
	var rawSlice []RawFile
	for path, _ := range rawData {
		keySlice = append(keySlice, path)
	}
	sort.Strings(keySlice)
	for _, path := range keySlice {
		file := RawFile{
			Path: path,
			Data: rawData[path],
		}
		rawSlice = append(rawSlice, file)
	}
	return rawSlice
}

func FinishResMap(resMap map[string][]byte, dplyStructList []ChartDeployMain, topicInfoList map[string]map[string]map[string]ExpTpcMain, prePath string) error {
	var chartTpc ChartTpc
	//按部署信息循环
	for _, platformIns := range dplyStructList {
		for _, nodeTypeIns := range platformIns.NodeTypeList {
			//检查是否有组播信息
			isTpcExst := false
			if _, ok1 := topicInfoList[platformIns.Platform]; ok1 {
				if _, ok2 := topicInfoList[platformIns.Platform][nodeTypeIns.NodeType]; ok2 {
					isTpcExst = true
				}
			}
			for _, setIns := range nodeTypeIns.SetList {
				commonPath := util.Join("/", prePath, platformIns.Platform, nodeTypeIns.NodeType, setIns.SetName)
				pathDp := commonPath + "/" + repository.Deployment
				dpFile, err := yaml.Marshal(setIns.Deployment)
				if err != nil {
					return errors.New(fmt.Sprintf("json marshal dpfile err: %s, path: %s", err.Error(), commonPath))
				}
				resMap[pathDp] = dpFile
				if isTpcExst {
					pathTpc := commonPath + "/" + repository.TopicInfo
					tpcFile, err := json.Marshal(topicInfoList[platformIns.Platform][nodeTypeIns.NodeType][setIns.SetName])
					if err != nil {
						return errors.New(fmt.Sprintf("json marshal tpcfile err: %s, path: %s", err.Error(), commonPath))
					}
					resMap[pathTpc] = tpcFile
					//处理组播总表
					isPlat, isNodeType := false, false
					iP, iN := 0, 0
					for i, _ := range chartTpc.Platforms {
						if chartTpc.Platforms[i].Platform == platformIns.Platform {
							isPlat = true
							iP = i
							for j, _ := range chartTpc.Platforms[i].NodeTypeList {
								if chartTpc.Platforms[i].NodeTypeList[j].NodeType == nodeTypeIns.NodeType {
									isNodeType = true
									iN = j
									break
								}
							}
							break
						}
					}
					if !isPlat && !isNodeType {
						chartTpc.Platforms = append(chartTpc.Platforms, ChartTpcMain{
							Platform: platformIns.Platform,
							NodeTypeList: []ChartTpcNodeType{{
								NodeType: nodeTypeIns.NodeType,
								SetList: []ChartTpcSet{{
									SetID:     setIns.SetID,
									SetName:   setIns.SetName,
									BroadInfo: topicInfoList[platformIns.Platform][nodeTypeIns.NodeType][setIns.SetName],
								}},
							}},
						})
					}
					if isPlat && !isNodeType {
						chartTpc.Platforms[iP].NodeTypeList = append(chartTpc.Platforms[iP].NodeTypeList, ChartTpcNodeType{
							NodeType: nodeTypeIns.NodeType,
							SetList: []ChartTpcSet{{
								SetID:     setIns.SetID,
								SetName:   setIns.SetName,
								BroadInfo: topicInfoList[platformIns.Platform][nodeTypeIns.NodeType][setIns.SetName],
							}},
						})
					}
					if isPlat && isNodeType {
						chartTpc.Platforms[iP].NodeTypeList[iN].SetList = append(chartTpc.Platforms[iP].NodeTypeList[iN].SetList, ChartTpcSet{
							SetID:     setIns.SetID,
							SetName:   setIns.SetName,
							BroadInfo: topicInfoList[platformIns.Platform][nodeTypeIns.NodeType][setIns.SetName],
						})
					}
				}
			}
		}
	}
	//添加部署总表
	dpMain := ChartDeploy{
		Scheme:    "",
		Platforms: dplyStructList,
	}
	mainDpFile, err := yaml.Marshal(dpMain)
	if err != nil {
		return errors.New(fmt.Sprintf("yaml marshal maindpfile err: %s", err.Error()))
	}
	resMap[prePath+"/"+DEPLOYLIST] = mainDpFile
	mainTpcFile, err := json.Marshal(chartTpc)
	if err != nil {
		return errors.New(fmt.Sprintf("json marshal maintpcfile err: %s", err.Error()))
	}
	resMap[prePath+"/"+TOPICLIST] = mainTpcFile
	return nil
}

func addThirdPartFiles(resMap map[string][]byte, infrastructure []byte, dplyStructList []ChartDeployMain, versionHead string) error {
	//生成华锐所需的配置文件
	huaRuiFileMap, err := huaRuiMain(infrastructure, dplyStructList, resMap[util.Join("/", versionHead, ORIGIN, TOPICLIST)])
	if err != nil {
		return err
	}
	for k, _ := range huaRuiFileMap {
		if _, ok := resMap[versionHead+"/"+k]; ok {
			return fmt.Errorf("repeated key in outputs between raw and huarui, key:%s", k)
		}
		resMap[versionHead+"/"+k] = huaRuiFileMap[k]
	}
	return nil
}
