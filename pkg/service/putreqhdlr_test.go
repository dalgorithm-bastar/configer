package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/configcenter/pkg/util"
	"github.com/golang/mock/gomock"
)

func Test_checkAndReplaceRootPath(t *testing.T) {
	type args struct {
		fileMap  map[string][]byte
		RootPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkAndReplaceRootPath(tt.args.fileMap, tt.args.RootPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkAndReplaceRootPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAndReplaceRootPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_put(t *testing.T) {
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
			name: "user name nil",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "",
				},
			},
			want:  errors.New("bad request, missing username or file or filedata"),
			want1: nil,
			want2: nil,
		},
		{
			name: "decompress err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					File: &pb.AnyFile{
						FileName: "testname",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  errors.New("get data err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "decompress nil",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					File: &pb.AnyFile{
						FileName: "testname",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  errors.New("get no data from filedata, please checkout file uploaded"),
			want1: nil,
			want2: nil,
		},
		{
			name: "ok",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					File: &pb.AnyFile{
						FileName: "testname",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
		{
			name: "get err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					File: &pb.AnyFile{
						FileName: "testname",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  errors.New("get err"),
			want1: nil,
			want2: nil,
		},
		{
			name: "commit err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "qrchen",
					File: &pb.AnyFile{
						FileName: "testname",
						FileData: []byte("testdata"),
					},
				},
			},
			want:  errors.New("commit err"),
			want1: nil,
			want2: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"olddata": []byte("olddata")}, nil),
		mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{"olddata": []byte("olddata")}, nil),
		mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(errors.New("commit err")),
	)

	outputsDeCompress := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("get data err")}},
		{Values: gomonkey.Params{nil, nil}},
		{Values: gomonkey.Params{map[string][]byte{
			"0.0.1/infrastructure.json":    []byte("test1"),
			"0.0.1/00/servicelist.json":    []byte("test2"),
			"0.0.1/00/templates/temp.toml": []byte("test3"),
		}, nil}},
		{Values: gomonkey.Params{map[string][]byte{
			"0.0.1/infrastructure.json":    []byte("test1"),
			"0.0.1/00/servicelist.json":    []byte("test2"),
			"0.0.1/00/templates/temp.toml": []byte("test3"),
		}, nil}},
		{Values: gomonkey.Params{map[string][]byte{
			"0.0.1/infrastructure.json":    []byte("test1"),
			"0.0.1/00/servicelist.json":    []byte("test2"),
			"0.0.1/00/templates/temp.toml": []byte("test3"),
		}, nil}},
	}
	patchesDeCompress := gomonkey.ApplyFuncSeq(util.DecompressFromStream, outputsDeCompress)
	defer patchesDeCompress.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := put(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("put() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("put() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("put() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
