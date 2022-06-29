package datasource

import (
    "context"
    "errors"
    "fmt"
    "os"
    "reflect"
    "testing"
    "time"

    "github.com/configcenter/internal/mock"
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/mvcc/mvccpb"
    "github.com/golang/mock/gomock"
    "github.com/spf13/viper"
)

func loadConfig() {
    viper.SetConfigFile("../../unittestfiles/config/configcenter.json")
    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}

func TestEtcdType_AtomicCommit(t *testing.T) {
    client3 := &clientv3.Client{}

    /*ctrl1 := gomock.NewController(t)
      defer ctrl1.Finish()
      mockTxn := mock.NewMockTxn(ctrl1)
      gomock.InOrder(
          mockTxn.EXPECT().Then(gomock.Any()).Return(mockTxn),
          mockTxn.EXPECT().Commit().Return(nil, nil),
      )

      ctrl := gomock.NewController(t)
      defer ctrl.Finish()
      mockKv := mock.NewMockKV(ctrl)
      gomock.InOrder(
          mockKv.EXPECT().Txn(gomock.Any()).Return(mockTxn),
      )*/

    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
    }
    type args struct {
        putMap      map[string]string
        deleteSlice []string
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr bool
    }{
        /*{
            name: "normal",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args: args{
                putMap:      map[string]string{"test": "test"},
                deleteSlice: []string{"test1"},
            },
            wantErr: false,
        },*/
        {
            name: "nil",
            fields: fields{
                client: client3,
                kv:     nil,
            },
            args: args{
                putMap:      nil,
                deleteSlice: nil,
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if err := e.AtomicCommit(tt.args.putMap, tt.args.deleteSlice); (err != nil) != tt.wantErr {
                t.Errorf("AcidCommit() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestEtcdType_ConnectToEtcd(t *testing.T) {
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
    }
    tests := []struct {
        name    string
        fields  fields
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if err := e.ConnectToEtcd(); (err != nil) != tt.wantErr {
                t.Errorf("ConnectToEtcd() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestEtcdType_Delete(t *testing.T) {
    client3 := &clientv3.Client{}

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockKv := mock.NewMockKV(ctrl)
    gomock.InOrder(
        mockKv.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil, errors.New("connect err")),
    )

    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
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
            name: "err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if err := e.Delete(tt.args.key); (err != nil) != tt.wantErr {
                t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestEtcdType_DeletebyPrefix(t *testing.T) {
    client3 := &clientv3.Client{}

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockKv := mock.NewMockKV(ctrl)
    gomock.InOrder(
        mockKv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("connect err")),
    )

    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
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
            name: "conn err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args:    args{prefix: ""},
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if err := e.DeletebyPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
                t.Errorf("DeletebyPrefix() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestEtcdType_Get(t *testing.T) {
    client3 := &clientv3.Client{}

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockKv := mock.NewMockKV(ctrl)
    gomock.InOrder(
        mockKv.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("connect err")),
        mockKv.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&clientv3.GetResponse{
            Kvs: []*mvccpb.KeyValue{{Value: []byte("test")}},
        }, nil),
    )
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
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
            name: "conn err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args:    args{key: ""},
            wantErr: true,
            want:    nil,
        },
        {
            name: "ok",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args:    args{key: ""},
            want:    []byte("test"),
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            got, err := e.Get(tt.args.key)
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

func TestEtcdType_GetSourceDataorOperator(t *testing.T) {
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
    }
    tests := []struct {
        name   string
        fields fields
        want   interface{}
    }{
        {
            name: "ok",
            fields: fields{
                client: &clientv3.Client{},
                kv:     nil,
            },
            want: &clientv3.Client{},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if got := e.GetSourceDataorOperator(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetSourceDataorOperator() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestEtcdType_GetbyPrefix(t *testing.T) {
    client3 := &clientv3.Client{}

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockKv := mock.NewMockKV(ctrl)
    gomock.InOrder(
        mockKv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("connect err")),
        mockKv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(&clientv3.GetResponse{
            Kvs: []*mvccpb.KeyValue{{Key: []byte("test"), Value: []byte("test")}},
        }, nil),
    )
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
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
            name: "conn err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args:    args{prefix: ""},
            wantErr: true,
            want:    nil,
        },
        {
            name: "ok",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args:    args{prefix: ""},
            wantErr: false,
            want:    map[string][]byte{"test": []byte("test")},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            got, err := e.GetbyPrefix(tt.args.prefix)
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

func TestEtcdType_GracefullyClose(t *testing.T) {
    ctx3 := context.Background()
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
    }
    type args struct {
        ctx context.Context
    }
    tests := []struct {
        name   string
        fields fields
        args   args
    }{
        {
            name: "ok",
            fields: fields{
                client: nil,
                kv:     nil,
            },
            args: args{ctx: ctx3},
        },
        {
            name: "ok2",
            fields: fields{
                client: &clientv3.Client{
                    Cluster:     nil,
                    KV:          nil,
                    Lease:       nil,
                    Watcher:     nil,
                    Auth:        nil,
                    Maintenance: nil,
                    Username:    "",
                    Password:    "",
                },
                kv: nil,
            },
            args: args{ctx: ctx3},
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            go e.GracefullyClose(tt.args.ctx)
            time.Sleep(100 * time.Microsecond)
            ctx3.Done()
        })
    }
}

func TestEtcdType_Put(t *testing.T) {
    client3 := &clientv3.Client{}

    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockKv := mock.NewMockKV(ctrl)
    gomock.InOrder(
        mockKv.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("connect err")),
        mockKv.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),
    )
    type fields struct {
        client *clientv3.Client
        kv     clientv3.KV
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
            name: "err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args: args{
                key:   "",
                value: "",
            },
            wantErr: true,
        },
        {
            name: "err",
            fields: fields{
                client: client3,
                kv:     mockKv,
            },
            args: args{
                key:   "",
                value: "",
            },
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            e := &EtcdType{
                client: tt.fields.client,
                kv:     tt.fields.kv,
            }
            if err := e.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
                t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestNewEtcdType(t *testing.T) {
    tests := []struct {
        name    string
        want    *EtcdType
        wantErr bool
    }{
        {
            name:    "err",
            wantErr: true,
            want:    nil,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewEtcdType()
            if (err != nil) != tt.wantErr {
                t.Errorf("NewEtcdType() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewEtcdType() got = %v, want %v", got, tt.want)
            }
        })
    }
}
