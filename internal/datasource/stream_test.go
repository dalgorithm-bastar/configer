package datasource

import (
	"reflect"
	"testing"
)

func TestNewStream(t *testing.T) {
	type args struct {
		fileMap map[string][]byte
	}
	tests := []struct {
		name string
		args args
		want *Stream
	}{
		{
			name: "new ins",
			args: args{
				fileMap: map[string][]byte{
					"0.0.1/infrastructure.json": []byte("test1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewStream(tt.args.fileMap)
			//if got := NewStream(tt.args.fileMap); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewStream() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestStream_AcidCommit(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
				&CompressFileType{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stream{
				c: tt.fields.c,
			}
			if err := s.AcidCommit(tt.args.putMap, tt.args.deleteMap); (err != nil) != tt.wantErr {
				t.Errorf("AcidCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStream_Delete(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
				&CompressFileType{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stream{
				c: tt.fields.c,
			}
			if err := s.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStream_DeletebyPrefix(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
				&CompressFileType{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stream{
				c: tt.fields.c,
			}
			if err := s.DeletebyPrefix(tt.args.prefix); (err != nil) != tt.wantErr {
				t.Errorf("DeletebyPrefix() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStream_Get(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
			name: "get",
			fields: fields{
				c: &CompressFileType{
					path: "",
					data: map[string][]byte{
						"0.0.1/infrastructure.json":    []byte("test1"),
						"0.0.1/00/servicelist.json":    []byte("test2"),
						"0.0.1/00/templates/temp.toml": []byte("test3"),
					},
				},
			},
			args: args{
				key: "0.0.1/infrastructure.json",
			},
			want:    []byte("test1"),
			wantErr: false,
		},
		{
			name: "get nil",
			fields: fields{
				c: &CompressFileType{
					path: "",
					data: map[string][]byte{
						"0.0.1/infrastructure.json":    []byte("test1"),
						"0.0.1/00/servicelist.json":    []byte("test2"),
						"0.0.1/00/templates/temp.toml": []byte("test3"),
					},
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
			s := &Stream{
				c: tt.fields.c,
			}
			got, err := s.Get(tt.args.key)
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

func TestStream_GetSourceDataorOperator(t *testing.T) {
	type fields struct {
		c *CompressFileType
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "get normal",
			fields: fields{
				c: &CompressFileType{
					path: "",
					data: map[string][]byte{
						"0.0.1/infrastructure.json": []byte("test"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stream{
				c: tt.fields.c,
			}
			got := s.GetSourceDataorOperator()
			if resmap, ok := got.(map[string][]byte); ok {
				if string(resmap["0.0.1/infrastructure.json"]) == "test" {
					return
				}
			}
			//if got := s.GetSourceDataorOperator(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetSourceDataorOperator() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestStream_GetbyPrefix(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
			name: "get prefix",
			fields: fields{
				c: &CompressFileType{
					path: "",
					data: map[string][]byte{
						"0.0.1/infrastructure.json":    []byte("test1"),
						"0.0.1/00/servicelist.json":    []byte("test2"),
						"0.0.1/00/templates/temp.toml": []byte("test3"),
					},
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
			s := &Stream{
				c: tt.fields.c,
			}
			got, err := s.GetbyPrefix(tt.args.prefix)
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

func TestStream_Put(t *testing.T) {
	type fields struct {
		c *CompressFileType
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
				&CompressFileType{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stream{
				c: tt.fields.c,
			}
			if err := s.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
