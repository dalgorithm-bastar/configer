//虚拟的流式文件对象
package util

import (
	"os"
	"time"
)

type StreamFile struct {
	data    []byte
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func NewStreamFile(inputData []byte, name string, size int64) *StreamFile {
	s := new(StreamFile)
	s.data = inputData
	s.name = name
	s.size = size
	s.mode = os.FileMode(777)
	s.modTime = time.Now()
	return s
}

func (s *StreamFile) Name() string {
	return s.name
}

func (s *StreamFile) Size() int64 {
	return s.size
}

func (s *StreamFile) Mode() os.FileMode {
	return s.mode
}
func (s *StreamFile) ModTime() time.Time {
	return s.modTime
}
func (s *StreamFile) IsDir() bool {
	if len(s.name) >= 1 && s.name[len(s.name)-1] == '/' {
		return true
	}
	return false
}
func (s *StreamFile) Sys() interface{} {
	return nil
}
