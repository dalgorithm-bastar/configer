package generation

/****************** 基础设施信息 ******************/

type InfraNetUnit struct {
	Name    string `json:"name" yaml:"name"`
	Adaptor string `json:"adapter" yaml:"adapter"`
	Ipv4    string `json:"ipv4" yaml:"ipv4"`
	Ipv6    string `json:"ipv6" yaml:"ipv6"`
}

type InfraUserUnit struct {
	Name string `json:"name" yaml:"name"`
}

type InfraHostUnit struct {
	NodeId    uint16          `json:"nodeId" yaml:"nodeId"`
	NodeIndex uint16          `json:"nodeIndex" yaml:"nodeIndex"`
	HostName  string          `json:"hostName" yaml:"hostName"`
	Location  string          `json:"location" yaml:"location"`
	Network   []InfraNetUnit  `json:"network" yaml:"network"`
	User      []InfraUserUnit `json:"user" yaml:"user"`
}

type InfraMain struct {
	Host []InfraHostUnit `json:"host" yaml:"host"`
}

/********************************************/

/**************** 部署信息 *******************/

type DeployArtificialUnit struct {
	Url        string `json:"url" yaml:"url"`
	DeployPath string `json:"deployPath" yaml:"deployPath"`
	ConfigPath string `json:"configPath" yaml:"configPath"`
}

type DeployMain struct {
	Artificial DeployArtificialUnit `json:"artifact" yaml:"artifact"`
	Node       []InfraHostUnit      `json:"node" yaml:"node"`
}

/********************************************/

/****************** 服务信息 ******************/

type SrvTpcStatUnit struct {
	TpcName string `json:"topicName" yaml:"topicName"`
	IsRMB   uint16 `json:"isRMB" yaml:"isRMB"`
}

type SrvStatementUnit struct {
	NetName  string           `json:"netName" yaml:"netName"`
	BizTopic []SrvTpcStatUnit `json:"bizTopic" yaml:"bizTopic"`
	SetTopic []SrvTpcStatUnit `json:"setTopic" yaml:"setTopic"`
}

type SrvMain struct {
	InnerTopicNet string             `json:"innerTopicNet" yaml:"innerTopicNet"`
	PubTopic      []SrvStatementUnit `json:"pubTopic" yaml:"pubTopic"`
	SubTopic      []SrvStatementUnit `json:"subTopic" yaml:"subTopic"`
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

/****************** 权限记录文件 ******************/

type PermFile struct {
	FilePerms []PermUnit `yaml:"filePerms"`
}

type PermUnit struct {
	Path  string `yaml:"path"`
	IsDir string `yaml:"isDir"`
	Perm  string `yaml:"perm"`
}

/**********************************************/
