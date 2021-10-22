// Package log 日志包，提供一层封装
package log

import (
	"encoding/json"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"io/ioutil"
	"os"
	"time"

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
	MaxSize      float64 `json:"maxsize"`
	MaxBackups   float64 `json:"maxbackups"`
	MaxAge       float64 `json:"maxage"`
}

// Logger Logger结构体封装zaplog、sugarlog、日志配置信息
type Logger struct {
	zapLog   *zap.Logger
	sugarLog *zap.SugaredLogger
	//zapFieldSlice []zapcore.Field //用于数据类型转换
	logInfo LogInfo //读取配置数据
}

func NewLogger(logConfigLocation string) error {
	logRec = new(Logger)
	//logRec.zapFieldSlice = make([]zapcore.Field, 0)
	file, err := os.Open(logConfigLocation)
	if err != nil {
		return err
	}
	binaryFlie, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(binaryFlie, &logRec.logInfo)
	if err != nil {
		return err
	}
	err = logRec.Init()
	if err != nil {
		return err
	}
	err = file.Close()
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
		rotatelogs.WithMaxAge(time.Hour*24*30),
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

func getEncoder(encoderType string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if encoderType == encodeTypeJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Zap() *zap.Logger {
	if logRec == nil || logRec.zapLog == nil {
		return zap.NewExample()
	}
	return logRec.zapLog
}

func Sugar() *zap.SugaredLogger {
	if logRec == nil || logRec.zapLog == nil {
		return zap.NewExample().Sugar()
	}
	return logRec.sugarLog
}

func GetLog() *Logger {
	return logRec
}