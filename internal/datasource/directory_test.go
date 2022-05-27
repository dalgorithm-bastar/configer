package datasource

import (
	"context"
	"reflect"
	"testing"
)

func TestDirType_AcidCommit(t *testing.T) {
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
			name: "err",
			args: args{
				putMap:    nil,
				deleteMap: nil,
			},
			fields: fields{
				path: "",
				data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if err := c.AcidCommit(tt.args.putMap, tt.args.deleteMap); (err != nil) != tt.wantErr {
				t.Errorf("AcidCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirType_Delete(t *testing.T) {
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
			name: "err",
			args: args{
				key: "",
			},
			fields: fields{
				path: "",
				data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if err := c.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirType_DeletebyPrefix(t *testing.T) {
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
			name: "err",
			args: args{
				prefix: "",
			},
			fields: fields{
				path: "",
				data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if err := c.DeletebyPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("DeletebyPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirType_Get(t *testing.T) {
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
			name: "ok",
			fields: fields{
				path: "",
				data: map[string][]byte{
					"test": []byte("test"),
				},
			},
			args:    args{key: "test"},
			wantErr: false,
			want:    []byte("test"),
		},
		{
			name: "err",
			fields: fields{
				path: "",
				data: nil,
			},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
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

func TestDirType_GetSourceDataorOperator(t *testing.T) {
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
			name: "ok",
			fields: fields{
				path: "",
				data: map[string][]byte{"test": []byte("test")},
			},
			want: map[string][]byte{"test": []byte("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if got := c.GetSourceDataorOperator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSourceDataorOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirType_GetbyPrefix(t *testing.T) {
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
			name: "ok",
			fields: fields{
				path: "",
				data: map[string][]byte{"test": []byte("test")},
			},
			args:    args{prefix: "te"},
			want:    map[string][]byte{"test": []byte("test")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
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

func TestDirType_GracefullyClose(t *testing.T) {
	type fields struct {
		path string
		data map[string][]byte
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
				path: "",
				data: nil,
			},
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			c.GracefullyClose(tt.args.ctx)
		})
	}
}

func TestDirType_Put(t *testing.T) {
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
			name: "err",
			fields: fields{
				path: "",
				data: nil,
			},
			args: args{
				key:   "",
				value: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if err := c.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirType_getData(t *testing.T) {
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
			name: "ok",
			fields: fields{
				path: "../../unittestfiles/config",
				data: nil,
			},
			wantErr: false,
		},
		{
			name: "err",
			fields: fields{
				path: "../../unittesetfiles/configcenter.json",
				data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			if err := c.getData(); (err != nil) != tt.wantErr {
				t.Errorf("getData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirType_setPath(t *testing.T) {
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
			name: "ok",
			fields: fields{
				path: "",
				data: nil,
			},
			args: args{path: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &DirType{
				path: tt.fields.path,
				data: tt.fields.data,
			}
			c.setPath(tt.args.path)
		})
	}
}

func TestNewDirType(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *DirType
		wantErr bool
	}{
		{
			name: "ok",
			args: args{path: "../../unittestfiles/config"},
			want: &DirType{
				path: "../../unittestfiles/config",
				data: map[string][]byte{"config/configcenter.json": []byte(`{
  "etcd": {
    "endpoints": [
      "127.0.0.1:2379"
    ],
    "username": "root",
    "password": "chbw0011",
    "operationtimeout": 1
  },
  "grpc": {
    "socket": "127.0.0.1:2333",
    "locktimeout": 30
  },
  "log": {
    "logpath": "log/",
    "recordlevel": "info",
    "encodingtype": "normal",
    "filename": "configcenter",
    "maxage": 30
  }
}`)},
			},
			wantErr: false,
		},
		{
			name:    "err",
			args:    args{path: "../../unittestfiles/config/configcenter.json"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDirType(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDirType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDirType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
