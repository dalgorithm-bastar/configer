package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/generation"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/golang/mock/gomock"
)

func Test_commit(t *testing.T) {
	type args struct {
		ctxRoot context.Context
		req     *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "err username",
			args: args{
				ctxRoot: context.Background(),
				req:     &pb.CfgReq{UserName: "qr,chen"},
			},
			want:  errors.New("empty username of username contains point or comma, commit request deny"),
			want1: nil,
			want2: nil,
		},
		{
			name: "err infra data",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetInfrastructure,
				},
			},
			want:  errors.New("infrastructure.yaml with err content,please checkout"),
			want1: nil,
			want2: nil,
		},
		{
			name: "infra ok",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetInfrastructure,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
		{
			name: "infra put err",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetInfrastructure,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  errors.New("put infra err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get raw err",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
				},
			},
			want:  errors.New("get raw err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get raw nil",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
				},
			},
			want:  errors.New("no rawData on remote yet at path:qrchen"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get next version nil",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
				},
			},
			want:  errors.New("get next version err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get infra err",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
				},
			},
			want:  errors.New("get infra err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "generate err",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
					ArgRange: &pb.ArgRange{
						TopicIp:   nil,
						TopicPort: nil,
					},
				},
			},
			want:  errors.New("generate cfgFile err:generate err,please checkout cache under username:qrchen"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get source operator err",
			args: args{
				ctxRoot: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					File: &pb.AnyFile{
						FileName: "test",
						FileData: []byte("testdata"),
					},
					Version: "",
					ArgRange: &pb.ArgRange{
						TopicIp:   nil,
						TopicPort: nil,
					},
				},
			},
			want:  errors.New("get etcd clientv3 err, type of datasource: string"),
			want1: nil,
			want2: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil),
		mockSrc.EXPECT().Put(gomock.Any(), gomock.Any()).Return(errors.New("put infra err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get raw err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetSourceDataorOperator().Return("string"),
	)

	outputsVersion := []gomonkey.OutputCell{
		{Values: gomonkey.Params{"", "", nil}},
		{Values: gomonkey.Params{"", "", nil}},
		{Values: gomonkey.Params{"", "", errors.New("get next version err")}},
		{Values: gomonkey.Params{"", "", nil}},
		{Values: gomonkey.Params{"", "", nil}},
		{Values: gomonkey.Params{"", "", nil}},
	}
	patchesVersion := gomonkey.ApplyFuncSeq(getNextVersion, outputsVersion)
	defer patchesVersion.Reset()

	outputsInfra := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get infra err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, &pb.AnyFile{FileName: "infrastructure.yaml", FileData: []byte("test")}}},
		{Values: gomonkey.Params{nil, nil, &pb.AnyFile{FileName: "infrastructure.yaml", FileData: []byte("test")}}},
	}
	patchesInfra := gomonkey.ApplyFuncSeq(getInfrastructure, outputsInfra)
	defer patchesInfra.Reset()

	outputsGen := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("generate err")}},
		{Values: gomonkey.Params{nil, nil}},
	}
	patchesGen := gomonkey.ApplyFuncSeq(generation.Generate, outputsGen)
	defer patchesGen.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := commit(tt.args.ctxRoot, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("commit() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("commit() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("commit() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getNextVersion(t *testing.T) {
	_ = NewManager(context.Background())
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name:    "err version",
			args:    args{version: "112.3"},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name:    "init state",
			args:    args{version: ""},
			want:    "",
			wantErr: false,
			want1:   "0.0.1",
		},
		{
			name:    "init point",
			args:    args{version: "1.0.0"},
			want:    "",
			wantErr: false,
			want1:   "1.0.0",
		},
		{
			name:    "already point",
			args:    args{version: "1.0.0"},
			want:    "0.0.1,qr,t1,0.0.2,qr,t2",
			wantErr: false,
			want1:   "1.0.0",
		},
		{
			name:    "reapeat point",
			args:    args{version: "0.0.2"},
			want:    "",
			wantErr: true,
			want1:   "",
		},
		{
			name:    "already no point",
			args:    args{version: ""},
			want:    "0.0.1,qr,t1,0.0.2,qr,t2",
			wantErr: false,
			want1:   "0.0.3",
		},
		{
			name:    "get err",
			args:    args{version: ""},
			want:    "",
			wantErr: true,
			want1:   "",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,qr,t1,0.0.2,qr,t2"), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,qr,t1,0.0.2,qr,t2"), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,qr,t1,0.0.2,qr,t2"), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getNextVersion(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNextVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNextVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getNextVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
