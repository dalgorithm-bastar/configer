package generation

/****************** 基础设施信息 ******************/

type InfraNetUnit struct {
    Name    string `json:"name"`
    Adaptor string `json:"adapter"`
    Ipv4    string `json:"ipv4"`
    Ipv6    string `json:"ipv6"`
}

type InfraUserUnit struct {
    Name string `json:"name"`
}

type InfraHostUnit struct {
    NodeId    uint16          `json:"nodeId"`
    NodeIndex uint16          `json:"nodeIndex"`
    HostName  string          `json:"hostName"`
    Location  string          `json:"location"`
    Network   []InfraNetUnit  `json:"network"`
    User      []InfraUserUnit `json:"user"`
}

type InfraMain struct {
    Host []InfraHostUnit `json:"host"`
}

/********************************************/

/**************** 部署信息 *******************/

type DeployArtificialUnit struct {
    Url        string `json:"url"`
    DeployPath string `json:"deployPath"`
    ConfigPath string `json:"configPath"`
}

type DeployMain struct {
    Artificial DeployArtificialUnit `json:"artifact"`
    Node       []InfraHostUnit      `json:"node"`
}

/********************************************/

/****************** 服务信息 ******************/

type SrvTpcStatUnit struct {
    TpcName string `json:"topicName"`
    IsRMB   uint16 `json:"isRMB"`
}

type SrvStatementUnit struct {
    NetName  string           `json:"netName"`
    BizTopic []SrvTpcStatUnit `json:"bizTopic"`
    SetTopic []SrvTpcStatUnit `json:"setTopic"`
}

type SrvMain struct {
    InnerTopicNet string             `json:"innerTopicNet"`
    PubTopic      []SrvStatementUnit `json:"pubTopic"`
    SubTopic      []SrvStatementUnit `json:"subTopic"`
}

/* 默认填充模板 */

type FillArgs struct {
    PlatName  string
    NodeType  string
    SetId     uint16
    SetIndex  uint16
    SetName   string
    NodeId    uint16
    NodeIndex uint16
}

/********************************************/

/************ 扩展组播信息相关 ************/

type ExpTpcTopicUnit struct {
    TopicName   string       `json:"topicName"`
    PubCluster  string       `json:"pubCluster"`
    PubSetId    uint16       `json:"pubSetId"`
    PubSetIndex uint16       `json:"pubSetIndex"`
    SubCluster  string       `json:"subCluster"`
    SubSetId    uint16       `json:"subSetId"`
    SubSetIndex uint16       `json:"subSetIndex"`
    ListenPort  string       `json:"listenPort"`
    TopicId     uint16       `json:"topicId"`
    EndPoint    string       `json:"endPoint"`
    NodeId      uint16       `json:"nodeId"`
    NodeIndex   uint16       `json:"nodeIndex"`
    IsRMB       uint16       `json:"isRMB"`
    Net         InfraNetUnit `json:"net"`
}

type ExpTpcInteractList struct {
    BizTopic []ExpTpcTopicUnit `json:"biz_topic"`
    SetTopic []ExpTpcTopicUnit `json:"set_topic"`
}

type ExpTpcMain struct {
    Inner     []ExpTpcTopicUnit  `json:"inner"`
    PubExtern ExpTpcInteractList `json:"pub_extern"`
    SubExtern ExpTpcInteractList `json:"sub_extern"`
}

/********************************************/

/****************** 部署信息总表 ******************/

type ChartDeploySet struct {
    SetName    string     `json:"set_name"`
    SetID      uint16     `json:"set_id"`
    SetIndex   uint16     `json:"set_index"`
    Deployment DeployMain `json:"deployment"`
}

type ChartDeployNodeType struct {
    NodeType string           `json:"node_type"`
    SetList  []ChartDeploySet `json:"set_list"`
}

type ChartDeployMain struct {
    Platform     string                `json:"platform"`
    NodeTypeList []ChartDeployNodeType `json:"node_type_list"`
}

type ChartDeploy struct {
    Scheme    string            `json:"scheme"`
    Platforms []ChartDeployMain `json:"platforms"`
}

/**********************************************/

/****************** 组播信息总表 ******************/

type ChartTpcSet struct {
    SetName   string     `json:"set_name"`
    SetID     uint16     `json:"set_id"`
    BroadInfo ExpTpcMain `json:"broad_info"`
}

type ChartTpcNodeType struct {
    NodeType string        `json:"node_type"`
    SetList  []ChartTpcSet `json:"set_list"`
}

type ChartTpcMain struct {
    Platform     string             `json:"platform"`
    NodeTypeList []ChartTpcNodeType `json:"node_type_list"`
}

type ChartTpc struct {
    Scheme    string         `json:"scheme"`
    Platforms []ChartTpcMain `json:"platforms"`
}

/**********************************************/
