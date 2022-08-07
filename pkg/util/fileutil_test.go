package util

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	patchesRead := gomonkey.ApplyFuncSeq(ioutil.ReadAll, outputsRead)
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
	fileMap, err := DecompressFromPath("../../test/unittestfiles/pkgs/0.0.1.tar.gz")
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

/*func TestWalkFromPath(t *testing.T) {
    type args struct {
        inputPath string
    }
    tests := []struct {
        name    string
        args    args
        want    map[string][]byte
        wantErr bool
    }{
        {
            name: "walk ok",
            args: args{
                inputPath: filepath.FromSlash("../../test/unittestfiles/config"),
            },
            wantErr: false,
            want: map[string][]byte{
                "config/configcenter.json": []byte(`{
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
}`),
            },
        },
        {
            name:    "read err",
            args:    args{inputPath: "../../test/unittestfiles/config/configcenter.json"},
            wantErr: true,
            want:    map[string][]byte{},
        },
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
}*/

func TestCheckYaml(t *testing.T) {
	yml2 := []byte(`mem1: "test1"
mem2: 0`)
	yml1 := []byte(`mem1: "test1"
mem2": 1`)
	yml3 := []byte(`mem1: "test1"
mem2: 0
mem3: "test"`)
	type args struct {
		srcYml []byte
		yml2   []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key diff",
			args: args{
				srcYml: yml1,
				yml2:   yml2,
			},
			want: false,
		},
		{
			name: "key diff reverse",
			args: args{
				srcYml: yml2,
				yml2:   yml1,
			},
			want: false,
		},
		{
			name: "key lack",
			args: args{
				srcYml: yml1,
				yml2:   yml3,
			},
			want: false,
		},
		{
			name: "key lack reverse",
			args: args{
				srcYml: yml3,
				yml2:   yml1,
			},
			want: false,
		},
		{
			name: "ok",
			args: args{
				srcYml: yml1,
				yml2:   yml1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckYaml(tt.args.srcYml, tt.args.yml2); got != tt.want {
				t.Errorf("CheckYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadDirWithPermFile(t *testing.T) {
	var outInfo bytes.Buffer
	cmd := exec.Command("chmod", "644", "../../test/unittestfiles/pkgs/permission/1.0.0/scheme1/DTP/MC/template/tmpl1/t1.toml")
	cmd.Stdout = &outInfo
	cmd.Stderr = &outInfo
	if err := cmd.Run(); err != nil {
		fmt.Println(outInfo.String())
		fmt.Println(err)
		os.Exit(1)
	}
	cmd = exec.Command("chmod", "755", "../../test/unittestfiles/pkgs/permission/1.0.0/scheme1/DTP/MC/template/tmpl1")
	if err := cmd.Run(); err != nil {
		fmt.Println(outInfo.String())
		fmt.Println(err)
		os.Exit(1)
	}
	type args struct {
		dirPath    string
		separator  string
		keyForTmpl string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]byte
		want1   []byte
		wantErr bool
	}{
		{
			name: "nromal",
			args: args{
				dirPath:    "../../test/unittestfiles/pkgs/permission/1.0.0",
				separator:  "/",
				keyForTmpl: "template",
			},
			want: map[string][]byte{
				"1.0.0/scheme1/DTP/MC/template/tmpl1/t1.toml": []byte(`platName = "{{.PlatName}}"`),
				"1.0.0/perm.yaml": []byte(`filePerms:
    - path: 1.0.0/scheme1/DTP/MC/template/tmpl1
      isDir: "1"
      perm: "755"
    - path: 1.0.0/scheme1/DTP/MC/template/tmpl1/t1.toml
      isDir: "0"
      perm: "644"
`)},
			want1: []byte(`filePerms:
    - path: 1.0.0/scheme1/DTP/MC/template/tmpl1
      isDir: "1"
      perm: "755"
    - path: 1.0.0/scheme1/DTP/MC/template/tmpl1/t1.toml
      isDir: "0"
      perm: "644"
`),
			wantErr: false,
		},
		/*		{
				name: "file perm 244",
				args: args{
					dirPath:    "../../test/unittestfiles/pkgs/permission/1.0.1",
					separator:  "/",
					keyForTmpl: "template",
				},
				want:    nil,
				want1:   nil,
				wantErr: true,
			},*/
		{
			name: "nil path",
			args: args{
				dirPath:    "",
				separator:  "/",
				keyForTmpl: "template",
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := LoadDirWithPermFile(tt.args.dirPath, tt.args.separator, tt.args.keyForTmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDirWithPermFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadDirWithPermFile() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("LoadDirWithPermFile() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
