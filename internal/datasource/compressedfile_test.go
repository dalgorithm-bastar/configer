package datasource

import (
    "errors"
    "reflect"
    "testing"

    "github.com/agiledragon/gomonkey/v2"
    "github.com/configcenter/pkg/util"
)

func TestCompressFileType_AcidCommit(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        putMap    map[string]string
        deleteMap []string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {
            name: "acidcommit",
            fields: fields{
                path: "",
                data: nil,
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            if err := c.AtomicCommit(tt.args.putMap, tt.args.deleteMap); (err != nil) != tt.wantErr {
                t.Errorf("AcidCommit() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompressFileType_Delete(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        key string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {
            name: "delete",
            fields: fields{
                path: "",
                data: nil,
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
                t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompressFileType_DeletebyPrefix(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        prefix string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {
            name: "deletebyprefix",
            fields: fields{
                path: "",
                data: nil,
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            if err := c.DeletebyPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
                t.Errorf("DeletebyPrefix() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompressFileType_Get(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        key string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    []byte
        wantErr bool
    }{
        {
            name: "get normal",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json": []byte("test"),
                },
            },
            args: args{
                key: "0.0.1/infrastructure.json",
            },
            want:    []byte("test"),
            wantErr: false,
        },
        {
            name: "get nil",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json": []byte("test"),
                },
            },
            args: args{
                key: "0.0.2/infrastructure.json",
            },
            want:    nil,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            got, err := c.Get(tt.args.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Get() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestCompressFileType_GetSourceDataorOperator(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    tests := []struct {
        name   string
        fields fields
        want   interface{}
    }{
        {
            name: "get normal",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json": []byte("test"),
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            got := c.GetSourceDataorOperator()
            if resmap, ok := got.(map[string][]byte); ok {
                if string(resmap["0.0.1/infrastructure.json"]) == "test" {
                    return
                }
            }
        })
    }
}

func TestCompressFileType_GetbyPrefix(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        prefix string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        want    map[string][]byte
        wantErr bool
    }{
        {
            name: "get normal",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json":    []byte("test1"),
                    "0.0.1/00/servicelist.json":    []byte("test2"),
                    "0.0.1/00/templates/temp.toml": []byte("test3"),
                },
            },
            args: args{
                prefix: "0.0.1/00/",
            },
            want: map[string][]byte{
                "0.0.1/00/servicelist.json":    []byte("test2"),
                "0.0.1/00/templates/temp.toml": []byte("test3"),
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            got, err := c.GetbyPrefix(tt.args.prefix)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetbyPrefix() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetbyPrefix() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestCompressFileType_Put(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        key   string
        value string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        {
            name: "put",
            fields: fields{
                path: "",
                data: nil,
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            if err := c.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
                t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompressFileType_getData(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        {
            name: "get data",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json":    []byte("test1"),
                    "0.0.1/00/servicelist.json":    []byte("test2"),
                    "0.0.1/00/templates/temp.toml": []byte("test3"),
                },
            },
            wantErr: true,
        },
        {
            name: "get data",
            fields: fields{
                path: "",
                data: map[string][]byte{
                    "0.0.1/infrastructure.json":    []byte("test1"),
                    "0.0.1/00/servicelist.json":    []byte("test2"),
                    "0.0.1/00/templates/temp.toml": []byte("test3"),
                },
            },
            wantErr: false,
        },
    }
    outputsDeCompress := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, errors.New("get data err")}},
        {Values: gomonkey.Params{map[string][]byte{
            "0.0.1/infrastructure.json":    []byte("test1"),
            "0.0.1/00/servicelist.json":    []byte("test2"),
            "0.0.1/00/templates/temp.toml": []byte("test3"),
        }, nil}},
    }
    patchesDeCompress := gomonkey.ApplyFuncSeq(util.DecompressFromPath, outputsDeCompress)
    defer patchesDeCompress.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            if err := c.getData(); (err != nil) != tt.wantErr {
                t.Errorf("getData() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestCompressFileType_setPath(t *testing.T) {
    type fields struct {
        path string
        data map[string][]byte
    }
    type args struct {
        path string
    }
    tests := []struct {
        name   string
        fields fields
        args   args
    }{
        {
            name: "setpath",
            fields: fields{
                path: "",
                data: nil,
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := &CompressFileType{
                path: tt.fields.path,
                data: tt.fields.data,
            }
            c.setPath(tt.fields.path)
        })
    }
}

func TestNewCompressedFileType(t *testing.T) {
    type args struct {
        path string
    }
    tests := []struct {
        name    string
        args    args
        want    *CompressFileType
        wantErr bool
    }{
        {
            name: "get err",
            args: args{
                path: "",
            },
            wantErr: true,
        },
        {
            name: "normal",
            args: args{
                path: "",
            },
            wantErr: false,
        },
    }
    outputsDeCompress := []gomonkey.OutputCell{
        {Values: gomonkey.Params{nil, errors.New("get data err")}},
        {Values: gomonkey.Params{map[string][]byte{
            "0.0.1/infrastructure.json":    []byte("test1"),
            "0.0.1/00/servicelist.json":    []byte("test2"),
            "0.0.1/00/templates/temp.toml": []byte("test3"),
        }, nil}},
    }
    patchesDeCompress := gomonkey.ApplyFuncSeq(util.DecompressFromPath, outputsDeCompress)
    defer patchesDeCompress.Reset()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewCompressedFileType(tt.args.path)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewCompressedFileType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            //if !reflect.DeepEqual(got, tt.want) {
            //	t.Errorf("NewCompressedFileType() got = %v, want %v", got, tt.want)
            //}
        })
    }
}
