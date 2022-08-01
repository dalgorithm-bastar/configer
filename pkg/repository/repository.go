// Package repository 数据层接口，隔离领域层和具体的数据操作
package repository

import (
	"context"

	"github.com/configcenter/internal/datasource"
	"github.com/configcenter/pkg/define"
	"xchg.ai/sse/gracefully"
)

const ()

var Src Storage

type Storage interface {
	Put(key, value string) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	GetbyPrefix(prefix string) (map[string][]byte, error)
	DeletebyPrefix(prefix string) error
	GetSourceDataorOperator() interface{}
	AtomicCommit(map[string]string, []string) error
	GracefullyClose(ctx context.Context)
}

// NewStorage 初始化底层数据结构。对于etcd，config参数为配置文件路径，对于本地压缩包，config参数为压缩包绝对路径
func NewStorage(ctx context.Context, dataSourceType, config string) error {
	switch dataSourceType {
	case define.EtcdType:
		e, err := datasource.NewEtcdType()
		if err != nil {
			return err
		}
		Src = e
	case define.CompressedFileType:
		c, err := datasource.NewCompressedFileType(config)
		if err != nil {
			return err
		}
		Src = c
	case define.DirType:
		d, err := datasource.NewDirType(config)
		if err != nil {
			return err
		}
		Src = d
	}
	gracefully.Go(func() { Src.GracefullyClose(ctx) })
	return nil
}

/*// NewStream 返回一个仅用于检测合法性的临时数据接口
func NewStream(fileMap map[string][]byte) Storage {
    s := datasource.NewStream(fileMap)
    var srcForCheck Storage = s
    return srcForCheck
}*/
