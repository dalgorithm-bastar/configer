package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/configcenter/pkg/pb"
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
		// TODO: Add test cases.
	}
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
