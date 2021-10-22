package util

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewStreamFile(t *testing.T) {
	type args struct {
		inputData []byte
		name      string
		size      int64
	}
	tests := []struct {
		name string
		args args
		want *StreamFile
	}{
		{
			name: "new ins",
			args: args{
				inputData: []byte("test"),
				name:      "testFile",
				size:      int64(len("test")),
			},
			want: &StreamFile{
				data: []byte("test"),
				name: "testFile",
				size: int64(len("test")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = NewStreamFile(tt.args.inputData, tt.args.name, tt.args.size)
			//if got := NewStreamFile(tt.args.inputData, tt.args.name, tt.args.size); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewStreamFile() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestStreamFile_IsDir(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// true
		{
			name: "true",
			fields: fields{
				data:    nil,
				name:    "testFile/",
				size:    0,
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: true,
		},
		// false
		{
			name: "true",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.IsDir(); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamFile_ModTime(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// test
		{
			name: "test",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
			},
			want: time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.ModTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamFile_Mode(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   os.FileMode
	}{
		{
			name: "modeTest",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: os.FileMode(777),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.Mode(); got != tt.want {
				t.Errorf("Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamFile_Name(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "filename",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: "testFile.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamFile_Size(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "size",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: int64(len("test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamFile_Sys(t *testing.T) {
	type fields struct {
		data    []byte
		name    string
		size    int64
		mode    os.FileMode
		modTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "sys",
			fields: fields{
				data:    []byte("test"),
				name:    "testFile.txt",
				size:    int64(len("test")),
				mode:    os.FileMode(777),
				modTime: time.Now(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamFile{
				data:    tt.fields.data,
				name:    tt.fields.name,
				size:    tt.fields.size,
				mode:    tt.fields.mode,
				modTime: tt.fields.modTime,
			}
			if got := s.Sys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sys() = %v, want %v", got, tt.want)
			}
		})
	}
}
