package util

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

var initialFile2 string = `[rmb_publisher]      #RMB发布端
TOPIC_TYPE = "1"                         #组播类型，可靠组播
HOST = {{GetInfo true "DTP.MC.set0" "deployment_info.hostname"}}
LOCAL_IP = {{GetInfo true "DTP.MC.set0" "deployment_info.MUDP_BIND_IP"}}  #实例化模板时选一个IP    #本机LOCAL_IP
KAFKA_BROKER=[{{GetInfo false "DTP.MC.set0" "kafka_broker"}}]
MUDP_IP = {{GetInfo false "DTP.MC.set0" "MUDP_IP"}}                    #组播IP
MUDP_NORMAL_PORT = 9999                  #组播端口
MUDP_SUPPLEMNET_PORT = 0
MUDP_NOTIFY_PORT = 0
FLOW_FILE_SIZE = 100
UNAIDED_IO = 0

[rmb_receiver]                         #同步器接收端
TOPIC_TYPE = "1"                         #组播类型，可靠组播
LOCAL_IP = {{GetInfo true "DTP.MC.set0" "deployment_info.MUDP_BIND_IP"}}                   #本机LOCAL_IP
MUDP_IP = {{GetInfo false "EzEI.set1" "MUDP_IP"}}                    #组播IP
MUDP_NORMAL_PORT = {{GetInfo false "EzEI.set1" "MUDP_NORMAL_PORT"}}                  #组播端口
MUDP_SUPPLEMNET_PORT = 0
MUDP_NOTIFY_PORT = 0
FLOW_FILE_SIZE = 100
UNAIDED_IO = 0
SUB_GROUP_ID = 0`

func TestCompressToStream(t *testing.T) {
	type args struct {
		stringWithFormat string
		fileMap          map[string][]byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// normal tar
		{
			name: "normal tar",
			args: args{
				stringWithFormat: "test.tar.gz",
				fileMap: map[string][]byte{
					"0.0.1/test1": []byte("test1"),
					"0.0.1/test2": []byte("test2"),
				},
			},
			wantErr: false,
		},
		//normal zip
		{
			name: "normal zip",
			args: args{
				stringWithFormat: "test.zip",
				fileMap: map[string][]byte{
					"0.0.1/test1": []byte("test1"),
					"0.0.1/test2": []byte("test2"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompressToStream(tt.args.stringWithFormat, tt.args.fileMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressToStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			format := "test.tar.gz"
			if tt.name == "normal zip" {
				format = "test.zip"
			}
			resmap, err := DecompressFromStream(format, got)
			if err != nil {
				t.Errorf("CompressToStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(string(resmap["0.0.1/test1"]) == "test1" && string(resmap["0.0.1/test2"]) == "test2") {
				t.Errorf("resmap err of %v", resmap)
				return
			}
		})
	}
}

func TestDecompressFromPath(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		// open err
		{
			name: "open err",
			args: args{
				inputPath: "",
			},
			want:    nil,
			wantErr: true,
		},
		// read err
		{
			name: "read err",
			args: args{
				inputPath: "streamfile.go",
			},
			want:    nil,
			wantErr: true,
		},
	}
	outputsOpen := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("open err")}},
		{Values: gomonkey.Params{&os.File{}, nil}},
	}
	patchesOpen := gomonkey.ApplyFuncSeq(os.Open, outputsOpen)
	defer patchesOpen.Reset()
	outputsRead := []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("read err")}},
		{Values: gomonkey.Params{nil, errors.New("read err")}},
	}
	patchesRead := gomonkey.ApplyFuncSeq(os.Open, outputsRead)
	defer patchesRead.Reset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecompressFromPath(tt.args.inputPath)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecompressFromPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecompressFromStream(t *testing.T) {
	fileMap, err := DecompressFromPath("../../compressfiledemo/0.0.1.tar.gz")
	if err != nil {
		t.Errorf("Decompress from path err:%v", err)
	}
	fileMap["0.0.1/99/"] = nil
	fileMap["0.0.1/99/servicelist.json"] = []byte(initialFile2)
	fileMap["0.0.1/99/template1.toml"] = []byte(initialFile2)
	data, err := CompressToStream("0.0.1.tar.gz", fileMap)
	fileMapPrecessed := make(map[string][]byte)
	for k, v := range fileMap {
		if k[len(k)-1] != '/' {
			fileMapPrecessed[k] = v
		}
	}
	if err != nil {
		t.Errorf("Decompress from path err:%v", err)
	}
	type args struct {
		stringWithFormat string
		binaryFile       []byte
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		// nil tar
		{
			name: "nil tar",
			args: args{
				stringWithFormat: "test.tar.gz",
				binaryFile:       nil,
			},
			wantErr: false,
		},
		//nil zip
		{
			name: "nil zip",
			args: args{
				stringWithFormat: "test.zip",
				binaryFile:       nil,
			},
			wantErr: false,
		},
		//
		{
			name: "insert stream to files test",
			args: args{
				stringWithFormat: "0.0.1.tar.gz",
				binaryFile:       data,
			},
			want:    fileMapPrecessed,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecompressFromStream(tt.args.stringWithFormat, tt.args.binaryFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecompressFromStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecompressFromStream() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWalkFromPath(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WalkFromPath(tt.args.inputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalkFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WalkFromPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
