package datasource

import (
	"errors"
	"strings"
)

type Stream struct {
	c *CompressFileType
}

func NewStream(fileMap map[string][]byte) *Stream {
	stream := &Stream{
		c: &CompressFileType{
			path: "",
			data: make(map[string][]byte),
		},
	}
	for k, v := range fileMap {
		stream.c.data[k] = v
	}
	return stream
}

func (s *Stream) Put(key, value string) error {
	return errors.New("cannot Put file in stream file")
}

func (s *Stream) Get(key string) ([]byte, error) {
	if v, ok := s.c.data[key]; ok {
		return v, nil
	}
	//return nil, errors.New(fmt.Sprintf("cannot find %s in local compressedFile", key))
	return nil, nil
}

func (s *Stream) Delete(key string) error {
	return errors.New("cannot delete file in stream file")
}

// GetbyPrefix 范围获取
func (s *Stream) GetbyPrefix(prefix string) (map[string][]byte, error) {
	//todo 提升效率
	resMap := make(map[string][]byte)
	for k, v := range s.c.data {
		if strings.HasPrefix(k, prefix) {
			resMap[k] = v
		}
	}
	return resMap, nil
}

// DeletebyPrefix 范围删除
func (s *Stream) DeletebyPrefix(prefix string) error {
	return errors.New("cannot DeletebyPrefix in stream file")
}

func (s *Stream) GetSourceDataorOperator() interface{} {
	source := make(map[string][]byte)
	for k, v := range s.c.data {
		source[k] = v
	}
	return source
}

func (s *Stream) AcidCommit(putMap map[string]string, deleteMap []string) error {
	return errors.New("cannot AcidCommit in stream file")
}
