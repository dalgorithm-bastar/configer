package manage

import (
    "context"
    "reflect"
    "testing"

    "github.com/configcenter/internal/mock"
    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
    "github.com/golang/mock/gomock"
)

func Test_deleteInManager(t *testing.T) {
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
        {
            name: "normal test",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().DeletebyPrefix(gomock.Any()).Return(nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := deleteInManager(tt.args.ctx, tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("delete() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("delete() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("delete() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("delete() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}
