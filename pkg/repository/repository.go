// Package repository 数据层接口，隔离领域层和具体的数据操作
package repository

import (
    "context"

    "github.com/configcenter/internal/datasource"
    "xchg.ai/sse/gracefully"
)

//与存储路径有关的常量均应放置在该处，统一管理和修改，避免高层之间产生循环引用
const (
    EtcdType           = "etcd"
    CompressedFileType = "compressFile"
    ServiceList        = "servicelist.json"
    Clusters           = "clusters"
    Templates          = "templates"
    Infrastructure     = "infrastructure.json"
    Manipulations      = "manipulations"
    Versions           = "versions"
    Confs              = "configs"
)

var Src Storage

type Storage interface {
    Put(key, value string) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    GetbyPrefix(prefix string) (map[string][]byte, error)
    DeletebyPrefix(prefix string) error
    GetSourceDataorOperator() interface{}
    AcidCommit(map[string]string, []string) error
    GracefullyClose(ctx context.Context)
}

// NewStorage 初始化底层数据结构。对于etcd，config参数为配置文件路径，对于本地压缩包，config参数为压缩包绝对路径
func NewStorage(ctx context.Context, dataSourceType, config string) error {
    switch dataSourceType {
    case EtcdType:
        e, err := datasource.NewEtcdType(config)
        if err != nil {
            return err
        }
        Src = e
    case CompressedFileType:
        c, err := datasource.NewCompressedFileType(config)
        if err != nil {
            return err
        }
        Src = c
    }
    gracefully.Go(func() { Src.GracefullyClose(ctx) })
    return nil
}

// NewStream 返回一个仅用于检测合法性的临时数据接口
func NewStream(fileMap map[string][]byte) Storage {
    s := datasource.NewStream(fileMap)
    var srcForCheck Storage = s
    return srcForCheck
}
