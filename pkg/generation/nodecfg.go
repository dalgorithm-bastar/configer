package generation

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/util"
)

func FillTemplates(dplyStructList []ChartDeployMain, rawFile []RawFile, hostTcpPortMap map[string]hostTcpUnit, tcpRangeSli []int, envNum string, resMap map[string][]byte) (string, error) {
	//按照部署信息，循环到node，生成每个node的配置
	dirPrePath := ""
	for _, platformIns := range dplyStructList {
		for _, nodeTypeIns := range platformIns.NodeTypeList {
			//检查该类型是否需要填充，否则跳过。一次性填充该类型的所有模板。
			for ri, _ := range rawFile {
				path := rawFile[ri].Path
				inFile := rawFile[ri].Data
				if !strings.Contains(path, util.Join("/", platformIns.Platform, nodeTypeIns.NodeType, define.Template)) {
					continue
				}
				sli := strings.SplitN(path, "/", 6)
				if len(sli) < 6 {
					return "", errors.New(fmt.Sprintf("err template path: %s", path))
				}
				for _, setIns := range nodeTypeIns.SetList {
					for _, nodeIns := range setIns.Deployment.Node {
						fillArgs := FillArgs{
							PlatName:  platformIns.Platform,
							NodeType:  nodeTypeIns.NodeType,
							NodeNum:   uint16(len(setIns.Deployment.Node)),
							SetId:     setIns.SetID,
							SetIndex:  setIns.SetIndex,
							SetName:   setIns.SetName,
							NodeId:    nodeIns.NodeId,
							NodeIndex: nodeIns.NodeIndex,
						}
						sli[4] = util.Join("_", nodeIns.HostName, strconv.Itoa(int(nodeIns.NodeId)))
						resPath := util.Join("/", sli[0]+"_"+sli[1], _origin, sli[2], sli[3], setIns.SetName, sli[4], sli[5])
						if dirPrePath == "" {
							dirPrePath = sli[0] + "_" + sli[1]
						}
						//处理空文件
						if inFile == nil || len(inFile) == 0 {
							resMap[resPath] = []byte{}
							continue
						}
						model := template.New("fillConf")
						errorInfo := util.Join(".", platformIns.Platform, nodeTypeIns.NodeType, setIns.SetName, strconv.Itoa(int(nodeIns.NodeIndex)))
						model = model.Funcs(template.FuncMap{
							"GetIpByNet": func(netName string) (string, error) {
								return getIpByNet(netName, errorInfo, nodeIns)
							},
							"GetNextTcpPort": func() (string, error) {
								return getNextTcpPort(nodeIns.HostName, tcpRangeSli, envNum, hostTcpPortMap)
							},
						})
						model, err := model.Parse(string(inFile))
						if err != nil {
							return "", errors.New(fmt.Sprintf("Parse cfg err in path: %s, with err %s", path, err.Error()))
						}
						var data bytes.Buffer
						err = model.ExecuteTemplate(&data, "fillConf", fillArgs)
						if err != nil {
							return "", errors.New(fmt.Sprintf("Execute cfg err in path: %s, with err %s", path, err.Error()))
						}
						resMap[resPath] = data.Bytes()
					}
				}
			}
		}
	}
	return dirPrePath, nil
}

func getIpByNet(netName string, errorInfo string, hostInfo InfraHostUnit) (string, error) {
	for _, netInfo := range hostInfo.Network {
		if netInfo.Name == netName {
			return netInfo.Ipv4, nil
		}
	}
	return "", fmt.Errorf("cannot find ip of net:%s on host:%s, when processing node:%s", netName, hostInfo.HostName, errorInfo)
}

func getNextTcpPort(hostName string, tcpRangeSli []int, envNum string, hostTcpPortMap map[string]hostTcpUnit) (string, error) {
	if _, ok := hostTcpPortMap[hostName]; !ok {
		hostTcpPortMap[hostName] = hostTcpUnit{
			actualTcpPort: uint16(tcpRangeSli[0]),
			PortIdx:       0,
			coverTcpMap:   make(map[int]int),
		}
	}
	actualListenPort := int(hostTcpPortMap[hostName].actualTcpPort)
	actualListenPort, coverListenPort, currentIdx, listenPortOverFlow := getNextListenPort(tcpRangeSli, int(hostTcpPortMap[hostName].PortIdx), actualListenPort, envNum, hostTcpPortMap[hostName].coverTcpMap)
	if listenPortOverFlow {
		return "", fmt.Errorf("get next tcp port err on host:%s, tcp port used up", hostName)
	}
	//更新端口池
	tmp := hostTcpUnit{
		actualTcpPort: uint16(actualListenPort),
		PortIdx:       uint16(currentIdx),
		coverTcpMap:   hostTcpPortMap[hostName].coverTcpMap,
	}
	hostTcpPortMap[hostName] = tmp
	return coverListenPort, nil
}
