package service

import (
    "context"
    "errors"
    "fmt"
    "reflect"
    "testing"

    "github.com/agiledragon/gomonkey/v2"
    "github.com/configcenter/internal/mock"
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
        want1 []string
        want2 string
        want3 []byte
    }{
        {
            name: "no username",
            args: args{
                req: &pb.CfgReq{
                    UserName: "",
                },
            },
            want:  errors.New("no username when getting"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "no target",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   nil,
                },
            },
            want:  errors.New("no target specified when getting"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "check err",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   []string{generation.NodeConfig},
                },
            },
            want:  errors.New("args error: too many cfgVersions or missing cfgVersions"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: generation.NodeConfig,
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   []string{generation.NodeConfig},
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.toml",
            want3: []byte("test"),
        },
        {
            name: generation.Templates,
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   []string{generation.Templates},
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.tar.gz",
            want3: []byte("test"),
        },
        {
            name: generation.Services,
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   []string{generation.Services},
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "servicelist.json",
            want3: []byte("test"),
        },
        {
            name: generation.Manipulations,
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    Target:   []string{generation.Manipulations},
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "manipulations.tar.gz",
            want3: []byte("test"),
        },
    }
    outputsGetNodeConfig := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "test.toml", []byte("test")}},
    }
    patchesGetNodeConfig := gomonkey.ApplyFuncSeq(getNodeConfig, outputsGetNodeConfig)
    defer patchesGetNodeConfig.Reset()
    outputsgetTemplates := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "test.tar.gz", []byte("test")}},
    }
    patchesgetTemplates := gomonkey.ApplyFuncSeq(getTemplates, outputsgetTemplates)
    defer patchesgetTemplates.Reset()
    outputsgetServices := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "servicelist.json", []byte("test")}},
    }
    patchesgetServices := gomonkey.ApplyFuncSeq(getServices, outputsgetServices)
    defer patchesgetServices.Reset()
    outputsgetManipulations := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "manipulations.tar.gz", []byte("test")}},
    }
    patchesgetManipulations := gomonkey.ApplyFuncSeq(getManipulations, outputsgetManipulations)
    defer patchesgetManipulations.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := Get(tt.args.ctx, tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Get() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("Get() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("Get() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_checkData(t *testing.T) {
    type args struct {
        data []*pb.CfgVersion
        tag  string
    }
    tests := []struct {
        name    string
        args    args
        wantErr bool
    }{
        {
            name: generation.NodeConfig + " nromal",
            args: args{
                tag: generation.NodeConfig,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "3204",
                                                LocalId:  "1",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
        {
            name: generation.NodeConfig + " nromal err",
            args: args{
                tag: generation.NodeConfig,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num:      "00",
                                Clusters: []*pb.Cluster{},
                            },
                        },
                    },
                },
            },
            wantErr: true,
        },
        {
            name: generation.NodeConfig + " globalid err",
            args: args{
                tag: generation.NodeConfig,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "fake",
                                                LocalId:  "1",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: true,
        },
        {
            name: generation.NodeConfig + " localid err",
            args: args{
                tag: generation.NodeConfig,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "3204",
                                                LocalId:  "fake",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: true,
        },
        {
            name: generation.Templates + " nromal",
            args: args{
                tag: generation.Templates,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "3204",
                                                LocalId:  "fake",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
        {
            name: generation.Templates + " nromal err",
            args: args{
                tag: generation.Templates,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num:      "00",
                                Clusters: []*pb.Cluster{},
                            },
                        },
                    },
                },
            },
            wantErr: true,
        },
        {
            name: generation.Clusters + " nromal",
            args: args{
                tag: generation.Clusters,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "3204",
                                                LocalId:  "fake",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
        {
            name: generation.Clusters + " nromal err",
            args: args{
                tag: generation.Clusters,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs:    []*pb.Environment{},
                    },
                },
            },
            wantErr: true,
        },
        {
            name: generation.Infrastructure + " nromal",
            args: args{
                tag: generation.Infrastructure,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num:      "00",
                                Clusters: []*pb.Cluster{},
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
        {
            name: generation.Infrastructure + " nromal err",
            args: args{
                tag:  generation.Infrastructure,
                data: []*pb.CfgVersion{},
            },
            wantErr: true,
        },
        {
            name: generation.Versions + " nromal",
            args: args{
                tag: generation.Versions,
                data: []*pb.CfgVersion{
                    {
                        Version: "0.0.1",
                        Envs: []*pb.Environment{
                            {
                                Num: "00",
                                Clusters: []*pb.Cluster{
                                    {
                                        ClusterName: "DTP.MC.set0",
                                        Nodes: []*pb.Node{
                                            {
                                                GlobalId: "3204",
                                                LocalId:  "fake",
                                                Template: "template1.toml",
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := checkData(tt.args.data, tt.args.tag); (err != nil) != tt.wantErr {
                t.Errorf("checkData() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func Test_createNodeConfig(t *testing.T) {
    repository.Src = *new(repository.Storage)
    type args struct {
        req         *pb.CfgReq
        tmplObj     string
        tmplContent []byte
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
            name: "normal1",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
                tmplObj:     "test.txt",
                tmplContent: []byte("test"),
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_00_DTP.MC.set0_1_test.txt",
            want3: []byte("test"),
        },
        {
            name: "normal2",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    File: &pb.CompressedFile{
                        FileName: "test1.txt",
                        FileData: nil,
                    },
                },
                tmplObj:     "test.txt",
                tmplContent: []byte("test"),
            },
            want:  nil,
            want1: nil,
            want2: "test1.txt",
            want3: []byte("test"),
        },
        {
            name: "init template err",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
                tmplObj:     "test.txt",
                tmplContent: []byte("test"),
            },
            want:  errors.New(fmt.Sprintf("no infrastructureData under path 0.0.1/infrastructure.json, please checkout in etcd or compressedfile")),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := createNodeConfig(tt.args.req, tt.args.tmplObj, tt.args.tmplContent)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("createNodeConfig() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("createNodeConfig() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("createNodeConfig() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("createNodeConfig() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getClusters(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: []string{"DTP.MC.set1", "EzEI.set0"},
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("No clusters under save path 0.0.1/00/clusters"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("DTP.MC.set1,EzEI.set0"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getClusters(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getClusters() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getClusters() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getClusters() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getClusters() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getDeploymentInfo(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_00_DTP.MC.set0_deploymentInfo.json",
            want3: []byte("test"),
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("No servicelist under path 0.0.1/00/DTP.MC.set0/servicelist.json"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "fill err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("fill err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
    )
    repository.Src = mockSrc
    outputsGetDeploymentInfo := []gomonkey.OutputCell{
        {Values: gomonkey.Params{"test", nil}},
        {Values: gomonkey.Params{nil, errors.New("fill err")}},
    }
    patchesGetDeploymentInfo := gomonkey.ApplyFuncSeq(generation.GetDeploymentInfo, outputsGetDeploymentInfo)
    defer patchesGetDeploymentInfo.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getDeploymentInfo(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getDeploymentInfo() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getDeploymentInfo() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getDeploymentInfo() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getDeploymentInfo() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getEnvironments(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: []string{"00", "88"},
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("No envs under save path 0.0.1/envs"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("00,88"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getEnvironments(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getEnvironments() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getEnvironments() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getEnvironments() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getEnvironments() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getInfrastructure(t *testing.T) {
    type args struct {
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
            name: "normal default name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_infrastructure.json.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "normal set name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "compress err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("compress err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("No infrastructure under save path 0.0.1/infrastructure.json"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    outputsCompressToStream := []gomonkey.OutputCell{
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{nil, errors.New("compress err")}},
    }
    patchesCompressToStream := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompressToStream)
    defer patchesCompressToStream.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getInfrastructure(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getInfrastructure() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getInfrastructure() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getInfrastructure() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getInfrastructure() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getKeyValueResult(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    File: &pb.CompressedFile{
                        FileName: "test.txt",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.txt",
            want3: []byte("test"),
        },
        {
            name: "bad req",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                },
            },
            want:  errors.New("nil File or FileContent transformed, please check your request"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getKeyValueResult(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getKeyValueResult() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getKeyValueResult() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getKeyValueResult() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getKeyValueResult() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getManipulations(t *testing.T) {
    type args struct {
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
            name: "normal default name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_00_DTP.MC.set0_manipulations.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "normal set name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "compress err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("compress err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("No manipulations under save path 0.0.1/00/DTP.MC.set0/manipulations/"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    outputsCompressToStream := []gomonkey.OutputCell{
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{nil, errors.New("compress err")}},
    }
    patchesCompressToStream := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompressToStream)
    defer patchesCompressToStream.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getManipulations(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getManipulations() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getManipulations() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getManipulations() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getManipulations() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getNodeConfig(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.toml",
            want3: []byte("test"),
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("No Template under save path 0.0.1/00/DTP.MC.set0/templates/template1.toml"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    outputsCreateNodeConfig := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "test.toml", []byte("test")}},
    }
    patchesCreateNodeConfig := gomonkey.ApplyFuncSeq(createNodeConfig, outputsCreateNodeConfig)
    defer patchesCreateNodeConfig.Reset()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getNodeConfig(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getNodeConfig() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getNodeConfig() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getNodeConfig() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getNodeConfig() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getNodeConfigPartlyOnline(t *testing.T) {
    type args struct {
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
            name: "get from local and server",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.toml",
                        FileData: []byte("test"),
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.toml",
            want3: []byte("test"),
        },
        {
            name: "get from server",
            args: args{
                req: &pb.CfgReq{
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "testforserver.toml",
            want3: []byte("testforserver"),
        },
    }
    outputsCreateNodeConfig := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "test.toml", []byte("test")}},
    }
    patchesCreateNodeConfig := gomonkey.ApplyFuncSeq(createNodeConfig, outputsCreateNodeConfig)
    defer patchesCreateNodeConfig.Reset()
    outputsgetNodeConfig := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, nil, "testforserver.toml", []byte("testforserver")}},
    }
    patchesgetNodeConfig := gomonkey.ApplyFuncSeq(getNodeConfig, outputsgetNodeConfig)
    defer patchesgetNodeConfig.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getNodeConfigPartlyOnline(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getNodeConfigPartlyOnline() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getNodeConfigPartlyOnline() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getNodeConfigPartlyOnline() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getNodeConfigPartlyOnline() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getServices(t *testing.T) {
    type args struct {
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
            name: "normal default name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_00_DTP.MC.set0_servicelist.json.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "normal set name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "compress err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("compress err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("No servicelist under save path 0.0.1/00/DTP.MC.set0/servicelist.json"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("test"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    outputsCompressToStream := []gomonkey.OutputCell{
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{nil, errors.New("compress err")}},
    }
    patchesCompressToStream := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompressToStream)
    defer patchesCompressToStream.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getServices(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getServices() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getServices() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getServices() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getServices() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getTemplates(t *testing.T) {
    type args struct {
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
            name: "normal default name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: nil,
            want2: "0.0.1_00_DTP.MC.set0_templates.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "normal set name",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  nil,
            want1: nil,
            want2: "test.tar.gz",
            want3: []byte("test"),
        },
        {
            name: "compress err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("compress err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                    File: &pb.CompressedFile{
                        FileName: "test.tar.gz",
                        FileData: nil,
                    },
                },
            },
            want:  errors.New("No Templates under save path 0.0.1/00/DTP.MC.set0/templates/"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(map[string][]byte{
            "test": []byte("test"),
        }, nil),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().GetbyPrefix(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    outputsCompressToStream := []gomonkey.OutputCell{
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{[]byte("test"), nil}},
        {Values: gomonkey.Params{nil, errors.New("compress err")}},
    }
    patchesCompressToStream := gomonkey.ApplyFuncSeq(util.CompressToStream, outputsCompressToStream)
    defer patchesCompressToStream.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getTemplates(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getTemplates() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getTemplates() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getTemplates() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getTemplates() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}

func Test_getVersions(t *testing.T) {
    type args struct {
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
            name: "normal",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  nil,
            want1: []string{"0.0.1", "0.0.2"},
            want2: "",
            want3: nil,
        },
        {
            name: "get err",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("get err"),
            want1: nil,
            want2: "",
            want3: nil,
        },
        {
            name: "get nil",
            args: args{
                req: &pb.CfgReq{
                    CfgVersions: []*pb.CfgVersion{
                        {
                            Version: "0.0.1",
                            Envs: []*pb.Environment{
                                {
                                    Num: "00",
                                    Clusters: []*pb.Cluster{
                                        {
                                            ClusterName: "DTP.MC.set0",
                                            Nodes: []*pb.Node{
                                                {
                                                    GlobalId: "3204",
                                                    LocalId:  "1",
                                                    Template: "template1.toml",
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    UserName: "chqr",
                },
            },
            want:  errors.New("No version saved in configcenter yet"),
            want1: nil,
            want2: "",
            want3: nil,
        },
    }
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockSrc := mock.NewMockStorage(ctrl)
    gomock.InOrder(
        mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("0.0.1,0.0.2"), nil),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
        mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
    )
    repository.Src = mockSrc
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, got1, got2, got3 := getVersions(tt.args.req)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("getVersions() got = %v, want %v", got, tt.want)
            }
            if !reflect.DeepEqual(got1, tt.want1) {
                t.Errorf("getVersions() got1 = %v, want %v", got1, tt.want1)
            }
            if got2 != tt.want2 {
                t.Errorf("getVersions() got2 = %v, want %v", got2, tt.want2)
            }
            if !reflect.DeepEqual(got3, tt.want3) {
                t.Errorf("getVersions() got3 = %v, want %v", got3, tt.want3)
            }
        })
    }
}
