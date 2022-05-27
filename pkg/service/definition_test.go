package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/configcenter/pkg/pb"
)

func TestCheckFlag(t *testing.T) {
	type args struct {
		flagName string
		value    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "action",
			args: args{
				flagName: _reqAction,
				value:    _actionAddAndReplaceAtLeaf,
			},
			want: true,
		},
		{
			name: "type",
			args: args{
				flagName: _reqType,
				value:    _typeDeployment,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckFlag(tt.args.flagName, tt.args.value); got != tt.want {
				t.Errorf("CheckFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGrpcInfo(t *testing.T) {
	_ = NewManager(context.Background())
	tests := []struct {
		name string
		want *GrpcInfoStruct
	}{
		{
			name: "ok",
			want: &manager.grpcInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGrpcInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGrpcInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManager(t *testing.T) {
	_ = NewManager(context.Background())
	tests := []struct {
		name string
		want *Manager
	}{
		{
			name: "ok",
			want: manager,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_COMMIT(t *testing.T) {
	_ = NewManager(context.Background())
	type fields struct {
		ctx      context.Context
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
			name: "nil req",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: nil,
			},
			want:    &pb.CfgResp{Status: "nil req deliverd"},
			wantErr: false,
		},
		{
			name: "commit err",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "commit err",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
		{
			name: "ok",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "ok",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
	}

	outputsCommit := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("commit err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, nil}},
	}
	patchesCommit := gomonkey.ApplyFuncSeq(commit, outputsCommit)
	defer patchesCommit.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				ctx:      tt.fields.ctx,
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			got, err := m.COMMIT(tt.args.ctx, tt.args.CfgReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("COMMIT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("COMMIT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_DELETE(t *testing.T) {
	_ = NewManager(context.Background())
	type fields struct {
		ctx      context.Context
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
			name: "nil req",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: nil,
			},
			want:    &pb.CfgResp{Status: "nil req deliverd"},
			wantErr: false,
		},
		{
			name: "delete err",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "delete err",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
		{
			name: "ok",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "ok",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
	}

	outputsDelete := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("delete err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, nil}},
	}
	patchesDelete := gomonkey.ApplyFuncSeq(DeleteInManager, outputsDelete)
	defer patchesDelete.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				ctx:      tt.fields.ctx,
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
	_ = NewManager(context.Background())
	type fields struct {
		ctx      context.Context
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
			name: "nil req",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: nil,
			},
			want:    &pb.CfgResp{Status: "nil req deliverd"},
			wantErr: false,
		},
		{
			name: "get err",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "get err",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
		{
			name: "ok",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "ok",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
	}

	outputsGet := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("get err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, nil}},
	}
	patchesGet := gomonkey.ApplyFuncSeq(Get, outputsGet)
	defer patchesGet.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				ctx:      tt.fields.ctx,
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

func TestManager_PUT(t *testing.T) {
	_ = NewManager(context.Background())
	type fields struct {
		ctx      context.Context
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
			name: "nil req",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: nil,
			},
			want:    &pb.CfgResp{Status: "nil req deliverd"},
			wantErr: false,
		},
		{
			name: "put err",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "put err",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
		{
			name: "ok",
			fields: fields{
				ctx:      manager.ctx,
				grpcInfo: manager.grpcInfo,
				regExp:   manager.regExp,
			},
			args: args{
				ctx:    context.Background(),
				CfgReq: &pb.CfgReq{},
			},
			want: &pb.CfgResp{
				Status:      "ok",
				VersionList: nil,
				File:        nil,
			},
			wantErr: false,
		},
	}

	outputsPut := []gomonkey.OutputCell{
		{Values: gomonkey.Params{errors.New("put err"), nil, nil}},
		{Values: gomonkey.Params{nil, nil, nil}},
	}
	patchesPut := gomonkey.ApplyFuncSeq(put, outputsPut)
	defer patchesPut.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				ctx:      tt.fields.ctx,
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

func TestManager_SetGrpcInfo(t *testing.T) {
	type fields struct {
		ctx      context.Context
		grpcInfo GrpcInfoStruct
		regExp   regExpStruct
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				ctx:      tt.fields.ctx,
				grpcInfo: tt.fields.grpcInfo,
				regExp:   tt.fields.regExp,
			}
			m.SetGrpcInfo()
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		ctxIn context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewManager(tt.args.ctxIn); (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
