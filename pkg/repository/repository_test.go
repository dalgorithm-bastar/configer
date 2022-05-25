package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/internal/datasource"
)

const (
	etcdConfigLocation = "../../unittestfiles/config/etcdClientv3.json"
	fileLocation       = "../../unittestfiles/config/"
)

func TestNewStorage(t *testing.T) {
	type args struct {
		dataSourceType string
		config         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// etcd type ok
		{
			name: "etcd type ok",
			args: args{
				dataSourceType: EtcdType,
				config:         etcdConfigLocation,
			},
			wantErr: false,
		},
		// etcd type failed
		{
			name: "etcd type fail",
			args: args{
				dataSourceType: EtcdType,
				config:         etcdConfigLocation,
			},
			wantErr: true,
		},
		// compressed file type succeed
		{
			name: "compress type ok",
			args: args{
				dataSourceType: CompressedFileType,
				config:         etcdConfigLocation,
			},
			wantErr: false,
		},
		// compressed file type failed
		{
			name: "compress type fail",
			args: args{
				dataSourceType: CompressedFileType,
				config:         etcdConfigLocation,
			},
			wantErr: true,
		},
		// normal dir type
		{
			name: "dir type ok",
			args: args{
				dataSourceType: DirType,
				config:         fileLocation,
			},
			wantErr: false,
		},
		// normal dir type
		{
			name: "dir type err",
			args: args{
				dataSourceType: DirType,
				config:         etcdConfigLocation,
			},
			wantErr: true,
		},
	}
	outputsNewEtcdType := []gomonkey.OutputCell{
		{Values: gomonkey.Params{&datasource.EtcdType{}, nil}},
		{Values: gomonkey.Params{nil, errors.New("init err")}},
	}
	patchesNewEtcdType := gomonkey.ApplyFuncSeq(datasource.NewEtcdType, outputsNewEtcdType)
	defer patchesNewEtcdType.Reset()
	outputsNewCompressedFileType := []gomonkey.OutputCell{
		{Values: gomonkey.Params{&datasource.CompressFileType{}, nil}},
		{Values: gomonkey.Params{nil, errors.New("init err")}},
	}
	patchesNewCompressedFileType := gomonkey.ApplyFuncSeq(datasource.NewCompressedFileType, outputsNewCompressedFileType)
	defer patchesNewCompressedFileType.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewStorage(context.Background(), tt.args.dataSourceType, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewStream(t *testing.T) {
	type args struct {
		fileMap map[string][]byte
	}
	tests := []struct {
		name string
		args args
		want Storage
	}{
		{
			name: "test",
			args: args{
				fileMap: map[string][]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewStream(tt.args.fileMap)
			//if got := NewStream(tt.args.fileMap); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewStream() = %v, want %v", got, tt.want)
			//}
		})
	}
}
