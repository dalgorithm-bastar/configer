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
	"github.com/configcenter/pkg/util"
	"github.com/golang/mock/gomock"
)

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "check data err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "",
				},
			},
			want:  errors.New("no username when getting"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get cfg err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetConfig,
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("get cfg err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get infra err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetInfrastructure,
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("get infra err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get version err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetVersion,
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("get version err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get cache err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetCache,
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("get cache err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "get raw err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   TargetRaw,
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("get raw err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "target err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Target:   "err",
					EnvNum:   "01",
					Version:  "1.0.0",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("wrong arg of type, target"),
			want1: nil,
			want2: nil,
		},
	}

	outputsCfg := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get cfg err"), nil, nil}},
	}
	patchesCfg := gomonkey.ApplyFuncSeq(getConfig, outputsCfg)
	defer patchesCfg.Reset()

	outputsInfra := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get infra err"), nil, nil}},
	}
	patchesInfra := gomonkey.ApplyFuncSeq(getInfrastructure, outputsInfra)
	defer patchesInfra.Reset()

	outputsVersion := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get version err"), nil, nil}},
	}
	patchesVersion := gomonkey.ApplyFuncSeq(getVersion, outputsVersion)
	defer patchesVersion.Reset()

	outputsCache := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get cache err"), nil, nil}},
	}
	patchesCache := gomonkey.ApplyFuncSeq(getCache, outputsCache)
	defer patchesCache.Reset()

	outputsRaw := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get raw err"), nil, nil}},
	}
	patchesRaw := gomonkey.ApplyFuncSeq(getRaw, outputsRaw)
	defer patchesRaw.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := Get(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Get() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_checkData(t *testing.T) {
	type args struct {
		data *pb.CfgReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil target",
			args: args{data: &pb.CfgReq{
				UserName: "qrchen",
				Target:   "",
			}},
			wantErr: true,
		},
		{
			name: "cfg env err",
			args: args{data: &pb.CfgReq{
				UserName: "qrchen",
				Target:   TargetConfig,
				EnvNum:   "000",
			}},
			wantErr: true,
		},
		{
			name: "cfg range err",
			args: args{data: &pb.CfgReq{
				UserName: "qrchen",
				Target:   TargetConfig,
				EnvNum:   "00",
			}},
			wantErr: true,
		},
		{
			name: "cfg version err",
			args: args{data: &pb.CfgReq{
				UserName: "qrchen",
				Target:   TargetConfig,
				EnvNum:   "00",
				ArgRange: &pb.ArgRange{
					TopicIp:   []string{"1.1.1.1,2.2.2.2"},
					TopicPort: []string{"10000", "55555"},
				},
				Version: "1..0",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("checkData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCache(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "get err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("get err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get nil",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("no cache on remote yet at path:qrchen"),
			want2: nil,
			want1: nil,
		},
		{
			name: "compress err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("compress err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want: nil,
			want2: &pb.AnyFile{
				FileName: "cache.tar.gz",
				FileData: nil,
			},
			want1: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
	)

	outputsCompress := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("compress err")}},
		{Values: gomonkey.Params{nil, nil}},
	}
	patchesCompress := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompress)
	defer patchesCompress.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getCache(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCache() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getCache() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getCache() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getConfig(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "get err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("get err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get nil",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
					Scheme:   "scheme1",
				},
			},
			want:  errors.New("no rawData on remote yet at path:1.0.0/scheme1"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get infra err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
					Scheme:   "scheme1",
				},
			},
			want:  errors.New("get infra err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "generate err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
					Scheme:   "scheme1",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("generate err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "compress err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
					Scheme:   "scheme1",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want:  errors.New("compress err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
					Scheme:   "scheme1",
					ArgRange: &pb.ArgRange{
						TopicIp:   []string{"1.1.1.1,2.2.2.2"},
						TopicPort: []string{"10000", "55555"},
					},
				},
			},
			want: nil,
			want2: &pb.AnyFile{
				FileName: "config.tar.gz",
				FileData: nil,
			},
			want1: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
	)

	outputsInfra := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get infra err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, &pb.AnyFile{
			FileName: repository.Infrastructure,
			FileData: []byte("test"),
		}}},
		{Values: gomonkey.Params{nil, nil, &pb.AnyFile{
			FileName: repository.Infrastructure,
			FileData: []byte("test"),
		}}},
		{Values: gomonkey.Params{nil, nil, &pb.AnyFile{
			FileName: repository.Infrastructure,
			FileData: []byte("test"),
		}}},
	}
	patchesInfra := gomonkey.ApplyFuncSeq(getInfrastructure, outputsInfra)
	defer patchesInfra.Reset()

	outputsGen := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("generate err")}},
		{Values: gomonkey.Params{nil, nil}},
		{Values: gomonkey.Params{nil, nil}},
	}
	patchesGen := gomonkey.ApplyFuncSeq(generation.Generate, outputsGen)
	defer patchesGen.Reset()

	outputsCompress := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("compress err")}},
		{Values: gomonkey.Params{nil, nil}},
	}
	patchesCompress := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompress)
	defer patchesCompress.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getConfig(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getConfig() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getConfig() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getInfrastructure(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "get err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("get err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get nil",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("no infrastructure on remote yet"),
			want2: nil,
			want1: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want: nil,
			want2: &pb.AnyFile{
				FileName: repository.Infrastructure,
				FileData: []byte("test"),
			},
			want1: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getInfrastructure(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInfrastructure() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getInfrastructure() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getInfrastructure() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getRaw(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "get err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("get err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get nil",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
				},
			},
			want:  errors.New("no rawData on remote yet at path:1.0.0"),
			want2: nil,
			want1: nil,
		},
		{
			name: "compress err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("compress err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want: nil,
			want2: &pb.AnyFile{
				FileName: "rawdata.tar.gz",
				FileData: nil,
			},
			want1: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
	)

	outputsCompress := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("compress err")}},
		{Values: gomonkey.Params{nil, nil}},
	}
	patchesCompress := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompress)
	defer patchesCompress.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getRaw(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRaw() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getRaw() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getRaw() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getVersion(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []*pb.VersionInfo
		want2 *pb.AnyFile
	}{
		{
			name: "get err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{UserName: "qrchen"},
			},
			want:  errors.New("get err"),
			want2: nil,
			want1: nil,
		},
		{
			name: "get nil",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
				},
			},
			want:  errors.New("no versions on remote yet"),
			want2: nil,
			want1: nil,
		},
		{
			name: "data err",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
				},
			},
			want:  errors.New("error happened on remote, with unexpected version data: map[test:[116 101 115 116]]"),
			want2: nil,
			want1: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					Version:  "1.0.0",
				},
			},
			want:  nil,
			want2: nil,
			want1: []*pb.VersionInfo{{
				Name: "1.0.0",
				User: "qrchen",
				Time: "time",
			}},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"test": []byte("test")}, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{repository.Versions: []byte("1.0.0,qrchen,time")}, nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getVersion(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVersion() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getVersion() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("getVersion() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
