package generation

/****************** 基础设施信息 ******************/

type InfraNetUnit struct {
	Name    string `json:"name" yaml:"name,omitempty"`
	Adaptor string `json:"adapter" yaml:"adapter,omitempty"`
	Ipv4    string `json:"ipv4" yaml:"ipv4,omitempty"`
	Ipv6    string `json:"ipv6" yaml:"ipv6,omitempty"`
}

type InfraUserUnit struct {
	Name string `json:"name" yaml:"name,omitempty"`
}

type InfraHostUnit struct {
	NodeId    uint16          `json:"nodeId" yaml:"nodeId,omitempty"`
	NodeIndex uint16          `json:"nodeIndex" yaml:"nodeIndex,omitempty"`
	HostName  string          `json:"hostName" yaml:"hostName,omitempty"`
	Location  string          `json:"location" yaml:"location,omitempty"`
	Network   []InfraNetUnit  `json:"network" yaml:"network,omitempty"`
	User      []InfraUserUnit `json:"user" yaml:"user,omitempty"`
}

type InfraMain struct {
	Host []InfraHostUnit `json:"host" yaml:"host,omitempty"`
}

/********************************************/

/**************** 部署信息 *******************/

type DeployArtifactUnit struct {
	Url        string `json:"url" yaml:"url,omitempty"`
	DeployPath string `json:"deployPath" yaml:"deployPath,omitempty"`
	ConfigPath string `json:"configPath" yaml:"configPath,omitempty"`
}

type DeployMain struct {
	Artifact DeployArtifactUnit `json:"artifact" yaml:"artifact,omitempty"`
	UserName string             `json:"userName" yaml:"userName,omitempty"`
	SetID    uint16             `json:"setID" yaml:"setID,omitempty"`
	Node     []InfraHostUnit    `json:"node" yaml:"node,omitempty"`
}

/********************************************/

/****************** 服务信息 ******************/

type SrvTpcStatUnit struct {
	TpcName   string  `json:"topicName" yaml:"topicName,omitempty"`
	IsRMB     *uint16 `json:"isRMB" yaml:"isRMB,omitempty"`
	UnaidedIO *bool   `json:"unaidedIO" yaml:"unaidedIO,omitempty"`
	Priority  *uint16 `json:"priority" yaml:"priority,omitempty"`
}

type SrvStatementUnit struct {
	NetName  string           `json:"netName" yaml:"netName,omitempty"`
	BizTopic []SrvTpcStatUnit `json:"bizTopic" yaml:"bizTopic,omitempty"`
	SetTopic []SrvTpcStatUnit `json:"setTopic" yaml:"setTopic,omitempty"`
}

type SrvMain struct {
	InnerTopicNet string             `json:"innerTopicNet" yaml:"innerTopicNet,omitempty"`
	PubTopic      []SrvStatementUnit `json:"pubTopic" yaml:"pubTopic,omitempty"`
	SubTopic      []SrvStatementUnit `json:"subTopic" yaml:"subTopic,omitempty"`
}

/* 默认填充模板 */

type FillArgs struct {
	PlatName  string
	NodeType  string
	NodeNum   uint16
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
	IsRMB       *uint16      `json:"isRMB,omitempty"`
	UnaidedIO   *bool        `json:"unaidedIO,omitempty"`
	Priority    *uint16      `json:"priority,omitempty"`
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
	SetName    string     `json:"set_name" yaml:"set_name"`
	SetID      uint16     `json:"set_id" yaml:"set_id"`
	SetIndex   uint16     `json:"set_index" yaml:"set_index"`
	Deployment DeployMain `json:"deployment" yaml:"deployment"`
}

type ChartDeployNodeType struct {
	NodeType string           `json:"node_type" yaml:"node_type"`
	SetList  []ChartDeploySet `json:"set_list" yaml:"set_list"`
}

type ChartDeployMain struct {
	Platform     string                `json:"platform" yaml:"platform"`
	NodeTypeList []ChartDeployNodeType `json:"node_type_list" yaml:"node_type_list"`
}

type ChartDeploy struct {
	Scheme    string            `json:"scheme" yaml:"scheme"`
	Platforms []ChartDeployMain `json:"platforms" yaml:"platforms"`
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
