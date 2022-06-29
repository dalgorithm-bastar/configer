package generation

import (
    "reflect"
    "testing"
)

func TestFindIpv4Seeds(t *testing.T) {
    type args struct {
        ipRanges []string
    }
    tests := []struct {
        name    string
        args    args
        want    []string
        want1   []int32
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, err := FindIpv4Seeds(tt.args.ipRanges)
            if (err != nil) != tt.wantErr {
                t.Errorf("FindIpv4Seeds() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("FindIpv4Seeds() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("FindIpv4Seeds() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}

func TestGenerateTopicInfo(t *testing.T) {
    type args struct {
        dplyStructList []ChartDeployMain
        rawSlice       []RawFile
        envNum         string
        topicIpRange   []string
        topicPortRange []string
        tcpPortRange   []string
    }
    tests := []struct {
        name    string
        args    args
        want    map[string]map[string]map[string]ExpTpcMain
        want1   map[string]hostTcpUnit
        want2   []int
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, err := GenerateTopicInfo(tt.args.dplyStructList, tt.args.rawSlice, tt.args.envNum, tt.args.topicIpRange, tt.args.topicPortRange, tt.args.tcpPortRange)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateTopicInfo() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GenerateTopicInfo() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("GenerateTopicInfo() got1 = %v, want %v", got1, tt.want1)
            }
            if !reflect.DeepEqual(got2, tt.want2) {
                t.Errorf("GenerateTopicInfo() got2 = %v, want %v", got2, tt.want2)
            }
        })
    }
}

func TestGetNextIpv4(t *testing.T) {
    type args struct {
        idx        int
        seedRanges []int32
        oldSeed    int32
    }
    tests := []struct {
        name  string
        args  args
        want  int
        want1 string
        want2 int32
        want3 bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := GetNextIpv4(tt.args.idx, tt.args.seedRanges, tt.args.oldSeed)
            if got != tt.want {
                t.Errorf("GetNextIpv4() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("GetNextIpv4() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("GetNextIpv4() got2 = %v, want %v", got2, tt.want2)
            }
            if got3 != tt.want3 {
                t.Errorf("GetNextIpv4() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func TestMergeNetMap(t *testing.T) {
    type args struct {
        inputIdtfy idtfy
        netMap     map[string][]topicStat
        net        string
        bizTopics  []SrvTpcStatUnit
    }
    tests := []struct {
        name string
        args args
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            MergeNetMap(tt.args.inputIdtfy, tt.args.netMap, tt.args.net, tt.args.bizTopics)
        })
    }
}

func TestMergeSrvStruct(t *testing.T) {
    tgtSrv := SrvMain{
        InnerTopicNet: "innerNet",
        PubTopic: []SrvStatementUnit{{
            NetName: _bizNet,
            BizTopic: []SrvTpcStatUnit{{
                TpcName: "biz1",
                IsRMB:   1,
            }},
            SetTopic: []SrvTpcStatUnit{{
                TpcName: "set1",
                IsRMB:   1,
            }},
        }},
        SubTopic: []SrvStatementUnit{{
            NetName: _bizNet,
            BizTopic: []SrvTpcStatUnit{{
                TpcName: "biz2",
                IsRMB:   1,
            }},
        }},
    }
    m1Srv := SrvMain{
        InnerTopicNet: "innerNet",
        PubTopic: []SrvStatementUnit{{
            NetName: _bizNet,
            BizTopic: []SrvTpcStatUnit{{
                TpcName: "biz1",
                IsRMB:   1,
            }, {
                TpcName: "biz3",
                IsRMB:   1,
            }},
            SetTopic: []SrvTpcStatUnit{{
                TpcName: "set1",
                IsRMB:   1,
            }, {
                TpcName: "set2",
                IsRMB:   1,
            }},
        }},
        SubTopic: []SrvStatementUnit{{
            NetName: _bizNet,
            BizTopic: []SrvTpcStatUnit{{
                TpcName: "biz2",
                IsRMB:   1,
            }, {
                TpcName: "biz4",
                IsRMB:   1,
            }},
        }},
    }
    res1Srv := SrvMain{
        InnerTopicNet: "innerNet",
        PubTopic: []SrvStatementUnit{
            {
                NetName: _bizNet,
                BizTopic: []SrvTpcStatUnit{
                    {
                        TpcName: "biz1",
                        IsRMB:   1,
                    },
                    {
                        TpcName: "biz3",
                        IsRMB:   1,
                    },
                },
                SetTopic: []SrvTpcStatUnit{
                    {
                        TpcName: "set1",
                        IsRMB:   1,
                    },
                    {
                        TpcName: "set2",
                        IsRMB:   1,
                    },
                },
            },
        },
        SubTopic: []SrvStatementUnit{
            {
                NetName: _bizNet,
                BizTopic: []SrvTpcStatUnit{
                    {
                        TpcName: "biz2",
                        IsRMB:   1,
                    },
                    {
                        TpcName: "biz4",
                        IsRMB:   1,
                    },
                },
            },
        },
    }
    type args struct {
        TgtSrv  SrvMain
        curtSrv SrvMain
    }
    tests := []struct {
        name    string
        args    args
        want    SrvMain
        wantErr bool
    }{
        {
            name: "merge biz and set",
            args: args{
                TgtSrv:  tgtSrv,
                curtSrv: m1Srv,
            },
            wantErr: false,
            want:    res1Srv,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := MergeSrvStruct(tt.args.TgtSrv, tt.args.curtSrv)
            if (err != nil) != tt.wantErr {
                t.Errorf("MergeSrvStruct() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("MergeSrvStruct() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_checkAndGetSeed(t *testing.T) {
    type args struct {
        ipv4String string
    }
    tests := []struct {
        name    string
        args    args
        want    int32
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := checkAndGetSeed(tt.args.ipv4String)
            if (err != nil) != tt.wantErr {
                t.Errorf("checkAndGetSeed() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("checkAndGetSeed() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_checkDeployment(t *testing.T) {
    type args struct {
        dplyStructList []ChartDeployMain
        platName       string
        nodeTypeName   string
    }
    tests := []struct {
        name  string
        args  args
        want  bool
        want1 int
        want2 int
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2 := checkDeployment(tt.args.dplyStructList, tt.args.platName, tt.args.nodeTypeName)
            if got != tt.want {
                t.Errorf("checkDeployment() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("checkDeployment() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("checkDeployment() got2 = %v, want %v", got2, tt.want2)
            }
        })
    }
}

func Test_checkNetOnNode(t *testing.T) {
    type args struct {
        tgtNet    string
        netOnNode []InfraNetUnit
    }
    tests := []struct {
        name  string
        args  args
        want  bool
        want1 int
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1 := checkNetOnNode(tt.args.tgtNet, tt.args.netOnNode)
            if got != tt.want {
                t.Errorf("checkNetOnNode() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("checkNetOnNode() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}

func Test_checkPort(t *testing.T) {
    type args struct {
        strPort string
    }
    tests := []struct {
        name  string
        args  args
        want  int
        want1 bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1 := checkPort(tt.args.strPort)
            if got != tt.want {
                t.Errorf("checkPort() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("checkPort() got1 = %v, want %v", got1, tt.want1)
            }
        })
    }
}

func Test_getNextEndPoint(t *testing.T) {
    type args struct {
        idxI        int
        seedRanges  []int32
        oldSeed     int32
        ports       []int
        idxP        int
        envNum      string
        actualPort  int
        coverPortIn string
        ipMap       map[string]int
        portMap     map[int]int
    }
    tests := []struct {
        name  string
        args  args
        want  int
        want1 int32
        want2 int
        want3 int
        want4 string
        want5 string
        want6 bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3, got4, got5, got6 := getNextEndPoint(tt.args.idxI, tt.args.seedRanges, tt.args.oldSeed, tt.args.ports, tt.args.idxP, tt.args.envNum, tt.args.actualPort, tt.args.coverPortIn, tt.args.ipMap, tt.args.portMap)
            if got != tt.want {
                t.Errorf("getNextEndPoint() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("getNextEndPoint() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getNextEndPoint() got2 = %v, want %v", got2, tt.want2)
            }
            if got3 != tt.want3 {
                t.Errorf("getNextEndPoint() got3 = %v, want %v", got3, tt.want3)
            }
            if got4 != tt.want4 {
                t.Errorf("getNextEndPoint() got4 = %v, want %v", got4, tt.want4)
            }
            if got5 != tt.want5 {
                t.Errorf("getNextEndPoint() got5 = %v, want %v", got5, tt.want5)
            }
            if got6 != tt.want6 {
                t.Errorf("getNextEndPoint() got6 = %v, want %v", got6, tt.want6)
            }
        })
    }
}

func Test_getNextListenPort(t *testing.T) {
    type args struct {
        ports         []int
        idxin         int
        inputPort     int
        envNum        string
        listenPortMap map[int]int
    }
    tests := []struct {
        name                 string
        args                 args
        wantActualListenPort int
        wantCoverPort        string
        wantIdx              int
        wantOverFlow         bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gotActualListenPort, gotCoverPort, gotIdx, gotOverFlow := getNextListenPort(tt.args.ports, tt.args.idxin, tt.args.inputPort, tt.args.envNum, tt.args.listenPortMap)
            if gotActualListenPort != tt.wantActualListenPort {
                t.Errorf("getNextListenPort() gotActualListenPort = %v, want %v", gotActualListenPort, tt.wantActualListenPort)
            }
            if gotCoverPort != tt.wantCoverPort {
                t.Errorf("getNextListenPort() gotCoverPort = %v, want %v", gotCoverPort, tt.wantCoverPort)
            }
            if gotIdx != tt.wantIdx {
                t.Errorf("getNextListenPort() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
            }
            if gotOverFlow != tt.wantOverFlow {
                t.Errorf("getNextListenPort() gotOverFlow = %v, want %v", gotOverFlow, tt.wantOverFlow)
            }
        })
    }
}

func Test_getNextPort(t *testing.T) {
    type args struct {
        ports      []int
        idx        int
        envNum     string
        actualPort int
    }
    tests := []struct {
        name  string
        args  args
        want  int
        want1 int
        want2 int
        want3 bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getNextPort(tt.args.ports, tt.args.idx, tt.args.envNum, tt.args.actualPort)
            if got != tt.want {
                t.Errorf("getNextPort() got = %v, want %v", got, tt.want)
            }
            if got1 != tt.want1 {
                t.Errorf("getNextPort() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getNextPort() got2 = %v, want %v", got2, tt.want2)
            }
            if got3 != tt.want3 {
                t.Errorf("getNextPort() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getTpcInfo(t *testing.T) {
    type args struct {
        topicInfoMap map[string]map[string]map[string]ExpTpcMain
        plat         string
        nodeType     string
        set          string
    }
    tests := []struct {
        name string
        args args
        want ExpTpcMain
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := getTpcInfo(tt.args.topicInfoMap, tt.args.plat, tt.args.nodeType, tt.args.set); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getTpcInfo() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_insertTpcInfo(t *testing.T) {
    type args struct {
        topicInfoMapIn map[string]map[string]map[string]ExpTpcMain
        tpcInfo        ExpTpcMain
        plat           string
        nodeType       string
        set            string
    }
    tests := []struct {
        name string
        args args
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            insertTpcInfo(tt.args.topicInfoMapIn, tt.args.tpcInfo, tt.args.plat, tt.args.nodeType, tt.args.set)
        })
    }
}

func Test_isIpv4Legal(t *testing.T) {
    type args struct {
        input string
    }
    tests := []struct {
        name string
        args args
        want bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := isIpv4Legal(tt.args.input); got != tt.want {
                t.Errorf("isIpv4Legal() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_sortPorts(t *testing.T) {
    type args struct {
        portRanges []string
    }
    tests := []struct {
        name    string
        args    args
        want    []int
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := sortPorts(tt.args.portRanges)
            if (err != nil) != tt.wantErr {
                t.Errorf("sortPorts() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("sortPorts() got = %v, want %v", got, tt.want)
            }
        })
    }
}
