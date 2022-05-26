package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/generation"
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"

	"github.com/configcenter/pkg/pb"
)

// Get 用于从服务端获取信息
func Get(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	//确定请求结构体内容是否合法
	err := checkData(req)
	if err != nil {
		log.Sugar().Infof("checkData failed in GET req: %+v, err: %v", req, err)
		return err, nil, nil
	}
	//根据约定常量确定要获取的资源类型
	switch req.Target {
	case TargetConfig:
		return getConfig(ctx, req)
	case TargetInfrastructure:
		return getInfrastructure(ctx, req)
	case TargetVersion:
		return getVersion(ctx, req)
	case TargetCache:
		return getCache(ctx, req)
	case TargetRaw:
		return getRaw(ctx, req)
	//默认返回错误
	default:
		err := fmt.Sprintf("Target of %s can not be recognized in GET req", req.Target)
		log.Sugar().Infof(err)
		return errors.New(err), nil, nil
	}
}

//判断请求体是否符合约定，不符合则提前返回错误
func checkData(data *pb.CfgReq) error {
	//无用户名时返回错误
	if data.UserName == "" {
		return errors.New("no username when getting")
	}
	//target字段为空时返回错误
	if data.Target == "" {
		return errors.New("no target specified when getting")
	}
	if !CheckFlag(_reqTarget, data.Target) {
		return errors.New(fmt.Sprintf("wrong arg of %s, %s", _reqType, _reqTarget))
	}
	if data.Target == TargetRaw || data.Target == TargetConfig {
		if data.Target == TargetConfig {
			if !manager.regExp.RegExpOfEnvNum.MatchString(data.EnvNum) {
				return errors.New(fmt.Sprintf("wrong envNum:  %s", data.EnvNum))
			}
			if data.ArgRange == nil || len(data.ArgRange.TopicIp) == 0 || len(data.ArgRange.TopicPort) == 0 {
				return errors.New(fmt.Sprintf("wrong ip or port range"))
			}
		}
		if !manager.regExp.RegExpOfVersion.MatchString(data.Version) {
			return errors.New(fmt.Sprintf("wrong version input:  %s", data.Version))
		}
	}
	return nil
}

func getConfig(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	rawData, err := repository.Src.GetbyPrefix(util.Join("/", req.Version, req.Scheme))
	if err != nil {
		log.Sugar().Errorf("get rawData from repository err of %v, under path %s", err, util.Join(req.Version, req.Scheme))
		return err, nil, nil
	}
	if rawData == nil || len(rawData) <= 0 {
		return errors.New(fmt.Sprintf("no rawData on remote yet at path:%s/%s", req.Version, req.Scheme)), nil, nil
	}
	err, _, file := getInfrastructure(ctx, req)
	if err != nil {
		return err, nil, nil
	}
	configData, err := generation.Generate(file.FileData, rawData, req.EnvNum, req.ArgRange.TopicIp, req.ArgRange.TopicPort)
	if err != nil {
		return err, nil, nil
	}
	compressedData, err := util.CompressToStream("config.tar.gz", configData)
	if err != nil {
		return err, nil, nil
	}
	return nil, nil, &pb.AnyFile{
		FileName: "config.tar.gz",
		FileData: compressedData,
	}
}

func getInfrastructure(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	file, err := repository.Src.Get(repository.Infrastructure)
	if err != nil {
		log.Sugar().Errorf("get Infrastructure from repository err of %v, under path %s", err, repository.Infrastructure)
		return err, nil, nil
	}
	if file == nil || len(file) <= 0 {
		return errors.New("no infrastructure on remote yet"), nil, nil
	}
	return nil, nil, &pb.AnyFile{
		FileName: repository.Infrastructure,
		FileData: file,
	}
}

func getVersion(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	fileMap, err := repository.Src.GetbyPrefix(repository.Versions)
	if err != nil {
		log.Sugar().Errorf("get versions from repository err of %v, under path %s", err, repository.Versions)
		return err, nil, nil
	}
	if fileMap == nil || len(fileMap) <= 0 {
		return errors.New("no versions on remote yet"), nil, nil
	}
	if _, ok := fileMap[repository.Versions]; !ok {
		errorIns := fmt.Sprintf("error happened on remote, with unexpected version date: %+v", fileMap)
		log.Sugar().Errorf(errorIns)
		return errors.New(errorIns), nil, nil
	}
	vSli := strings.Split(string(fileMap[repository.Versions]), ",")
	var versions []*pb.VersionInfo
	for i := 0; 3*i < len(vSli); i++ {
		versions = append(versions, &pb.VersionInfo{
			Name: vSli[i],
			User: vSli[i+1],
			Time: vSli[i+2],
		})
	}
	return nil, versions, nil
}

func getCache(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	cache, err := repository.Src.GetbyPrefix(req.UserName)
	if err != nil {
		log.Sugar().Errorf("get rawData from repository err of %v, under path %s", err, util.Join(req.Version, req.Scheme))
		return err, nil, nil
	}
	if cache == nil || len(cache) <= 0 {
		return errors.New(fmt.Sprintf("no cache on remote yet at path:%s", req.UserName)), nil, nil
	}
	compressedData, err := util.CompressToStream("cache.tar.gz", cache)
	if err != nil {
		return err, nil, nil
	}
	return nil, nil, &pb.AnyFile{
		FileName: "cache.tar.gz",
		FileData: compressedData,
	}
}

func getRaw(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	rawData, err := repository.Src.GetbyPrefix(req.Version)
	if err != nil {
		log.Sugar().Errorf("get rawData from repository err of %v, under path %s", err, util.Join(req.Version, req.Scheme))
		return err, nil, nil
	}
	if rawData == nil || len(rawData) <= 0 {
		return errors.New(fmt.Sprintf("no rawData on remote yet at path:%s", req.Version)), nil, nil
	}
	compressedData, err := util.CompressToStream("rawdata.tar.gz", rawData)
	if err != nil {
		return err, nil, nil
	}
	return nil, nil, &pb.AnyFile{
		FileName: "rawdata.tar.gz",
		FileData: compressedData,
	}
}
