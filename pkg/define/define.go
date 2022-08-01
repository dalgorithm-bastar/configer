// Package define 用于存放可能要被任意包引用的常量，避免循环引用。该包不得引入任何项目内的包
package define

//用于标示配置文件的键值
const (
	EtcdEndpoints        = "etcd.endpoints"
	EtcdUserName         = "etcd.username"
	EtcdPassWord         = "etcd.password"
	EtcdOperationTimeout = "etcd.operationtimeout"
	GrpcSocket           = "grpc.socket"
	GrpcLockTimeout      = "grpc.locktimeout"
	LogLogPath           = "log.logpath"
	LogRecordLevel       = "log.recordlevel"
	LogEncodingType      = "log.encodingtype"
	LogFileName          = "log.filename"
	LogMaxAge            = "log.maxage"

	// 与存储路径有关的常量均应放置在该处，统一管理和修改，避免高层之间产生循环引用
	EtcdType           = "etcd"
	CompressedFileType = "compressFile"
	DirType            = "dirType"
	Service            = "service.yaml"
	Template           = "template"
	Infrastructure     = "infrastructure.yaml"
	Deployment         = "deployment.yaml"
	TopicInfo          = "topicInfo.json"
	Manipulations      = "manipulations"
	Versions           = "versions"
	Perms              = "perm.yaml"
	EnvFile            = "envFile.yaml"

	// 用于校验文件包路径
	DeploymentFlag = "deployment"
	ServiceFlag    = "service"
	TemplateFlag   = "template"
)
