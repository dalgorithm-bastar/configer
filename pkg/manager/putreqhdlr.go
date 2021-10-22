package manage

import (
	"context"
	"errors"
	"fmt"
	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/template"
	"github.com/configcenter/pkg/util"
	"strings"
)

//put接口用于导入配置并缓存，缓存接口单独设置以便后续升级
func put(ctx context.Context, req *pb.CfgReq) (error, []string, string, []byte) {
	//检查请求是否合法
	if req.UserName == "" || req.File == nil || req.File.FileData == nil {
		return errors.New("bad request, missing username or file or filedata"), nil, "", nil
	}
	//获取文件存储路径和内容
	fileMap, err := util.DecompressFromStream(req.File.FileName, req.File.FileData)
	if err != nil {
		log.Sugar().Infof("Decompressed from file err:%v, filename:%s", err, req.File.FileName)
		return err, nil, "", nil
	}
	if fileMap == nil || len(fileMap) == 0 {
		log.Sugar().Infof("Decompressed from file and get nil filedata, filename:%s", req.File.FileName)
		return errors.New("get no data from filedata, please checkout file uploaded"), nil, "", nil
	}
	//返回根目录(版本号)，不同环境号下对应的模板路径，不同环境号对应的集群名称，以及错误类型
	rootPath, envToTmplMap, envToClusterMap, err := checkFilePath(fileMap)
	if err != nil {
		return err, nil, "", nil
	}
	//新建临时存储对象
	srcForCheck := repository.NewStream(fileMap)
	//根据不同环境号新建模板对象
	for envNum, templateFilePaths := range envToTmplMap {
		for _, templateFilePath := range templateFilePaths {
			templateIns, err := template.NewTemplateImpl(srcForCheck, "0", "0", "templateIns", rootPath, envNum)
			if err != nil {
				log.Sugar().Errorf("init tmpl err when putting, err:%v", err)
				return err, nil, "", nil
			}
			content, err := srcForCheck.Get(templateFilePath)
			if err != nil {
				log.Sugar().Infof("get tmpl err when putting, err:%v, tmplPath:%s", err, templateFilePath)
				return err, nil, "", nil
			}
			_, err = templateIns.Fill(content, templateFilePath)
			if err != nil {
				log.Sugar().Infof("fill tmpl err when putting, err:%v, tmplPath:%s, tmplContent:%s", err, templateFilePath, content)
				return err, nil, "", nil
			}
		}
	}
	//校验完成，开始操作etcd中的数据

	//获取数据源，然后以事务提交方式进行缓存
	putmap := make(map[string]string)
	if rootPath == req.UserName {
		for filePath, content := range fileMap {
			putmap[filePath] = string(content)
		}
	} else {
		//替换根目录名
		for filePath, content := range fileMap {
			filePathSlice := strings.SplitN(filePath, "/", 2)
			filePathSlice[0] = req.UserName
			newFilePath := strings.Join(filePathSlice, "/")
			putmap[newFilePath] = string(content)
		}
	}
	//存储环境号和对应的集群名称
	var envSlice, clusterSlice []string
	for envNum, clusterMap := range envToClusterMap {
		envSlice = append(envSlice, envNum)
		for clusterName, _ := range clusterMap {
			clusterSlice = append(clusterSlice, clusterName)
		}
		//csv
		putmap[util.Join("/", req.UserName, envNum, repository.Clusters)] = strings.Join(clusterSlice, ",")
		clusterSlice = clusterSlice[0:0]
	}
	putmap[req.UserName+"/"+repository.Envs] = strings.Join(envSlice, ",")
	//删除用户上次缓存的文件，避免污染
	userFileMap, err := repository.Src.GetbyPrefix(req.UserName)
	if err != nil {
		log.Sugar().Infof("delete files under username %s err:%v", req.UserName, err)
		return err, nil, "", nil
	}
	var deleteSlice []string
	if userFileMap != nil && len(userFileMap) != 0 {
		for k, _ := range userFileMap {
			deleteSlice = append(deleteSlice, k)
		}
	}
	err = repository.Src.AcidCommit(putmap, deleteSlice)
	if err != nil {
		log.Sugar().Errorf("AcidCommit err of %v when putting, putmap:%+v", err, putmap)
		return err, nil, "", nil
	}
	log.Sugar().Info("commit success")
	return nil, nil, "", nil
}

// 校验压缩包内文件路径是否满足既定规则，并返回根目录名，每个环境下模板文件的相对路径，每个环境下的集群名哈希表
func checkFilePath(fileMap map[string][]byte) (string, map[string][]string, map[string]map[string]string, error) {
	envToTmplsMap := make(map[string][]string)
	envToClusterMap := make(map[string]map[string]string)
	rootPathMap := make(map[string]string)
	for k, _ := range fileMap {
		//切分到能够分辨集群名称
		pathSlice := strings.SplitN(k, "/", 4)
		//根目录下不允许放文件
		if len(pathSlice) <= 1 {
			log.Sugar().Infof("illegal file path of %s when putting, file is put under rootpath", k)
			return "", nil, nil, errors.New(fmt.Sprintf("can not put file under rootpath, filename: %s", k))
		}
		//检测是否存在多个根文件夹
		if _, ok := rootPathMap[pathSlice[0]]; !ok {
			rootPathMap[pathSlice[0]] = ""
			if len(rootPathMap) > 1 {
				log.Sugar().Infof("multi rootpath found when putting, paths:%+v", rootPathMap)
				return "", nil, nil, errors.New(fmt.Sprintf("multi rootPath found in compressed file: %v", rootPathMap))
			}
		}
		//长度不满足要求，或属于文件夹对象时，直接跳过
		if len(pathSlice) < 4 || k[len(k)-1] == '/' {
			continue
		}
		//筛选模板文件，若templates关键字在路径中重复，返回错误
		subStringIndexs := manager.regExp.RegExpOfTemplate.FindAllStringIndex(k, -1)
		if subStringIndexs == nil || len(subStringIndexs) == 0 {
			continue
		}
		if len(subStringIndexs) > 1 {
			existFlag := false
			for _, index := range subStringIndexs {
				if index[0] == 0 {
					err := fmt.Sprintf("key word \"template\" found as rootPath %s", k)
					log.Sugar().Info(err)
					return "", nil, nil, errors.New(err)
				}
				if k[index[0]-1] == '/' {
					if existFlag {
						err := fmt.Sprintf("key word \"template\" repeated in Path %s", k)
						log.Sugar().Info(err)
						return "", nil, nil, errors.New(err)
					}
					existFlag = true
				}
			}
		}
		//按照环境号建立需要填充的模板的哈希表
		envToTmplsMap[pathSlice[1]] = append(envToTmplsMap[pathSlice[1]], k)
		//按照环境号整理对应的集群名称
		if clusterMap, envOk := envToClusterMap[pathSlice[1]]; envOk {
			if _, clusterOk := clusterMap[pathSlice[2]]; !clusterOk {
				clusterMap[pathSlice[2]] = ""
			}
		} else {
			envToClusterMap[pathSlice[1]] = make(map[string]string)
			envToClusterMap[pathSlice[1]][pathSlice[2]] = ""
		}
	}
	var rootPath string
	for k, _ := range rootPathMap {
		rootPath = k
	}
	return rootPath, envToTmplsMap, envToClusterMap, nil
}