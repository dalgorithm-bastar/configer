package manage

import (
	"context"
	"errors"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestManager_DELETE(t *testing.T) {
	type fields struct {
		grpcInfo GrpcInfoStruct
		regExp   regExpStruct
	}
	type args struct {
		ctx    context.Context
		CfgReq *pb.CfgReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CfgResp
		wantErr bool
	}{
		{
			name:   "normal delete",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want: &pb.CfgResp{
				Status:    "ok",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "err case",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "",
				},
			},
			want: &pb.CfgResp{
				Status:    "delete err",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "nil req",
			fields: fields{},
			args: args{
				ctx:    context.TODO(),
				CfgReq: nil,
			},
			want: &pb.CfgResp{
				Status: "nil req deliverd",
			},
			wantErr: false,
		},
	}
	//outputs := []gomonkey.OutputCell{
	//	{
	//		Values: gomonkey.Params{
	//			nil, nil, "", nil,
	//		},
	//	},
	//	{
	//		Values: gomonkey.Params{
	//			errors.New("delete error"), nil, "", nil,
	//		},
	//	},
	//}
	//patches := gomonkey.ApplyFuncSeq(deleteInManager, outputs)
	//// 执行完毕后释放桩序列
	//defer patches.Reset()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().DeletebyPrefix(gomock.Any()).Return(nil),
		mockSrc.EXPECT().DeletebyPrefix(gomock.Any()).Return(errors.New("delete err")),
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			got, err := m.DELETE(tt.args.ctx, tt.args.CfgReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("DELETE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DELETE() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GET(t *testing.T) {
	type fields struct {
		grpcInfo GrpcInfoStruct
		regExp   regExpStruct
	}
	type args struct {
		ctx    context.Context
		CfgReq *pb.CfgReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CfgResp
		wantErr bool
	}{
		{
			name:   "get succeed",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want: &pb.CfgResp{
				Status:    "ok",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "get err",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "",
				},
			},
			want: &pb.CfgResp{
				Status:    "get error",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "nil req",
			fields: fields{},
			args: args{
				ctx:    context.TODO(),
				CfgReq: nil,
			},
			want: &pb.CfgResp{
				Status: "nil req deliverd",
			},
			wantErr: false,
		},
	}
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{
				nil, nil, "", nil,
			},
		},
		{
			Values: gomonkey.Params{
				errors.New("get error"), nil, "", nil,
			},
		},
	}
	patches := gomonkey.ApplyFuncSeq(Get, outputs)
	// 执行完毕后释放桩序列
	defer patches.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			got, err := m.GET(tt.args.ctx, tt.args.CfgReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("GET() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GET() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_POST(t *testing.T) {
	type fields struct {
		grpcInfo GrpcInfoStruct
		regExp   regExpStruct
	}
	type args struct {
		ctx    context.Context
		CfgReq *pb.CfgReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CfgResp
		wantErr bool
	}{
		{
			name:   "post succeed",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want: &pb.CfgResp{
				Status:    "ok",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "post err",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "",
				},
			},
			want: &pb.CfgResp{
				Status:    "post error",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "nil req",
			fields: fields{},
			args: args{
				ctx:    context.TODO(),
				CfgReq: nil,
			},
			want: &pb.CfgResp{
				Status: "nil req deliverd",
			},
			wantErr: false,
		},
	}
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{
				nil, nil, "", nil,
			},
		},
		{
			Values: gomonkey.Params{
				errors.New("post error"), nil, "", nil,
			},
		},
	}
	patches := gomonkey.ApplyFuncSeq(post, outputs)
	defer patches.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			got, err := m.POST(tt.args.ctx, tt.args.CfgReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("POST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("POST() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_PUT(t *testing.T) {
	type fields struct {
		grpcInfo GrpcInfoStruct
		regExp   regExpStruct
	}
	type args struct {
		ctx    context.Context
		CfgReq *pb.CfgReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.CfgResp
		wantErr bool
	}{
		{
			name:   "put succeed",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "chqr",
				},
			},
			want: &pb.CfgResp{
				Status:    "ok",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "put err",
			fields: fields{},
			args: args{
				ctx: context.TODO(),
				CfgReq: &pb.CfgReq{
					UserName: "",
				},
			},
			want: &pb.CfgResp{
				Status:    "put error",
				SliceData: nil,
				File: &pb.CompressedFile{
					FileName: "",
					FileData: nil,
				},
			},
			wantErr: false,
		},
		{
			name:   "nil req",
			fields: fields{},
			args: args{
				ctx:    context.TODO(),
				CfgReq: nil,
			},
			want: &pb.CfgResp{
				Status: "nil req deliverd",
			},
			wantErr: false,
		},
	}
	outputs := []gomonkey.OutputCell{
		{
			Values: gomonkey.Params{
				nil, nil, "", nil,
			},
		},
		{
			Values: gomonkey.Params{
				errors.New("put error"), nil, "", nil,
			},
		},
	}
	patches := gomonkey.ApplyFuncSeq(put, outputs)
	defer patches.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			got, err := m.PUT(tt.args.ctx, tt.args.CfgReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("PUT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PUT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
