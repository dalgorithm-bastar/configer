package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"github.com/golang/mock/gomock"
)

func TestDeleteInManager(t *testing.T) {
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
			name: "ok",
			args: args{
				ctx: context.TODO(),
				req: &pb.CfgReq{},
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	repository.Src = mockSrc
	gomock.InOrder(
		mockSrc.EXPECT().DeletebyPrefix(gomock.Any()).Return(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := DeleteInManager(tt.args.ctx, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteInManager() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("DeleteInManager() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("DeleteInManager() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
