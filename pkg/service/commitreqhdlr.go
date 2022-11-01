package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/define"
	"github.com/configcenter/pkg/generation"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"go.uber.org/zap"
)

//post 提交函数post认为将要提交的对象是正确的，并在此基础上进一步执行提交操作。对象的正确性在put操作中保证。
func commit(ctxRoot context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	if req.UserName == "" || strings.Contains(req.UserName, ",") || strings.Contains(req.UserName, ".") {
		return errors.New("empty username of username contains point or comma, commit request deny"), nil, nil
	}
	if req.Target == TargetInfrastructure {
		if req.File == nil || req.File.FileData == nil || len(req.File.FileData) == 0 {
			return errors.New("infrastructure.yaml with err content,please checkout"), nil, nil
		}
		err := repository.Src.Put(define.Infrastructure, string(req.File.FileData))
		if err != nil {
			return err, nil, nil
		}
		return nil, nil, nil
	}
	versionStr, newVersion, err := getNextVersion(req.Version)
	if err != nil {
		return err, nil, nil
	}
	//获取已缓存的文件
	rawData, err := repository.Src.GetbyPrefix(req.UserName)
	if err != nil {
		log.Logger.Error("get rawData from repository err", zap.Any("err", err), zap.String("path(user)", req.UserName))
		return err, nil, nil
	}
	if rawData == nil || len(rawData) <= 0 {
		return errors.New(fmt.Sprintf("no rawData on remote yet at path:%s", req.UserName)), nil, nil
	}
	err, _, file := getInfrastructure(ctxRoot, req)
	if err != nil {
		return err, nil, nil
	}
	_, err = generation.Generate(file.FileData, rawData, false, "01", "", []string{"0.0.0.0", "255.255.255.255"},
		[]string{"1024", "65535"}, []string{"1024", "65535"}, []string{})
	if err != nil {
		return errors.New(fmt.Sprintf("generate cfgFile err:%s,please checkout cache under username:%s", err, req.UserName)), nil, nil
	}
	fileMap, _ := checkAndReplaceRootPath(rawData, newVersion)
	//加锁并上传
	cliRes := repository.Src.GetSourceDataorOperator()
	if _, ok := cliRes.(*clientv3.Client); !ok {
		//获取etcd客户端失败，返回获取的类型
		msg := fmt.Sprintf("get etcd clientv3 err, type of datasource: %v", reflect.TypeOf(cliRes))
		log.Logger.Error(msg)
		return errors.New(msg), nil, nil
	}
	client, _ := cliRes.(*clientv3.Client)
	//给etcd加锁
	//设置超时时间
	ctxTimeout, cancelFunc := context.WithTimeout(ctxRoot, time.Duration(manager.grpcInfo.LockTimeout)*time.Second)
	defer cancelFunc()
	response, err := client.Grant(ctxTimeout, int64(manager.grpcInfo.LockTimeout))
	if err != nil {
		log.Logger.Error("get lease from etcd err", zap.Any("err", err))
		return err, nil, nil
	}
	session, err := concurrency.NewSession(client, concurrency.WithLease(response.ID))
	if err != nil {
		log.Logger.Error("get session from etcd err", zap.Any("err", err))
		return err, nil, nil
	}
	defer func() {
		if session != nil {
			err := session.Close()
			if err != nil {
				log.Logger.Warn("Close session err", zap.Any("err", err))
			}
		}
	}()
	mutex := concurrency.NewMutex(session, _lockName)
	err = mutex.Lock(ctxTimeout)
	session.Orphan()
	if err != nil {
		log.Logger.Info("lock etcd err", zap.Any("err", err))
		return err, nil, nil
	}
	//加锁后进行操作
	timeStamp, newVersionStr := strconv.FormatInt(time.Now().UnixNano()/1e6, 10), ""
	if versionStr != "" {
		newVersionStr = util.Join(",", versionStr, newVersion, req.UserName, timeStamp)
	} else {
		newVersionStr = util.Join(",", newVersion, req.UserName, timeStamp)
	}
	fileMap[define.Versions] = newVersionStr
	var deleteKeySlice []string
	for path, _ := range rawData {
		deleteKeySlice = append(deleteKeySlice, path)
	}
	err = repository.Src.AtomicCommit(fileMap, deleteKeySlice)
	if err != nil {
		log.Logger.Error("AtomicCommit err when posting", zap.Any("err", err), zap.Any("filemap", fileMap), zap.Any("deleteKeySlice", deleteKeySlice))
		return err, nil, nil
	}
	defer func() {
		if mutex != nil {
			err := mutex.Unlock(ctxTimeout)
			if err != nil {
				log.Logger.Warn("Close session err", zap.Any("err", err))
			}
		}
	}()
	return nil, []*pb.VersionInfo{{
		Name: newVersion,
		User: req.UserName,
		Time: timeStamp,
	}}, nil
}

func getNextVersion(version string) (string, string, error) {
	newVersion := ""
	//指定格式不满足要求
	if version != "" && !manager.regExp.RegExpOfVersion.MatchString(version) {
		return "", "", errors.New(fmt.Sprintf("err version input:%s", version))
	}
	//获取历史版本号
	versionStr, err := repository.Src.Get(define.Versions)
	if err != nil {
		log.Logger.Error("get former versions from repository err", zap.Any("err", err), zap.String("version(path)", define.Versions))
		return "", "", err
	}
	if versionStr == nil || len(versionStr) <= 0 {
		//未获取到值，为第一次提交
		if version != "" {
			newVersion = version
		} else {
			newVersion = "0.0.1"
		}
	} else {
		//非第一次提交，查重或递增
		versionSli := strings.Split(string(versionStr), ",")
		if version != "" {
			for i := 0; 3*i < len(versionSli); i++ {
				if version == versionSli[3*i] {
					return "", "", errors.New(fmt.Sprintf("version input repeated:%s", version))
				}
			}
			newVersion = version
		} else {
			oldVersion := versionSli[len(versionSli)-3]
			oldSli := strings.SplitN(oldVersion, ".", 3)
			lastNum, _ := strconv.Atoi(oldSli[2])
			oldSli[2] = strconv.Itoa(lastNum + 1)
			newVersion = strings.Join(oldSli, ".")
		}
	}
	return string(versionStr), newVersion, nil
}

/*func commit(ctxRoot context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
    if req.UserName == "" {
        return errors.New("empty username, commit request deny"), nil, "", nil
    }
    //若当前用户不存在缓存数据，则快速返回处理结果
    fileMap, err := getDataByUserName(req.UserName)
    if err != nil {
        log.Sugar().Errorf("Get file from etcd err when posting:%v, UserName is %s", err, req.UserName)
        return err, nil, "", nil
    }
    if fileMap == nil || len(fileMap) == 0 {
        log.Sugar().Infof("Get nil file from etcd when posting, UserName is %s", req.UserName)
        return errors.New(fmt.Sprintf("no data under username %s, commit request deny", req.UserName)), nil, "", nil
    }
    res := repository.Src.GetSourceDataorOperator()
    if client, ok := res.(*clientv3.Client); ok {
        //给etcd加锁
        //设置超时时间
        ctxTimeout, cancelFunc := context.WithTimeout(ctxRoot, time.Duration(manager.grpcInfo.LockTimeout)*time.Second)
        defer cancelFunc()
        response, err := client.Grant(ctxTimeout, int64(manager.grpcInfo.LockTimeout))
        if err != nil {
            log.Sugar().Infof("get lease from etcd err:%v", err)
            return err, nil, "", nil
        }
        session, err := concurrency.NewSession(client, concurrency.WithLease(response.ID))
        if err != nil {
            log.Sugar().Infof("get session from etcd err:%v", err)
            return err, nil, "", nil
        }
        defer func() {
            if session != nil {
                err := session.Close()
                if err != nil {
                    log.Sugar().Warnf("Close session err:%v", err)
                }
            }
        }()
        mutex := concurrency.NewMutex(session, LockName)
        err = mutex.Lock(ctxTimeout)
        session.Orphan()
        if err != nil {
            log.Sugar().Infof("lock etcd err:%v", err)
            return err, nil, "", nil
        }
        //加锁后进行操作
        if req.Target == nil || req.Target[0] == "" {
            versionSlice, version, err := getNextVersion()
            if err != nil {
                erru := mutex.Unlock(ctxTimeout)
                return errors.New(err.Error() + erru.Error()), nil, "", nil
            }
            resErr, resStringSlice, resString, resByte := createVersion(fileMap, versionSlice, version)
            erru := mutex.Unlock(ctxTimeout)
            if erru != nil {
                return erru, nil, "", nil
            }
            return resErr, resStringSlice, resString, resByte
        } else {
            //若指定版本号上传，校验参数是否合法
            versionSlice, err := checkVersionFormat(req.Target[0])
            if err != nil {
                erru := mutex.Unlock(ctxTimeout)
                return errors.New(err.Error() + erru.Error()), nil, "", nil
            }
            version := req.Target[0]
            resErr, resStringSlice, resString, resByte := createVersion(fileMap, versionSlice, version)
            erru := mutex.Unlock(ctxTimeout)
            if erru != nil {
                return erru, nil, "", nil
            }
            return resErr, resStringSlice, resString, resByte
        }
    }
    //获取etcd客户端失败，返回获取的类型
    msg := fmt.Sprintf("get etcd clientv3 err, type of datasource: %v", reflect.TypeOf(res))
    log.Sugar().Infof(msg)
    return errors.New(msg), nil, "", nil
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
        log.Sugar().Warnf("get nil version slice when posting with assigned versionNum %s, if this is not the first post, there may be an err", version)
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
}*/
