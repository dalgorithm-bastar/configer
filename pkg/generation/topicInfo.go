package generation

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/configcenter/pkg/repository"
    "github.com/configcenter/pkg/util"
)

//描述biztopic订阅者信息
type idtfy struct {
    Plat     string
    NodeType string
}

//描述biztopic处理状况
type topicStat struct {
    IsProc bool
    Name   string
    Sub    []idtfy
}

//主切片排序单元
type mainSlicePlatUnit struct {
    PlatForm  string
    NodeTypes []mainSliceNodeTypeUnit
}

//主切片次级排序单元
type mainSliceNodeTypeUnit struct {
    NodeType string
    SrvFile  SrvMain
}

//用于主机tcp端口池
type hostTcpUnit struct {
    actualTcpPort uint16
    coverTcpMap   map[int]int
}

func GenerateTopicInfo(dplyStructList []ChartDeployMain, rawSlice []RawFile, ipRange []string, portRange []string, envNum string) (map[string]map[string]map[string]ExpTpcMain, error) {
    //初始化参数和返回结果
    topicInfoMap := make(map[string]map[string]map[string]ExpTpcMain)
    //netMap，网络-biztopic名称&处理状态-接收方
    netMap := make(map[string][]topicStat)
    //mainSlice,主切片，plat-nodetype-srvMainStruct，用于循环遍历
    var mainSlice []mainSlicePlatUnit
    //构造主切片
    for _, fileInfo := range rawSlice {
        path := fileInfo.Path
        data := fileInfo.Data
        //找出service.json
        if !strings.Contains(path, repository.Service) {
            continue
        }
        //校验长度
        pathSlice := strings.SplitN(path, "/", 6)
        if len(pathSlice) < 6 {
            return nil, errors.New(fmt.Sprintf("service.json in wrong place ,with path %s", path))
        }
        //读取并归档当前文件的服务声明
        var srvStruct SrvMain
        err := json.Unmarshal(data, &srvStruct)
        if err != nil {
            return nil, err
        }
        //检测当前平台或节点类型是否收录过
        isPlatExist, isNodeTypeExist, platIdx, nodeTypeIdx := false, false, 0, 0
        for i, _ := range mainSlice {
            if mainSlice[i].PlatForm == pathSlice[2] {
                isPlatExist = true
                platIdx = i
                for j, _ := range mainSlice[i].NodeTypes {
                    if mainSlice[i].NodeTypes[j].NodeType == pathSlice[3] {
                        isNodeTypeExist = true
                        nodeTypeIdx = j
                        break
                    }
                }
                break
            }
        }
        //平台与节点均未收录过
        if !isPlatExist && !isNodeTypeExist {
            mainSlice = append(mainSlice, mainSlicePlatUnit{
                PlatForm: pathSlice[2],
                NodeTypes: []mainSliceNodeTypeUnit{{
                    NodeType: pathSlice[3],
                    SrvFile:  srvStruct,
                }},
            })
        }
        //节点类型未被收录过
        if isPlatExist && !isNodeTypeExist {
            mainSlice[platIdx].NodeTypes = append(mainSlice[platIdx].NodeTypes, mainSliceNodeTypeUnit{
                NodeType: pathSlice[3],
                SrvFile:  srvStruct,
            })
        }
        //平台与节点类型均被收录过，进行合并
        if isPlatExist && isNodeTypeExist {
            //合并innerTopic
            if mainSlice[platIdx].NodeTypes[nodeTypeIdx].SrvFile.InnerTopicNet != "" {
                if srvStruct.InnerTopicNet != mainSlice[platIdx].NodeTypes[nodeTypeIdx].SrvFile.InnerTopicNet {
                    return nil, errors.New(fmt.Sprintf("Multi inner topic net assigned in nodeType: %s", path))
                }
            } else if srvStruct.InnerTopicNet != "" {
                mainSlice[platIdx].NodeTypes[nodeTypeIdx].SrvFile.InnerTopicNet = srvStruct.InnerTopicNet
            }
            //合并其余项
            //取map中的对应srv用于合并
            tmp := mainSlice[platIdx].NodeTypes[nodeTypeIdx].SrvFile
            //合并单个目标srv,同时构建netMap
            tmp, err := MergeSrvStruct(tmp, srvStruct)
            if err != nil {
                return nil, errors.New(fmt.Sprintf("merge srv err in path: %s;", path) + err.Error())
            }
            mainSlice[platIdx].NodeTypes[nodeTypeIdx].SrvFile = tmp
        }
        //同步更新netMap
        for _, subUnit := range srvStruct.SubTopic {
            MergeNetMap(idtfy{
                Plat:     pathSlice[2],
                NodeType: pathSlice[3],
            }, netMap, subUnit.NetName, subUnit.BizTopic)
        }
    }
    //初始化，准备生成 todo topicid越界检查
    //生成listenPort，为每个主机生成分配tcp端口池
    hostTcpPoolMap := make(map[string]hostTcpUnit)
    var topicID uint16 = 1
    //去重
    ipMap, portMap := make(map[string]int), make(map[int]int)
    //初始化ip、port生成
    _, seedRanges, err := FindIpv4Seeds(ipRange)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("iprange err: %s", err.Error()))
    }
    portRanges, err := sortPorts(portRange)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("portrange err: %s", err.Error()))
    }
    seed, actPort, coverPort := seedRanges[0], portRanges[0], strconv.Itoa(portRanges[0])
    seedIdx, portIdx := 0, 0
    endPoint, overflow := "", false
    //遍历mainMap，根据每个pub查询deploy，构建细分到pub和sub的set的topicInfo
    for _, platInfo := range mainSlice {
        platName := platInfo.PlatForm
        for _, nodeTypeIns := range platInfo.NodeTypes {
            nodeTypeName := nodeTypeIns.NodeType
            SrvIns := nodeTypeIns.SrvFile
            //找部署信息,未部署的集群跳过
            isPubDeployed, pubPi, pubNi := checkDeployment(dplyStructList, platName, nodeTypeName)
            if !isPubDeployed {
                continue
            }
            //整体循环srv的结构，最外层为网络
            for _, netUnit := range SrvIns.PubTopic {
                netName := netUnit.NetName
                //循环分配bizTopic
                for _, bizUnit := range netUnit.BizTopic {
                    if bizUnit.TpcName == "" {
                        return nil, fmt.Errorf("nil bizTopic got under service of: %s.%s", platName, nodeTypeName)
                    }
                    bizTpcName := bizUnit.TpcName
                    if _, ok := netMap[netName]; !ok {
                        continue
                    }
                    for i, _ := range netMap[netName] {
                        if netMap[netName][i].Name == bizTpcName {
                            if netMap[netName][i].IsProc == true {
                                return nil, errors.New(fmt.Sprintf("multi nodeType pubs the same biztopic: %s", bizTpcName))
                            }
                            //biztpc，发送方每个集群只处理一次，在最外层循环
                            for _, pubSetIns := range dplyStructList[pubPi].NodeTypeList[pubNi].SetList {
                                pubTpcInfo := getTpcInfo(topicInfoMap, platName, nodeTypeName, pubSetIns.SetName)
                                //获取endpoint,作为该set的biztpc组播endpoint
                                seedIdx, seed, portIdx, actPort, coverPort, endPoint, overflow = getNextEndPoint(seedIdx, seedRanges, seed, portRanges, portIdx, envNum, actPort, coverPort, ipMap, portMap)
                                if overflow {
                                    return nil, errors.New("endPoint used up")
                                }
                                //处理发送方，添加集群内每个节点的tpcUnit
                                for _, pubNodeIns := range pubSetIns.Deployment.Node {
                                    //检测节点上是否具有该网络
                                    isNetExst, netIdx := checkNetOnNode(netName, pubNodeIns.Network)
                                    if !isNetExst {
                                        return nil, errors.New(fmt.Sprintf("node in set: %s have no adapter for net: %s for topic: %s", util.Join("/", platName, nodeTypeName, pubSetIns.SetName), netName, bizTpcName))
                                    }
                                    pubTpcUnit := ExpTpcTopicUnit{
                                        TopicName:   bizTpcName,
                                        PubCluster:  util.Join(".", platName, nodeTypeName, pubSetIns.SetName),
                                        PubSetId:    pubSetIns.SetID,
                                        PubSetIndex: pubSetIns.SetIndex,
                                        TopicId:     topicID,
                                        EndPoint:    endPoint,
                                        NodeId:      pubNodeIns.NodeId,
                                        NodeIndex:   pubNodeIns.NodeIndex,
                                        IsRMB:       bizUnit.IsRMB,
                                        Net:         pubNodeIns.Network[netIdx],
                                    }
                                    pubTpcInfo.PubExtern.BizTopic = append(pubTpcInfo.PubExtern.BizTopic, pubTpcUnit)
                                }
                                insertTpcInfo(topicInfoMap, pubTpcInfo, platName, nodeTypeName, pubSetIns.SetName)
                                //依次获取接收方节点类型
                                for _, subUnit := range netMap[netName][i].Sub {
                                    isSubDply, subPi, subNi := checkDeployment(dplyStructList, subUnit.Plat, subUnit.NodeType)
                                    if !isSubDply {
                                        continue
                                    }
                                    //依次填充接收方每个set的信息
                                    for _, subSetIns := range dplyStructList[subPi].NodeTypeList[subNi].SetList {
                                        subTpcInfo := getTpcInfo(topicInfoMap, subUnit.Plat, subUnit.NodeType, subSetIns.SetName)
                                        //处理发送方，添加集群内每个节点的tpcUnit
                                        for _, subNodeIns := range subSetIns.Deployment.Node {
                                            //检测节点上是否具有该网络
                                            isNetExst, netIdx := checkNetOnNode(netName, subNodeIns.Network)
                                            if !isNetExst {
                                                return nil, errors.New(fmt.Sprintf("node in set: %s have no adapter for net: %s for topic: %s", util.Join("/", subUnit.Plat, subUnit.NodeType, subSetIns.SetName), netName, bizTpcName))
                                            }
                                            subTpcUnit := ExpTpcTopicUnit{
                                                TopicName:   bizTpcName,
                                                PubCluster:  util.Join(".", platName, nodeTypeName, pubSetIns.SetName),
                                                PubSetId:    pubSetIns.SetID,
                                                PubSetIndex: pubSetIns.SetIndex,
                                                SubCluster:  util.Join(".", subUnit.Plat, subUnit.NodeType, subSetIns.SetName),
                                                SubSetId:    subSetIns.SetID,
                                                SubSetIndex: subSetIns.SetIndex,
                                                TopicId:     topicID,
                                                EndPoint:    endPoint,
                                                NodeId:      subNodeIns.NodeId,
                                                NodeIndex:   subNodeIns.NodeIndex,
                                                IsRMB:       bizUnit.IsRMB,
                                                Net:         subNodeIns.Network[netIdx],
                                            }
                                            subTpcInfo.SubExtern.BizTopic = append(subTpcInfo.SubExtern.BizTopic, subTpcUnit)
                                        }
                                        insertTpcInfo(topicInfoMap, subTpcInfo, subUnit.Plat, subUnit.NodeType, subSetIns.SetName)
                                    }
                                }
                                netMap[netName][i].IsProc = true
                                topicID++
                            }
                        }
                    }
                }
                //循环分配setTpc,无需netMap
                for _, setUnit := range netUnit.SetTopic {
                    setTpcName := setUnit.TpcName
                    //校验是否自发自收
                    if setTpcName == platName+"."+nodeTypeName {
                        return nil, fmt.Errorf("cannot send setTopic to self: %s", setTpcName)
                    }
                    //校验subset是否部署,未部署返回错误
                    sli := strings.Split(setTpcName, ".")
                    if len(sli) != 2 {
                        return nil, errors.New(fmt.Sprintf("err settopic :%s in nodetype %s", setTpcName, nodeTypeName))
                    }
                    isSubDploy, subpi, subni := checkDeployment(dplyStructList, sli[0], sli[1])
                    if !isSubDploy {
                        return nil, errors.New(fmt.Sprintf("sub undeployed:%s", setTpcName))
                    }
                    //为每一对发送方和接收方的set对，分配tpcid，endpoint
                    for _, pubSetIns := range dplyStructList[pubPi].NodeTypeList[pubNi].SetList {
                        for _, subSetIns := range dplyStructList[subpi].NodeTypeList[subni].SetList {
                            //获取endpoint,作为该set对的settpc组播endpoint
                            seedIdx, seed, portIdx, actPort, coverPort, endPoint, overflow = getNextEndPoint(seedIdx, seedRanges, seed, portRanges, portIdx, envNum, actPort, coverPort, ipMap, portMap)
                            if overflow {
                                return nil, errors.New("endPoint used up")
                            }
                            pubTpcInfo := getTpcInfo(topicInfoMap, platName, nodeTypeName, pubSetIns.SetName)
                            subTpcInfo := getTpcInfo(topicInfoMap, sli[0], sli[1], subSetIns.SetName)
                            for _, pubNode := range pubSetIns.Deployment.Node {
                                //检测节点上是否具有该网络
                                isNetExst, netIdx := checkNetOnNode(netName, pubNode.Network)
                                if !isNetExst {
                                    return nil, errors.New(fmt.Sprintf("node in set: %s have no adapter for net: %s for topic: %s", util.Join("/", platName, nodeTypeName, pubSetIns.SetName), netName, setTpcName))
                                }
                                pubTpcInfo.PubExtern.SetTopic = append(pubTpcInfo.PubExtern.SetTopic, ExpTpcTopicUnit{
                                    TopicName:   setTpcName,
                                    PubCluster:  util.Join(".", platName, nodeTypeName, pubSetIns.SetName),
                                    PubSetId:    pubSetIns.SetID,
                                    PubSetIndex: pubSetIns.SetIndex,
                                    SubCluster:  util.Join(".", setTpcName, subSetIns.SetName),
                                    SubSetId:    subSetIns.SetID,
                                    SubSetIndex: subSetIns.SetIndex,
                                    TopicId:     topicID,
                                    EndPoint:    endPoint,
                                    NodeId:      pubNode.NodeId,
                                    NodeIndex:   pubNode.NodeIndex,
                                    IsRMB:       setUnit.IsRMB,
                                    Net:         pubNode.Network[netIdx],
                                })
                            }
                            insertTpcInfo(topicInfoMap, pubTpcInfo, platName, nodeTypeName, pubSetIns.SetName)
                            for _, subNode := range subSetIns.Deployment.Node {
                                //检测节点上是否具有该网络
                                isNetExst, netIdx := checkNetOnNode(netName, subNode.Network)
                                if !isNetExst {
                                    return nil, errors.New(fmt.Sprintf("node in set: %s have no adapter for net: %s for topic: %s", util.Join("/", platName, nodeTypeName, pubSetIns.SetName), netName, setTpcName))
                                }
                                subTpcInfo.SubExtern.SetTopic = append(subTpcInfo.SubExtern.SetTopic, ExpTpcTopicUnit{
                                    TopicName:   setTpcName,
                                    PubCluster:  util.Join(".", platName, nodeTypeName, pubSetIns.SetName),
                                    PubSetId:    pubSetIns.SetID,
                                    PubSetIndex: pubSetIns.SetIndex,
                                    SubCluster:  util.Join(".", setTpcName, subSetIns.SetName),
                                    SubSetId:    subSetIns.SetID,
                                    SubSetIndex: subSetIns.SetIndex,
                                    TopicId:     topicID,
                                    EndPoint:    endPoint,
                                    NodeId:      subNode.NodeId,
                                    NodeIndex:   subNode.NodeIndex,
                                    IsRMB:       setUnit.IsRMB,
                                    Net:         subNode.Network[netIdx],
                                })
                            }
                            insertTpcInfo(topicInfoMap, subTpcInfo, sli[0], sli[1], subSetIns.SetName)
                            topicID++
                        }
                    }
                }
            }
            //以set为单位循环生成
            for _, setIns := range dplyStructList[pubPi].NodeTypeList[pubNi].SetList {
                expTpcMain := getTpcInfo(topicInfoMap, platName, nodeTypeName, setIns.SetName)
                //构建集群内组播通道
                innerNet := SrvIns.InnerTopicNet
                //依次取节点添加组内通道，越界时添加main，其余时间添加follow
                for i := 0; i < len(setIns.Deployment.Node)+1; i++ {
                    seedIdx, seed, portIdx, actPort, coverPort, endPoint, overflow = getNextEndPoint(seedIdx, seedRanges, seed, portRanges, portIdx, envNum, actPort, coverPort, ipMap, portMap)
                    if overflow {
                        return nil, errors.New("endPoint used up")
                    }
                    var topicName, listenPort string
                    var nodeId, nodeIndex uint16
                    var netInfo InfraNetUnit
                    if i < len(setIns.Deployment.Node) {
                        //从节点对应主机的tcp端口池中取出下一个可用端口
                        hostName := setIns.Deployment.Node[i].HostName
                        //若当前主机尚未建立端口池，则新建端口池
                        if _, ok := hostTcpPoolMap[hostName]; !ok {
                            hostTcpPoolMap[hostName] = hostTcpUnit{
                                actualTcpPort: 1024,
                                coverTcpMap:   make(map[int]int),
                            }
                        }
                        actualListenPort := int(hostTcpPoolMap[hostName].actualTcpPort)
                        actualListenPort, coverListenPort, listenPortOverFlow := getNextListenPort(actualListenPort, envNum, hostTcpPoolMap[hostName].coverTcpMap)
                        if listenPortOverFlow {
                            return nil, fmt.Errorf("listenPort used up on host: %s", hostName)
                        }
                        //更新端口池
                        tmp := hostTcpUnit{
                            actualTcpPort: uint16(actualListenPort),
                            coverTcpMap:   hostTcpPoolMap[hostName].coverTcpMap,
                        }
                        hostTcpPoolMap[hostName] = tmp
                        //检测网络是否存在
                        topicName = "follow"
                        listenPort = coverListenPort
                        nodeId, nodeIndex = setIns.Deployment.Node[i].NodeId, setIns.Deployment.Node[i].NodeIndex
                        isNetExt := false
                        for j, _ := range setIns.Deployment.Node[i].Network {
                            if setIns.Deployment.Node[i].Network[j].Name == innerNet {
                                netInfo = setIns.Deployment.Node[i].Network[j]
                                isNetExt = true
                            }
                        }
                        if !isNetExt {
                            return nil, fmt.Errorf("innerTopic net out of host adapter list, path:%s", util.Join("/", platName, nodeTypeName, repository.Service))
                        }
                    } else {
                        topicName = "main"
                        nodeId, nodeIndex = 0, 0
                    }
                    expTpcMain.Inner = append(expTpcMain.Inner, ExpTpcTopicUnit{
                        TopicName:   topicName,
                        PubCluster:  util.Join(".", platName, nodeTypeName, setIns.SetName),
                        PubSetId:    setIns.SetID,
                        PubSetIndex: setIns.SetIndex,
                        ListenPort:  listenPort,
                        TopicId:     topicID,
                        EndPoint:    endPoint,
                        NodeId:      nodeId,
                        NodeIndex:   nodeIndex,
                        Net:         netInfo,
                    })
                    insertTpcInfo(topicInfoMap, expTpcMain, platName, nodeTypeName, setIns.SetName)
                    topicID++
                }
            }
        }
    }
    //校验订阅biztopic的set是否均已处理
    for _, bizVals := range netMap {
        for _, bizVal := range bizVals {
            if !bizVal.IsProc {
                return nil, errors.New(fmt.Sprintf("bizTopic: %s subed but not pub", bizVal.Name))
            }
        }
    }
    return topicInfoMap, nil
}

func MergeSrvStruct(TgtSrv, curtSrv SrvMain) (SrvMain, error) {
    //取pubTopic
    for _, pubTopic := range curtSrv.PubTopic {
        isNetExist := false
        netIdx := 0
        for i, stdPubTopic := range TgtSrv.PubTopic {
            if pubTopic.NetName == stdPubTopic.NetName {
                netIdx = i
                isNetExist = true
                break
            }
        }
        //当前网络信息不存在，直接拷贝
        if !isNetExist {
            TgtSrv.PubTopic = append(TgtSrv.PubTopic, pubTopic)
        } else {
            //合并bizUnit
            for _, bizUnit := range pubTopic.BizTopic {
                isBizUnitExist := false
                for _, stdBizUnit := range TgtSrv.PubTopic[netIdx].BizTopic {
                    if bizUnit.TpcName == stdBizUnit.TpcName {
                        if bizUnit.IsRMB != stdBizUnit.IsRMB {
                            return TgtSrv, errors.New(fmt.Sprintf("different isRMB value of pub biztopic name: %s", stdBizUnit.TpcName))
                        }
                        isBizUnitExist = true
                        break
                    }
                }
                if isBizUnitExist {
                    continue
                } else {
                    //添加尚未录入的bizUnit
                    TgtSrv.PubTopic[netIdx].BizTopic = append(TgtSrv.PubTopic[netIdx].BizTopic, bizUnit)
                }
            }
            //合并setUnit
            for _, setUnit := range pubTopic.SetTopic {
                isSetUnitExist := false
                for _, stdSetUnit := range TgtSrv.PubTopic[netIdx].SetTopic {
                    if setUnit.TpcName == stdSetUnit.TpcName {
                        if setUnit.IsRMB != stdSetUnit.IsRMB {
                            return TgtSrv, errors.New(fmt.Sprintf("different isRMB value of pub settopic name: %s", stdSetUnit.TpcName))
                        }
                        isSetUnitExist = true
                        break
                    }
                }
                if isSetUnitExist {
                    continue
                } else {
                    //添加尚未录入的bizUnit
                    TgtSrv.PubTopic[netIdx].SetTopic = append(TgtSrv.PubTopic[netIdx].SetTopic, setUnit)
                }
            }
        }
    }
    //取subTopic
    for _, subTopic := range curtSrv.SubTopic {
        isNetExist := false
        netIdx := 0
        for i, stdSubTopic := range TgtSrv.SubTopic {
            if subTopic.NetName == stdSubTopic.NetName {
                netIdx = i
                isNetExist = true
                break
            }
        }
        //当前网络信息不存在，直接拷贝
        if !isNetExist {
            TgtSrv.SubTopic = append(TgtSrv.SubTopic, subTopic)
        } else {
            //合并bizUnit
            for _, bizUnit := range subTopic.BizTopic {
                isBizUnitExist := false
                for _, stdBizUnit := range TgtSrv.SubTopic[netIdx].BizTopic {
                    if bizUnit.TpcName == stdBizUnit.TpcName {
                        if bizUnit.IsRMB != stdBizUnit.IsRMB {
                            return TgtSrv, errors.New(fmt.Sprintf("different isRMB value of sub biztopic name: %s", stdBizUnit.TpcName))
                        }
                        isBizUnitExist = true
                        break
                    }
                }
                if isBizUnitExist {
                    continue
                } else {
                    //添加尚未录入的bizUnit
                    TgtSrv.SubTopic[netIdx].BizTopic = append(TgtSrv.SubTopic[netIdx].BizTopic, bizUnit)
                }
            }
            //合并setUnit
            for _, setUnit := range subTopic.SetTopic {
                isSetUnitExist := false
                for _, stdSetUnit := range TgtSrv.SubTopic[netIdx].SetTopic {
                    if setUnit.TpcName == stdSetUnit.TpcName {
                        if setUnit.IsRMB != stdSetUnit.IsRMB {
                            return TgtSrv, errors.New(fmt.Sprintf("different isRMB value of sub settopic name: %s", stdSetUnit.TpcName))
                        }
                        isSetUnitExist = true
                        break
                    }
                }
                if isSetUnitExist {
                    continue
                } else {
                    //添加尚未录入的bizUnit
                    TgtSrv.SubTopic[netIdx].SetTopic = append(TgtSrv.SubTopic[netIdx].SetTopic, setUnit)
                }
            }
        }
        /*if len(subTopic.BizTopic) > 0 {
            //向netMap中添加信息
            MergeNetMap(inputIdtfy, netMap, subTopic.NetName, subTopic.BizTopic)
        }*/
    }
    return TgtSrv, nil
}

func MergeNetMap(inputIdtfy idtfy, netMap map[string][]topicStat, net string, bizTopics []SrvTpcStatUnit) {
    if len(bizTopics) == 0 {
        return
    }
    if _, ok := netMap[net]; !ok {
        netMap[net] = make([]topicStat, 0)
        for _, bizUnit := range bizTopics {
            netMap[net] = append(netMap[net], topicStat{
                Name:   bizUnit.TpcName,
                IsProc: false,
                Sub:    []idtfy{inputIdtfy},
            })
        }
    } else {
        for _, bizUnit := range bizTopics {
            isTpcRpt := false
            tpcIdx := 0
            for i, topicStatIns := range netMap[net] {
                if topicStatIns.Name == bizUnit.TpcName {
                    isTpcRpt = true
                    tpcIdx = i
                    break
                }
            }
            if !isTpcRpt {
                netMap[net] = append(netMap[net], topicStat{
                    IsProc: false,
                    Name:   bizUnit.TpcName,
                    Sub:    []idtfy{inputIdtfy},
                })
            } else {
                for _, sub := range netMap[net][tpcIdx].Sub {
                    if sub.Plat == inputIdtfy.Plat && sub.NodeType == inputIdtfy.NodeType {
                        return
                    }
                }
                netMap[net][tpcIdx].Sub = append(netMap[net][tpcIdx].Sub, inputIdtfy)
            }
        }
    }
}

// FindIpv4Seeds 寻找多个给定Ipv4范围的起始位置，返回起始值。此处假定输入参数和Ip值都合法。
func FindIpv4Seeds(ipRanges []string) ([]string, []int32, error) {
    if len(ipRanges) == 0 || len(ipRanges)%2 != 0 {
        return nil, nil, errors.New("nil or odd ipRange input, Please checkout")
    }
    seedRanges := make([]int32, 0)
    for i := 0; 2*i < len(ipRanges); i++ {
        if !isIpv4Legal(ipRanges[2*i]) || !isIpv4Legal(ipRanges[2*i+1]) {
            return nil, nil, errors.New(fmt.Sprintf("err ipv4 pair of %s, %s", ipRanges[2*i], ipRanges[2*i+1]))
        }
        frontSliceStr, backSliceStr := strings.Split(ipRanges[2*i], "."), strings.Split(ipRanges[2*i+1], ".")
        frontSlice, backSlice := make([]int32, 0), make([]int32, 0)
        for i, _ := range frontSliceStr {
            intFrt, _ := strconv.Atoi(frontSliceStr[i])
            intBck, _ := strconv.Atoi(backSliceStr[i])
            frontSlice = append(frontSlice, int32(intFrt))
            backSlice = append(backSlice, int32(intBck))
        }
        frontLess := true
        for i, _ := range frontSlice {
            if frontSlice[i] > backSlice[i] {
                frontLess = false
                break
            }
        }
        if !frontLess {
            tmp := ipRanges[2*i]
            ipRanges[2*i] = ipRanges[2*i+1]
            ipRanges[2*i+1] = tmp
        }
        frtSeed, err := checkAndGetSeed(ipRanges[2*i])
        if err != nil {
            return nil, nil, err
        }
        bckSeed, err := checkAndGetSeed(ipRanges[2*i+1])
        if err != nil {
            return nil, nil, err
        }
        seedRanges = append(seedRanges, frtSeed, bckSeed)
    }
    return ipRanges, seedRanges, nil
}

// GetNextIpv4 将种子值+1，返回生成的ip和新种子，并指示是否溢出
func GetNextIpv4(idx int, seedRanges []int32, oldSeed int32) (int, string, int32, bool) {
    var newSeed int32 = 0
    if seedRanges[2*idx+1] >= seedRanges[2*idx] {
        if oldSeed > seedRanges[2*idx+1] {
            if 2*(idx+1) >= len(seedRanges) {
                return 0, "", 0, true
            } else {
                idx++
                newSeed = seedRanges[2*idx]
            }
        } else {
            newSeed = oldSeed
        }
    } else {
        if oldSeed > seedRanges[2*idx+1] && oldSeed < seedRanges[2*idx] {
            if 2*(idx+1) >= len(seedRanges) {
                return 0, "", 0, true
            } else {
                idx++
                newSeed = seedRanges[2*idx]
            }
        } else {
            newSeed = oldSeed
        }
    }
    b1 := (newSeed >> 24) & 0xff
    b2 := (newSeed >> 16) & 0xff
    b3 := (newSeed >> 8) & 0xff
    b4 := (newSeed) & 0xff
    newSeed++
    return idx, strconv.Itoa(int(b1)) + "." + strconv.Itoa(int(b2)) + "." + strconv.Itoa(int(b3)) + "." + strconv.Itoa(int(b4)), newSeed, false
}

func isIpv4Legal(input string) bool {
    ipv4Slice := strings.Split(input, ".")
    if len(ipv4Slice) != 4 {
        return false
    }
    for _, strPart := range ipv4Slice {
        i, err := strconv.Atoi(strPart)
        if err != nil || i < 0 || i > 255 {
            return false
        }
    }
    return true
}

func checkAndGetSeed(ipv4String string) (int32, error) {
    if !isIpv4Legal(ipv4String) {
        return 0, errors.New(fmt.Sprintf("err ipv4 pair of %s", ipv4String))
    }
    strSlice := strings.Split(ipv4String, ".")
    var bits []int32
    for i, _ := range strSlice {
        b, err := strconv.ParseInt(strSlice[i], 10, 32)
        if err != nil {
            return 0, err
        }
        bits = append(bits, int32(b))
    }
    var res int32 = 0
    Pos := true
    for i, bit := range bits {
        switch i {
        case 0:
            if bit < 127 {
                res = 0
                res = res | (bit << 24)
            } else {
                Pos = false
                res = res + (-1 & (bit << 24))
            }
        case 1:
            if Pos {
                res = res | (bit << 16)
            } else {
                res = res + (-1 & (bit << 16))
            }

        case 2:
            if Pos {
                res = res | (bit << 8)
            } else {
                res = res + (-1 & (bit << 8))
            }
        case 3:
            if Pos {
                res = res | bit
            } else {
                res = res + (-1 & bit)
            }
        }
    }
    return res, nil
}

func checkPort(strPort string) (int, bool) {
    port, err := strconv.Atoi(strPort)
    if err != nil || port < 1024 || port > 65535 {
        return 0, false
    }
    return port, true
}

func sortPorts(portRanges []string) ([]int, error) {
    if len(portRanges) == 0 || len(portRanges)%2 != 0 {
        return nil, errors.New("nil or odd portRange input, Please checkout")
    }
    var ports []int
    for i := 0; 2*i < len(portRanges); i++ {
        p1, check1 := checkPort(portRanges[2*i])
        p2, check2 := checkPort(portRanges[2*i+1])
        if !check1 || !check2 || p1 == p2 {
            return nil, errors.New(fmt.Sprintf("err portRange of: %s, %s", portRanges[2*i], portRanges[2*i+1]))
        }
        if p1 > p2 {
            ports = append(ports, p2, p1)
        } else {
            ports = append(ports, p1, p2)
        }
    }
    return ports, nil
}

//假设envNum合法
func getNextPort(ports []int, idx int, envNum string, actualPort int) (int, int, int, bool) {
    if actualPort > ports[2*idx+1] {
        if 2*(idx+1) >= len(ports) {
            return 0, 0, 0, true
        }
        idx++
        actualPort = ports[2*idx]
    }
    env, _ := strconv.Atoi(envNum)
    env = env * 10
    coverPort := actualPort/1000*1000 + env + actualPort%10
    return idx, actualPort + 1, coverPort, false
}

//假设envNum合法
func getNextListenPort(inputPort int, envNum string, listenPortMap map[int]int) (actualListenPort int, coverPort string, overFlow bool) {
    env, _ := strconv.Atoi(envNum)
    env = env * 10
    for true {
        if inputPort > 65535 {
            return 0, "", true
        }
        coverPortInt := inputPort/1000*1000 + env + inputPort%10
        if coverPortInt < 1024 || coverPortInt > 65535 {
            inputPort++
            continue
        }
        if _, ok := listenPortMap[coverPortInt]; !ok {
            listenPortMap[coverPortInt] = 0
            actualListenPort = inputPort + 1
            coverPort = strconv.Itoa(coverPortInt)
            break
        }
        inputPort++
    }
    return
}

func getNextEndPoint(idxI int, seedRanges []int32, oldSeed int32, ports []int, idxP int, envNum string, actualPort int, coverPortIn string, ipMap map[string]int, portMap map[int]int) (int, int32, int, int, string, string, bool) {
    isIpRpt, isPortRpt, isIpOverFlow, isPortOverflow := true, true, false, false
    ipv4, res, coverPort := "", "", 0
    for isIpRpt {
        idxI, ipv4, oldSeed, isIpOverFlow = GetNextIpv4(idxI, seedRanges, oldSeed)
        if isIpOverFlow {
            break
        }
        if _, ok := ipMap[ipv4]; !ok {
            ipMap[ipv4] = 0
            break
        }
    }
    if !isIpOverFlow {
        for true {
            idxP, actualPort, coverPort, isPortOverflow = getNextPort(ports, idxP, envNum, actualPort)
            if coverPort < 65536 {
                portMap[coverPort] = 0
                break
            }
        }
    } else {
        for isPortRpt {
            //actualPort++
            idxP, actualPort, coverPort, isPortOverflow = getNextPort(ports, idxP, envNum, actualPort)
            if isPortOverflow {
                //溢出
                return 0, 0, 0, 0, "", "", true
            }
            if coverPort > 65535 {
                continue
            }
            if _, ok := portMap[coverPort]; !ok {
                portMap[coverPort] = 0
                break
            }
        }
        for k, _ := range ipMap {
            delete(ipMap, k)
        }
        idxI = 0
        oldSeed = seedRanges[0]
        idxI, ipv4, oldSeed, isIpOverFlow = GetNextIpv4(idxI, seedRanges, oldSeed)
        if _, ok := ipMap[ipv4]; !ok {
            ipMap[ipv4] = 0
        }
    }
    //在已排除溢出的情况下，保持现有port
    actualPort--
    coverPortIn = strconv.Itoa(coverPort)
    res = ipv4 + ":" + coverPortIn
    return idxI, oldSeed, idxP, actualPort, coverPortIn, res, false
}

func checkDeployment(dplyStructList []ChartDeployMain, platName, nodeTypeName string) (bool, int, int) {
    dpi, dni := 0, 0
    isDeployed := false
    for i, _ := range dplyStructList {
        if dplyStructList[i].Platform == platName {
            for j, _ := range dplyStructList[i].NodeTypeList {
                if dplyStructList[i].NodeTypeList[j].NodeType == nodeTypeName {
                    dpi = i
                    dni = j
                    isDeployed = true
                    break
                }
            }
            if isDeployed {
                break
            }
        }
    }
    return isDeployed, dpi, dni
}

func getTpcInfo(topicInfoMap map[string]map[string]map[string]ExpTpcMain, plat, nodeType, set string) ExpTpcMain {
    if _, ok1 := topicInfoMap[plat]; ok1 {
        if _, ok2 := topicInfoMap[plat][nodeType]; ok2 {
            if _, ok3 := topicInfoMap[plat][nodeType][set]; ok3 {
                return topicInfoMap[plat][nodeType][set]
            }
        }
    }
    return ExpTpcMain{}
}

func insertTpcInfo(topicInfoMapIn map[string]map[string]map[string]ExpTpcMain, tpcInfo ExpTpcMain, plat, nodeType, set string) {
    pubTpcInfoExst1, pubTpcInfoExst2 := false, false
    if _, ok1 := topicInfoMapIn[plat]; ok1 {
        pubTpcInfoExst1 = true
        if _, ok2 := topicInfoMapIn[plat][nodeType]; ok2 {
            pubTpcInfoExst2 = true
        }
    }
    if pubTpcInfoExst1 && pubTpcInfoExst2 {
        topicInfoMapIn[plat][nodeType][set] = tpcInfo
    }
    if pubTpcInfoExst1 && !pubTpcInfoExst2 {
        topicInfoMapIn[plat][nodeType] = make(map[string]ExpTpcMain)
        topicInfoMapIn[plat][nodeType][set] = tpcInfo
    }
    if !pubTpcInfoExst1 && !pubTpcInfoExst2 {
        topicInfoMapIn[plat] = make(map[string]map[string]ExpTpcMain)
        topicInfoMapIn[plat][nodeType] = make(map[string]ExpTpcMain)
        topicInfoMapIn[plat][nodeType][set] = tpcInfo
    }
}

func checkNetOnNode(tgtNet string, netOnNode []InfraNetUnit) (bool, int) {
    for i, _ := range netOnNode {
        if tgtNet == netOnNode[i].Name {
            return true, i
        }
    }
    return false, 0
}
