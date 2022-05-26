package generation

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"
)

func FillTemplates(dplyStructList []ChartDeployMain, rawData map[string][]byte, resMap map[string][]byte) (string, error) {
	//按照部署信息，循环到node，生成每个node的配置
	dirPrePath := ""
	for _, platformIns := range dplyStructList {
		for _, nodeTypeIns := range platformIns.NodeTypeList {
			//检查该类型是否需要填充，否则跳过。一次性填充该类型的所有模板。
			for path, inFile := range rawData {
				if !strings.Contains(path, util.Join("/", platformIns.Platform, nodeTypeIns.NodeType, repository.Template)) {
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
