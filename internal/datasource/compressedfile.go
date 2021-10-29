//数据源为本地压缩包时，将storage初始化为该类型
package datasource

import (
    "context"
    "errors"
    "strings"

    "github.com/configcenter/pkg/util"
)

type CompressFileType struct {
    path string
    data map[string][]byte
}

func NewCompressedFileType(path string) (*CompressFileType, error) {
    instance := new(CompressFileType)
    instance.setPath(path)
    err := instance.getData()
    if err != nil {
        return nil, err
    }
    return instance, nil
}

func (c *CompressFileType) setPath(path string) {
    c.path = path
}

func (c *CompressFileType) getData() error {
    source, err := util.DecompressFromPath(c.path)
    if err != nil {
        return err
    }
    c.data = source
    return nil
}

func (c *CompressFileType) Put(key, value string) error {
    return errors.New("cannot Put file in local compressed file")
}

func (c *CompressFileType) Get(key string) ([]byte, error) {
    if v, ok := c.data[key]; ok {
        return v, nil
    }
    //return nil, errors.New(fmt.Sprintf("cannot find %s in local compressedFile", key))
    return nil, nil
}

func (c *CompressFileType) Delete(key string) error {
    return errors.New("cannot delete file in local compressed file")
}

// GetbyPrefix 范围获取
func (c *CompressFileType) GetbyPrefix(prefix string) (map[string][]byte, error) {
    //todo 提升效率
    resMap := make(map[string][]byte)
    for k, v := range c.data {
        if strings.HasPrefix(k, prefix) {
            resMap[k] = v
        }
    }
    return resMap, nil
}

// DeletebyPrefix 范围删除
func (c *CompressFileType) DeletebyPrefix(prefix string) error {
    return errors.New("cannot DeletebyPrefix in local compressed file")
}

func (c *CompressFileType) GetSourceDataorOperator() interface{} {
    source := make(map[string][]byte)
    for k, v := range c.data {
        source[k] = v
    }
    return source
}

func (c *CompressFileType) AcidCommit(putMap map[string]string, deleteMap []string) error {
    return errors.New("cannot AcidCommit in local compressed file")
}

func (c *CompressFileType) GracefullyClose(ctx context.Context) {

}
