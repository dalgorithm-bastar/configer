// Package log 日志包，提供一层封装
package log

import (
    "os"
    "time"

    "github.com/configcenter/config"
    rotatelogs "github.com/lestrrat/go-file-rotatelogs"
    "github.com/spf13/viper"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

//仅在启动时初始化
const (
    debugLevel     = "debug"
    encodeTypeJson = "json"
)

var logRec *Logger

type LogInfo struct {
    LogPath      string  `json:"logpath"`
    RocordLevel  string  `json:"recordlevel"`
    EncodingType string  `json:"encodingtype"`
    FileName     string  `json:"filename"`
    MaxAge       float64 `json:"maxage"`
}

// Logger Logger结构体封装zaplog、sugarlog、日志配置信息
type Logger struct {
    zapLog   *zap.Logger
    sugarLog *zap.SugaredLogger
    //zapFieldSlice []zapcore.Field //用于数据类型转换
    logInfo LogInfo //读取配置数据
}

func NewLogger() error {
    logRec = new(Logger)
    logRec.setLogInfo()
    err := logRec.Init()
    if err != nil {
        return err
    }
    //创建日志文件夹
    err = os.MkdirAll(logRec.logInfo.LogPath, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}

func (l *Logger) Init() error {
    encoder := getEncoder(l.logInfo.EncodingType)

    // 保存日志30天，每24小时分割一次日志
    hook, err := rotatelogs.New(
        l.logInfo.LogPath+l.logInfo.FileName+"_%Y%m%d.log",
        rotatelogs.WithLinkName(l.logInfo.LogPath+l.logInfo.FileName),
        rotatelogs.WithMaxAge(time.Hour*24*time.Duration(l.logInfo.MaxAge)),
        rotatelogs.WithRotationTime(time.Hour*24),
    )
    if err != nil {
        return err
    }
    /*lumberJackLogger := &lumberjack.Logger{
    	Filename:   l.logInfo.LogPath + l.logInfo.FileName,
    	MaxSize:    int(l.logInfo.MaxSize), //MB
    	MaxBackups: int(l.logInfo.MaxBackups),
    	MaxAge:     int(l.logInfo.MaxAge), //day
    	Compress:   false,
    }*/
    writeSyncer := zapcore.AddSync(hook)
    level := zapcore.InfoLevel
    if l.logInfo.RocordLevel == debugLevel {
        level = zapcore.DebugLevel
    }
    core := zapcore.NewCore(encoder, writeSyncer, level)
    l.zapLog = zap.New(core, zap.AddCaller())
    l.sugarLog = l.zapLog.Sugar()
    return nil
}

func (l *Logger) setLogInfo() {
    l.logInfo = LogInfo{
        LogPath:      viper.GetString(config.LogLogPath),
        RocordLevel:  viper.GetString(config.LogRecordLevel),
        EncodingType: viper.GetString(config.LogEncodingType),
        FileName:     viper.GetString(config.LogFileName),
        MaxAge:       viper.GetFloat64(config.LogMaxAge),
    }
}

func getEncoder(encoderType string) zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    if encoderType == encodeTypeJson {
        return zapcore.NewJSONEncoder(encoderConfig)
    }
    return zapcore.NewConsoleEncoder(encoderConfig)
}

// Zap 获取结构化输出日志实例
func Zap() *zap.Logger {
    if logRec == nil || logRec.zapLog == nil {
        return zap.NewExample()
    }
    return logRec.zapLog
}

// Sugar 获取格式化输出日志实例
func Sugar() *zap.SugaredLogger {
    if logRec == nil || logRec.zapLog == nil {
        return zap.NewExample().Sugar()
    }
    return logRec.sugarLog
}

// GetLog 获取日志实例
func GetLog() *Logger {
    return logRec
}
