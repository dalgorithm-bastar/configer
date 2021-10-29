package util

import (
    "archive/tar"
    "archive/zip"
    "bytes"
    "errors"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"

    "github.com/configcenter/internal/log"
    "github.com/mholt/archiver"
)

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
        log.Sugar().Infof("create archiever writer err of %v, input format %s", err, stringWithFormat)
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
            log.Sugar().Infof("write archiever writer err of %v, input format %s", err, stringWithFormat)
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
        log.Sugar().Infof("archiever reader err of %v, input format %s", err, stringWithFormat)
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
            resMap[key] = bitSource
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
        log.Sugar().Infof("open file err of %v when Decompress, filepath %s", err, inputPath)
        return nil, err
    }
    binaryFlie, err := ioutil.ReadAll(file)
    if err != nil {
        log.Sugar().Infof("Read file err of %v when Decompress, filepath %s", err, inputPath)
        return nil, err
    }
    err = file.Close()
    if err != nil {
        log.Sugar().Infof("close file err of %v when Decompress, filepath %s", err, inputPath)
        return nil, err
    }
    return DecompressFromStream(inputPath, binaryFlie)
}
