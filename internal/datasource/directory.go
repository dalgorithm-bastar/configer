// Package datasource 数据源为本地文件夹时，将storage初始化为该类型
package datasource

import (
	"context"
	"errors"
	"strings"

	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/util"
)

type DirType struct {
	path string
	data map[string][]byte
}

func NewDirType(path string) (*DirType, error) {
	instance := new(DirType)
	instance.setPath(path)
	err := instance.getData()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *DirType) setPath(path string) {
	c.path = path
}

func (c *DirType) getData() error {
	//source, err := util.WalkFromPath(c.path)
	source, _, err := util.LoadDirWithPermFile(c.path, util.Separator, define.Template)
	//source[define.Perms] = pemFile
	if err != nil {
		return err
	}
	c.data = source
	return nil
}

func (c *DirType) Put(key, value string) error {
	return errors.New("cannot Put file in local compressed file")
}

func (c *DirType) Get(key string) ([]byte, error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	//return nil, errors.New(fmt.Sprintf("cannot find %s in local compressedFile", key))
	return nil, nil
}

func (c *DirType) Delete(key string) error {
	return errors.New("cannot delete file in local compressed file")
}

// GetbyPrefix 范围获取
func (c *DirType) GetbyPrefix(prefix string) (map[string][]byte, error) {
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
func (c *DirType) DeletebyPrefix(prefix string) error {
	return errors.New("cannot DeletebyPrefix in local compressed file")
}

func (c *DirType) GetSourceDataorOperator() interface{} {
	source := make(map[string][]byte)
	for k, v := range c.data {
		source[k] = v
	}
	return source
}

func (c *DirType) AtomicCommit(putMap map[string]string, deleteMap []string) error {
	return errors.New("cannot AtomicCommit in local compressed file")
}

func (c *DirType) GracefullyClose(ctx context.Context) {

}
