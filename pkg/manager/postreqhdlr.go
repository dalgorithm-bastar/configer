package manage

import (
	"context"
	"errors"
	"fmt"
	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"
	"strconv"
	"strings"
)

//post 提交函数post认为将要提交的对象是正确的，并在此基础上进一步执行提交操作。对象的正确性在put操作中保证。
func post(ctx context.Context, req *pb.CfgReq) (error, []string, string, []byte) {
	if req.UserName == "" {
		return errors.New("empty username, commit request deny"), nil, "", nil
	}
	//若当前用户不存在缓存数据，则快速返回处理结果
	fileMap, err := getDataByUserName(req.UserName)
	if err != nil {
		log.Sugar().Errorf("Get file from etcd err of %v when posting, UserName is %s", err, req.UserName)
		return err, nil, "", nil
	}
	if fileMap == nil || len(fileMap) == 0 {
		log.Sugar().Infof("Get nil file from etcd when posting, UserName is %s", req.UserName)
		return errors.New(fmt.Sprintf("no data under username %s, commit request deny", req.UserName)), nil, "", nil
	}
	if req.Target == nil || req.Target[0] == "" {
		versionSlice, version, err := getNextVersion()
		if err != nil {
			return err, nil, "", nil
		}
		return createVersion(fileMap, versionSlice, version)
	} else {
		//若指定版本号上传，校验参数是否合法
		versionSlice, err := checkVersionFormat(req.Target[0])
		if err != nil {
			return err, nil, "", nil
		}
		version := req.Target[0]
		return createVersion(fileMap, versionSlice, version)
	}
}

func getDataByUserName(userName string) (map[string][]byte, error) {
	//在上级函数中打印过错误日志，此处无需重复打印
	fileMap, err := repository.Src.GetbyPrefix(util.GetPrefix(userName))
	if err != nil {
		return nil, err
	}
	if fileMap == nil {
		return nil, nil
	}
	return fileMap, nil
}

func checkVersionFormat(version string) ([]string, error) {
	if !manager.regExp.RegExpOfVersion.MatchString(version) {
		log.Sugar().Infof("illegal version format of %s found when posting", version)
		return nil, errors.New(fmt.Sprintf("illegal version format of %s, commit request deny", version))
	}
	resVersion, err := repository.Src.Get(repository.Versions)
	if err != nil {
		log.Sugar().Errorf("Get versionSlice from etcd err of %v when posting, key is %s", err, repository.Versions)
		return nil, err
	}
	if resVersion == nil {
		log.Sugar().Warnf("get nil version slice when posting with assigned versionNum %s, if this is not the first post, there may be err", version)
		return nil, nil
	}
	versionSlice := strings.Split(string(resVersion), ",")
	for _, v := range versionSlice {
		if version == v {
			log.Sugar().Infof("Get repeated versionNum when posting of %s", version)
			return nil, errors.New(fmt.Sprintf("repaet version num of %s, commit request deny", version))
		}
	}
	return versionSlice, nil
}

func getNextVersion() ([]string, string, error) {
	resVersion, err := repository.Src.Get(repository.Versions)
	if err != nil {
		log.Sugar().Errorf("Get versionSlice from etcd err of %v when posting, key is %s", err, repository.Versions)
		return nil, "", err
	}
	if resVersion == nil {
		return nil, "0.0.1", nil
	}
	versionSlice := strings.Split(string(resVersion), ",")
	formerVersion := versionSlice[len(versionSlice)-1]
	formerVersionSlice := strings.Split(formerVersion, ".")
	newNum, err := strconv.Atoi(formerVersionSlice[len(formerVersionSlice)-1])
	if err != nil {
		log.Sugar().Errorf("trasfer versionNum to int err of %v when adding version num, versionString is %s", err, formerVersion)
		return nil, "", err
	}
	newNum += 1
	formerVersionSlice[len(formerVersionSlice)-1] = strconv.Itoa(newNum)
	return versionSlice, strings.Join(formerVersionSlice, "."), nil
}

func createVersion(fileMap map[string][]byte, formerVersions []string, version string) (error, []string, string, []byte) {
	//提交版本，删除缓存，更新版本号序列
	putmap := make(map[string]string)
	deleteKeySlice := make([]string, 0)
	for k, v := range fileMap {
		PathSlice := strings.SplitN(k, "/", 2)
		if len(PathSlice) <= 1 {
			log.Sugar().Errorf("illegal file path of %s when posting", k)
			return errors.New(fmt.Sprintf("error path: %s", k)), nil, "", nil
		}
		PathSlice[0] = version
		newPath := strings.Join(PathSlice, "/")
		//添加新版本数据，同时删除用户缓存数据
		putmap[newPath] = string(v)
		deleteKeySlice = append(deleteKeySlice, k)
	}
	formerVersions = append(formerVersions, version)
	newVersions := strings.Join(formerVersions, ",")
	putmap[repository.Versions] = newVersions
	err := repository.Src.AcidCommit(putmap, deleteKeySlice)
	if err != nil {
		log.Sugar().Errorf("AcidCommit err of %v when posting, putmap:%+v, deleteSlice:%+v", err, putmap, deleteKeySlice)
		return err, nil, "", nil
	}
	return nil, []string{version}, "", nil
}
