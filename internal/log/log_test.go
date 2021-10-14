package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"reflect"
	"testing"
)

const logConfigLocation = "../../config/log.json"

func TestLogger_Init(t *testing.T) {
	type fields struct {
		zapLog   *zap.Logger
		sugarLog *zap.SugaredLogger
		logInfo  LogInfo
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				zapLog:   nil,
				sugarLog: nil,
				logInfo: LogInfo{
					EncodingType: encodeTypeJson,
				},
			},
			wantErr: false,
		},
		{
			name: "debuglevel",
			fields: fields{
				zapLog:   nil,
				sugarLog: nil,
				logInfo: LogInfo{
					RocordLevel: debugLevel,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				zapLog:   tt.fields.zapLog,
				sugarLog: tt.fields.sugarLog,
				logInfo:  tt.fields.logInfo,
			}
			if err := l.Init(); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	type args struct {
		logConfigLocation string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				logConfigLocation: logConfigLocation,
			},
			wantErr: false,
		},
		{
			name: "open err",
			args: args{
				logConfigLocation: "",
			},
			wantErr: true,
		},
		{
			name: "unmarshall err",
			args: args{
				logConfigLocation: "log.go",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewLogger(tt.args.logConfigLocation); (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSugar(t *testing.T) {
	tests := []struct {
		name string
		want *zap.SugaredLogger
	}{
		{
			name: "nil",
			want: zap.NewExample().Sugar(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = Sugar()
			/*if got := Sugar(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sugar() = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func TestZap(t *testing.T) {
	tests := []struct {
		name string
		want *zap.Logger
	}{
		{
			name: "nil",
			want: zap.NewExample(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = Zap()
			/*if got := Zap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Zap() = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func Test_getEncoder(t *testing.T) {
	type args struct {
		encoderType string
	}
	tests := []struct {
		name string
		args args
		want zapcore.Encoder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEncoder(tt.args.encoderType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}
