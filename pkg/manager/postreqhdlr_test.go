package manage

import (
	"context"
	"errors"
	"fmt"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_checkVersionFormat(t *testing.T) {
	NewManager(grpcLocation)
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "err version format",
			args: args{
				version: "01.10.0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no version yet",
			args: args{
				version: "1.0.1",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "get err",
			args: args{
				version: "1.0.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "version repeat",
			args: args{
				version: "0.0.1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				version: "1.0.1",
			},
			want:    []string{"0.0.1", "0.0.2"},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),                     //当前没有版本
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get error")), //获取出错
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,0.0.2"), nil),   //版本号重复
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,0.0.2"), nil),   //正常
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkVersionFormat(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkVersionFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkVersionFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createVersion(t *testing.T) {
	NewManager(grpcLocation)
	type args struct {
		fileMap        map[string][]byte
		formerVersions []string
		version        string
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []string
		want2 string
		want3 []byte
	}{
		//路径错误
		{
			name: "filepath err",
			args: args{
				fileMap: map[string][]byte{
					"": []byte("test"),
				},
				formerVersions: nil,
				version:        "",
			},
			want:  errors.New(fmt.Sprintf("error path: %s", "")),
			want1: nil,
			want2: "",
			want3: nil,
		},
		//提交失败
		{
			name: "commit err",
			args: args{
				fileMap: map[string][]byte{
					"0.0.1/infrastructure.json": []byte("test"),
				},
				formerVersions: nil,
				version:        "",
			},
			want:  errors.New("commit err"),
			want1: nil,
			want2: "",
			want3: nil,
		},
		//正常
		{
			name: "normal",
			args: args{
				fileMap: map[string][]byte{
					"0.0.1/infrastructure.json": []byte("test"),
				},
				formerVersions: nil,
				version:        "0.0.1",
			},
			want:  nil,
			want1: []string{"0.0.1"},
			want2: "",
			want3: nil,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(errors.New("commit err")),
		mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(nil),
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := createVersion(tt.args.fileMap, tt.args.formerVersions, tt.args.version)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createVersion() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("createVersion() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("createVersion() got2 = %v, want %v", got2, tt.want2)
			}
			if !reflect.DeepEqual(got3, tt.want3) {
				t.Errorf("createVersion() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_getDataByUserName(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		// get err
		{
			name: "get err",
			args: args{
				userName: "chqr",
			},
			want:    nil,
			wantErr: true,
		},
		// get nil
		{
			name: "get nil",
			args: args{
				userName: "chqr",
			},
			want:    nil,
			wantErr: false,
		},
		// normal
		{
			name: "normal",
			args: args{
				userName: "chqr",
			},
			want: map[string][]byte{
				"0.0.1/infrastructure.json": []byte("test"),
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
			"0.0.1/infrastructure.json": []byte("test"),
		}, nil),
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDataByUserName(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDataByUserName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataByUserName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextVersion(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		want1   string
		wantErr bool
	}{
		// get err
		{
			name:    "get err",
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		// get nil
		{
			name:    "get nil",
			want:    nil,
			want1:   "0.0.1",
			wantErr: false,
		},
		// atoi err
		{
			name:    "atoi err",
			want:    nil,
			want1:   "",
			wantErr: true,
		},
		// normal
		{
			name:    "normal",
			want:    []string{"0.0.1", "0.0.2"},
			want1:   "0.0.3",
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(""), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,0.0.2"), nil),
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getNextVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("getNextVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNextVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getNextVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_post(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CfgReq
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 []string
		want2 string
		want3 []byte
	}{
		// username err
		{
			name: "username err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "",
				},
			},
			want:  errors.New("empty username, commit request deny"),
			want1: nil,
			want2: "",
			want3: nil,
		},
		// getdata err
		{
			name: "get data err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want:  errors.New("get data err"),
			want1: nil,
			want2: "",
			want3: nil,
		},
		// filemap err
		{
			name: "filemap err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want:  errors.New(fmt.Sprintf("no data under username %s, commit request deny", "chqr")),
			want1: nil,
			want2: "",
			want3: nil,
		},
		// get next version err
		{
			name: "get next version err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want:  errors.New("get next version err"),
			want1: nil,
			want2: "",
			want3: nil,
		},
		// check version format err
		{
			name: "check version format err",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{
					UserName: "chqr",
					Target:   []string{"0.0.4"},
				},
			},
			want:  errors.New("check version format err"),
			want1: nil,
			want2: "",
			want3: nil,
		},
	}
	outputsGetData := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("get data err")}},
		{Values: gomonkey.Params{nil, nil}},
		{Values: gomonkey.Params{map[string][]byte{
			"0.0.1/infrastructure.json": []byte("test"),
		}, nil}},
		{Values: gomonkey.Params{map[string][]byte{
			"0.0.1/infrastructure.json": []byte("test"),
		}, nil}},
	}
	patchesGetData := gomonkey.ApplyFuncSeq(getDataByUserName, outputsGetData)
	defer patchesGetData.Reset()
	outputsGetNextVersion := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, "", errors.New("get next version err")}},
	}
	patchesGetNextVersion := gomonkey.ApplyFuncSeq(getNextVersion, outputsGetNextVersion)
	defer patchesGetNextVersion.Reset()
	outputsCheckVersionFormat := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("check version format err")}},
	}
	patchesCheckVersionFormat := gomonkey.ApplyFuncSeq(checkVersionFormat, outputsCheckVersionFormat)
	defer patchesCheckVersionFormat.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := post(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("post() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("post() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("post() got2 = %v, want %v", got2, tt.want2)
			}
			if !reflect.DeepEqual(got3, tt.want3) {
				t.Errorf("post() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
