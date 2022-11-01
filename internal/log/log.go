// Package log 日志包，提供一层封装
package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_logDir  string = "../log"
	_logTail string = ".log"
)

var (
	Logger  *zap.Logger
	Version string = time.Now().Format("20060102_150405")
)

func init() {
	err := os.MkdirAll(_logDir, os.FileMode(0755))
	if err != nil {
		fmt.Printf("mkdir for %s err:%v \n", _logDir, err)
		os.Exit(1)
	}
	f, err := os.OpenFile(_logDir+"/"+Version+_logTail, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		fmt.Printf("mkdir for %s err:%v \n", _logDir+"/"+Version, err)
		os.Exit(1)
	}
	err = f.Close()
	if err != nil {
		fmt.Printf("close file for %s err:%v", _logDir+"/"+Version+_logTail, err)
		os.Exit(1)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		//EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.InfoLevel)
	config := zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      "console",     // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		//InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", _logDir + "/" + Version + _logTail}, // 输出到指定文档 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr", _logDir + "/" + Version + _logTail},
	}
	// 构建日志
	logger, err := config.Build()
	Logger = logger
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}
}
