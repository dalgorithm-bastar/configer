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

func Test_checkFilePath(t *testing.T) {
    NewManager(context.Background(), grpcLocation)
    type args struct {
        fileMap map[string][]byte
    }
    tests := []struct {
        name    string
        args    args
        want    string
        want1   map[string][]string
        want2   map[string]map[string]string
        wantErr bool
    }{
        // multi root
        {
            name: "multi root",
            args: args{
                fileMap: map[string][]byte{
                    "0.0.1/00/DTP.MC.Set0/templates/template1.toml": []byte("test"),
                    "0.0.2/00/DTP.MC.Set0/templates/template1.toml": []byte("test"),
                },
            },
            want:    "",
            want1:   nil,
            want2:   nil,
            wantErr: true,
        },
        // keyword templates repeated
        {
            name: "keyword templates repeated",
            args: args{
                fileMap: map[string][]byte{
                    "0.0.1/00/DTP.MC.Set0/templates/template1.toml":        []byte("test"),
                    "0.0.1/00/DTP.MC.Set0/templates/template2.toml":        []byte("test"),
                    "0.0.1/templates/DTP.MC.Set0/templates/template2.toml": []byte("test"),
                },
            },
            want:    "",
            want1:   nil,
            want2:   nil,
            wantErr: true,
        },
        // normal of multi env and short path
        {
            name: "normal of multi env and short path",
            args: args{
                fileMap: map[string][]byte{
                    "0.0.1/infrastructure.json":                     []byte("test"),
                    "0.0.1/00/DTP.MC.Set0/templates/template1.toml": []byte("test00"),
                    "0.0.1/88/DTP.MC.Set0/templates/template2.toml": []byte("test88"),
                },
            },
            want: "0.0.1",
            want1: map[string][]string{
                "00": {"0.0.1/00/DTP.MC.Set0/templates/template1.toml"},
                "88": {"0.0.1/88/DTP.MC.Set0/templates/template2.toml"},
            },
            want2: map[string]map[string]string{
                "00": {"DTP.MC.Set0": ""},
                "88": {"DTP.MC.Set0": ""},
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, err := checkFilePath(tt.args.fileMap)
            if (err != nil) != tt.wantErr {
                t.Errorf("checkFilePath() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("checkFilePath() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("checkFilePath() got1 = %v, want %v", got1, tt.want1)
            }
            if !reflect.DeepEqual(got2, tt.want2) {
                t.Errorf("checkFilePath() got2 = %v, want %v", got2, tt.want2)
            }
        })
    }
}

func Test_put(t *testing.T) {
    NewManager(context.Background(), grpcLocation)
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
        // req illegal
        {
            name: "req illegal",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "",
                },
            },
            want:  errors.New("bad request, missing username or file or filedata"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        // normal
        {
            name: "normal",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.zip",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "",
            want3: nil,
        },
        // commit err
        {
            name: "commit err",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.zip",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  errors.New("commit err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        // get err
        {
            name: "get err",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.zip",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        // get nil
        {
            name: "get nil",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.zip",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  errors.New("get no data from filedata, please checkout file uploaded"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "delete err",
            args: args{
                ctx: context.TODO(),
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.zip",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    outputsDecompressFromStream := []gomonkey.OutputCell{
        {Values: gomonkey.Params{map[string][]byte{
            "0.0.1/00/DTP.MC.set0/templates/temp.toml": []byte("test"),
            "0.0.1/infrastructure.json":                []byte("test"),
        }, nil}},
        {Values: gomonkey.Params{map[string][]byte{
            "0.0.1/00/DTP.MC.set0/templates/temp.toml": []byte("test"),
            "0.0.1/infrastructure.json":                []byte("test"),
        }, nil}},
        {Values: gomonkey.Params{nil, errors.New("get err")}},
        {Values: gomonkey.Params{nil, nil}},
        {Values: gomonkey.Params{map[string][]byte{
            "0.0.1/00/DTP.MC.set0/templates/temp.toml": []byte("test"),
            "0.0.1/infrastructure.json":                []byte("test"),
        }, nil}},
    }
    patchesDecompressFromStream := gomonkey.ApplyFuncSeq(util.DecompressFromStream, outputsDecompressFromStream)
    defer patchesDecompressFromStream.Reset()

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
        mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
        mockSrc.EXPECT().AcidCommit(gomock.Any(), gomock.Any()).Return(errors.New("commit err")),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
    )
    repository.Src = mockSrc

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := put(tt.args.ctx, tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("put() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("put() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("put() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("put() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}
