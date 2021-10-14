// Package manage 启动时构建Manage类并注册至grpc服务端
package manage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

type GrpcInfoStruct struct {
	Port   string `json:"port"`
	Socket string `json:"socket"`
}

// NewManager 读取grpc配置并初始化manager实例
func NewManager(grpcConfigLocation string) error {
	manager = new(Manager)
	file, err := os.Open(grpcConfigLocation)
	if err != nil {
		return err
	}
	binaryFlie, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(binaryFlie, &manager.grpcInfo)
	if err != nil {
		return err
	}
	versionFormat, err := regexp.Compile(VersionString)
	if err != nil {
		return err
	}
	templateFormat, err := regexp.Compile(TemplateString)
	if err != nil {
		return err
	}
	manager.regExp = regExpStruct{
		RegExpOfVersion:  versionFormat,
		RegExpOfTemplate: templateFormat,
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func GetManager() *Manager {
	return manager
}

func GetGrpcInfo() *GrpcInfoStruct {
	return &manager.grpcInfo
}
