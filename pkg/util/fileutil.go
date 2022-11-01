package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/define"
	"github.com/mholt/archiver"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const Separator = string(os.PathSeparator)

/****************** 权限记录文件 ******************/

type PermFile struct {
	FilePerms []PermUnit `yaml:"filePerms"`
}

type PermUnit struct {
	Path  string `yaml:"path"`
	IsDir string `yaml:"isDir"`
	Perm  string `yaml:"perm"`
}

/**********************************************/

// CompressToStream 将若干个文件按照给定的格式压缩，输入map的key应为文件在压缩包内的相对路径
func CompressToStream(stringWithFormat string, fileMap map[string][]byte) ([]byte, error) {
	if fileMap == nil {
		return nil, errors.New("can not compress nil inputmap")
	}
	out := bytes.NewBuffer([]byte{})
	arc, _ := archiver.ByExtension(stringWithFormat)
	arcW, _ := arc.(archiver.Writer)
	err := arcW.Create(out)
	if err != nil {
		log.Logger.Info("create archiever writer err", zap.Any("err", err), zap.String("stringwithformat", stringWithFormat))
		return nil, err
	}
	//defer arcW.Close()
	for key, v := range fileMap {
		if key[len(key)-1] == '/' {
			continue
		}
		buffer := bytes.NewBuffer(v)
		streamFile := NewStreamFile(v, filepath.Base(key), int64(len(v)))
		simufileInfo := os.FileInfo(streamFile)
		simuFile := archiver.File{
			FileInfo: archiver.FileInfo{
				FileInfo:   simufileInfo,
				CustomName: key,
			},
			ReadCloser: ioutil.NopCloser(buffer),
		}
		err = arcW.Write(simuFile)
		if err != nil {
			log.Logger.Error("write archiever writer err of %v, input format %s", zap.Any("err", err), zap.String("stringwithformat", stringWithFormat))
			return nil, err
		}
	}
	err = arcW.Close()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// DecompressFromStream 根据指定格式解内存压缩包并返回相对路径和文件内容，其中stringWithFormat只需包含合法扩展名即可
func DecompressFromStream(stringWithFormat string, binaryFile []byte) (map[string][]byte, error) {
	if binaryFile == nil {
		return nil, nil
	}
	length := len(binaryFile)
	b := bytes.Reader{}
	b.Reset(binaryFile)

	arc, err := archiver.ByExtension(stringWithFormat)
	if err != nil {
		return nil, err
	}
	arcImpl, ok := arc.(archiver.Reader)
	if !ok {
		log.Logger.Error("archiever reader err", zap.Any("err", err), zap.String("stringwithformat", stringWithFormat))
		return nil, errors.New("input compressed file format error")
	}
	err = arcImpl.Open(&b, int64(length))
	if err != nil {
		return nil, err
	}
	resMap := make(map[string][]byte)
	var key string
	for {
		f, err := arcImpl.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(f.Name(), ".") {
			switch f.Header.(type) {
			case *tar.Header:
				key = f.Header.(*tar.Header).Name
			case zip.FileHeader:
				key = f.Header.(zip.FileHeader).Name
			}
			bitSource, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			resMap[filepath.Clean(key)] = bitSource
		}
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	err = arcImpl.Close()
	if err != nil {
		return nil, err
	}
	return resMap, nil
}

// DecompressFromPath 从本地压缩包导入数据，输入参数为绝对路径，返回的map中，key为所有文件的相对路径，值为内容
func DecompressFromPath(inputPath string) (map[string][]byte, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Logger.Info("open file err of %v when Decompress, filepath %s", zap.Any("err", err), zap.String("inputpath", inputPath))
		return nil, err
	}
	binaryFlie, err := ioutil.ReadAll(file)
	if err != nil {
		log.Logger.Info("Read file err of %v when Decompress, filepath %s", zap.Any("err", err), zap.String("inputpath", inputPath))
		return nil, err
	}
	err = file.Close()
	if err != nil {
		log.Logger.Info("close file err of %v when Decompress, filepath %s", zap.Any("err", err), zap.String("inputpath", inputPath))
		return nil, err
	}
	return DecompressFromStream(inputPath, binaryFlie)
}

/*func WalkFromPath(inputPath string) (map[string][]byte, error) {
    inputPath = filepath.Clean(inputPath)
    fileMap := make(map[string][]byte)
    //读剩余的文件
    err := walkAndRead(fileMap, inputPath, filepath.Base(inputPath))
    return fileMap, err
}

func walkAndRead(fileMap map[string][]byte, inputPath, key string) error {
    fileRange, err := ioutil.ReadDir(inputPath)
    if err != nil {
        return err
    }
    var data []byte
    for _, file := range fileRange {
        if file.IsDir() {
            err = walkAndRead(fileMap, inputPath+Separator+file.Name(), key+"/"+file.Name())
        } else {
            data, err = ioutil.ReadFile(inputPath + Separator + file.Name())
            if err != nil {
                return err
            }
            fileMap[key+"/"+file.Name()] = data
        }
        if err != nil {
            return err
        }
    }
    return nil
}*/

// LoadDirWithPermFile 从文件夹中读取文件，整理成符合规范的输入文件map,separator表示分隔符,keyForTmpl表示存放模板的文件夹名称
func LoadDirWithPermFile(dirPath, separator, keyForTmpl string) (map[string][]byte, []byte, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil || !fileInfo.IsDir() {
		return nil, nil, fmt.Errorf("dir err or not exist, err:%v", err)
	}
	dirPath = filepath.Clean(dirPath)
	pathSli := strings.Split(dirPath, separator)
	var permData []byte
	permStruct := PermFile{}
	fileMap := make(map[string][]byte)
	//用于记录version在路径中所在位置，便于遍历时截断
	lenVersion := len(pathSli[len(pathSli)-1])
	idx := len(dirPath) - lenVersion + 1
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		}
		//整理权限文件
		if strings.Contains(path, separator+keyForTmpl+separator) {
			pathKey := path[idx-1:]
			fileInfo, statErr := os.Stat(path)
			if statErr != nil {
				return statErr
			}
			permStr, isDir := strconv.FormatUint(uint64(fileInfo.Mode()), 8), "0"
			if fileInfo.IsDir() {
				permStr = permStr[len(permStr)-3:]
				isDir = "1"
			}
			//fmt.Println(fileInfo.Mode())
			permStruct.FilePerms = append(permStruct.FilePerms, PermUnit{
				Path:  filepath.ToSlash(pathKey),
				IsDir: isDir,
				Perm:  permStr,
			})
		}
		//整理文件map
		if !info.IsDir() {
			pathKey := path[idx-1:]
			data, readErr := ioutil.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			//转化为标准路径分隔符
			fileMap[filepath.ToSlash(pathKey)] = data
		}
		return err
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error when getting permission value walking the path %q: %v\n", dirPath, err)
	}
	if len(permStruct.FilePerms) != 0 {
		permData, _ = yaml.Marshal(permStruct)
		fileMap[filepath.Base(dirPath)+"/"+define.Perms] = permData
	}
	return fileMap, permData, err
}

// CheckYaml 检查输入的两份yaml文件是否一致，需保证输入的yaml文件格式正确，通常需要跟随在文件解码后使用。该函数用于检查用户的文件编写是否符合预期。源文件必须放在第一个参数。
func CheckYaml(srcYml, yml2 []byte) bool {
	v1, v2 := viper.New(), viper.New()
	v1.SetConfigType("yaml")
	v2.SetConfigType("yaml")
	v1.ReadConfig(bytes.NewBuffer(srcYml))
	v2.ReadConfig(bytes.NewBuffer(yml2))
	m1, m2 := v1.AllSettings(), v2.AllSettings()
	//print(m1, m2)
	if !reflect.DeepEqual(m1, m2) {
		return false
	}
	return true
}
