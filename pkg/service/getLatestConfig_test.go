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
    "github.com/golang/mock/gomock"
    "gopkg.in/yaml.v3"
)

func Test_getLatestConfigByEnvNum(t *testing.T) {
    genSrc := GenSrc{
        UserName:       "qrchen",
        Version:        "3.1.0",
        Scheme:         "scheme1",
        EnvNum:         "02",
        Ip:             "1.1.1.1,2.2.2.2",
        CastPort:       "1024,65535",
        TcpPort:        "1024,65535",
        Infrastructure: "test",
    }
    genSrc2 := GenSrc{
        UserName:       "qrchen",
        Version:        "3.1.0",
        Scheme:         "scheme1",
        EnvNum:         "03",
        Ip:             "1.1.1.1,2.2.2.2",
        CastPort:       "1024,65535",
        TcpPort:        "1024,65535",
        Infrastructure: "test",
    }
    genGrp := GenSrcGrp{
        GenSrcs: []GenSrc{genSrc, genSrc2},
    }
    genGrpBytes, err := yaml.Marshal(genGrp)
    if err != nil {
        panic(err)
    }
    _ = NewManager(context.Background())
    type args struct {
        ctx context.Context
        req *pb.EnvNumReq
    }
    tests := []struct {
        name  string
        args  args
        want  error
        want1 []*pb.VersionInfo
        want2 *pb.AnyFile
    }{
        {
            name: "num err",
            args: args{
                ctx: context.Background(),
                req: &pb.EnvNumReq{
                    EnvNum: "000",
                },
            },
            want:  errors.New("envNum format err, input envNum:000, please input envNum of 2 bit number"),
            want1: nil,
            want2: nil,
        },
        {
            name: "get err",
            args: args{
                ctx: context.Background(),
                req: &pb.EnvNumReq{
                    EnvNum: "01",
                },
            },
            want:  errors.New("get envfile err or envfile not exist yet, err:get err; envfile:[]"),
            want1: nil,
            want2: nil,
        },
        {
            name: "ok",
            args: args{
                ctx: context.Background(),
                req: &pb.EnvNumReq{
                    EnvNum: "02",
                },
            },
            want:  nil,
            want1: nil,
            want2: nil,
        },
        {
            name: "not gen yet",
            args: args{
                ctx: context.Background(),
                req: &pb.EnvNumReq{
                    EnvNum: "05",
                },
            },
            want:  errors.New("envNum:05 has not been generated and recorded on remote yet, please generate first"),
            want1: nil,
            want2: nil,
        },
    }

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    repository.Src = mockSrc
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(genGrpBytes, nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(genGrpBytes, nil),
    )

    outputsCfg := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, nil}},
        {Values: gomonkey.Params{}},
    }
    patchesCfg := gomonkey.ApplyFuncSeq(getConfig, outputsCfg)
    defer patchesCfg.Reset()

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2 := getLatestConfigByEnvNum(tt.args.ctx, tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getLatestConfigByEnvNum() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getLatestConfigByEnvNum() got1 = %v, want %v", got1, tt.want1)
            }
            if !reflect.DeepEqual(got2, tt.want2) {
                t.Errorf("getLatestConfigByEnvNum() got2 = %v, want %v", got2, tt.want2)
            }
        })
    }
}
