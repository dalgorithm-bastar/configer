package datasource

import (
	"go.etcd.io/etcd/clientv3"
	"reflect"
	"testing"
)

const etcdLocation = "../../config/etcdClientv3_test.json"

func TestEtcdType_AcidCommit(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
		{
			name: "err",
			fields: fields{
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				putMap: map[string]string{
					"test": "test",
				},
				deleteSlice: []string{"test"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if err := e.AcidCommit(tt.args.putMap, tt.args.deleteSlice); (err != nil) != tt.wantErr {
				t.Errorf("AcidCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcdType_ConnectToEtcd(t *testing.T) {
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
	}
	type args struct {
		etcdConfigLocation string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "normal",
			fields:  fields{},
			args:    args{etcdConfigLocation: etcdLocation},
			wantErr: false,
		},
		{
			name:    "open err",
			fields:  fields{},
			args:    args{etcdConfigLocation: ""},
			wantErr: true,
		},
		{
			name:    "file err",
			fields:  fields{},
			args:    args{etcdConfigLocation: "stream.go"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if err := e.ConnectToEtcd(tt.args.etcdConfigLocation); (err != nil) != tt.wantErr {
				t.Errorf("ConnectToEtcd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcdType_Delete(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				key: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if err := e.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcdType_DeletebyPrefix(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
			name: "err",
			fields: fields{
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				prefix: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if err := e.DeletebyPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("DeletebyPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEtcdType_Get(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
			name: "err",
			fields: fields{
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				key: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
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
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "nil",
			fields: fields{
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			want: etcdType.client,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if got := e.GetSourceDataorOperator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSourceDataorOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEtcdType_GetbyPrefix(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
			name: "err",
			fields: fields{
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				prefix: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
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

func TestEtcdType_Put(t *testing.T) {
	etcdType, err := NewEtcdType(etcdLocation)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		client *clientv3.Client
		kv     clientv3.KV
		cfgMap *EtcdInfo
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
				client: etcdType.client,
				kv:     etcdType.kv,
				cfgMap: etcdType.cfgMap,
			},
			args: args{
				key:   "test",
				value: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EtcdType{
				client: tt.fields.client,
				kv:     tt.fields.kv,
				cfgMap: tt.fields.cfgMap,
			}
			if err := e.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEtcdType(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *EtcdType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEtcdType(tt.args.configPath)
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
