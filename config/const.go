package config

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
    LogMaxAge            = "log.maxage" //设置标示但不从命令行输入
)
