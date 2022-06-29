package generation

import (
    "context"
    "os"
    "reflect"
    "testing"

    "github.com/configcenter/pkg/repository"
)

const (
    infraPath = "../../test/unittestfiles/pkgs/pkg3/infrastructure.yaml"
    cfgPath   = "../../test/unittestfiles/pkgs/pkg3/3.1.0"
)

func TestFinishResMap(t *testing.T) {
    type args struct {
        resMap         map[string][]byte
        dplyStructList []ChartDeployMain
        topicInfoList  map[string]map[string]map[string]ExpTpcMain
        prePath        string
    }
    tests := []struct {
        name    string
        args    args
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := FinishResMap(tt.args.resMap, tt.args.dplyStructList, tt.args.topicInfoList, tt.args.prePath); (err != nil) != tt.wantErr {
                t.Errorf("FinishResMap() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestGenerate(t *testing.T) {
    //读取数据
    infraData, _ := os.ReadFile(infraPath)
    repository.NewStorage(context.Background(), repository.DirType, cfgPath)
    rawData, _ := repository.Src.GetbyPrefix("3.1.0/scheme1")
    repository.NewStorage(context.Background(), repository.DirType,
        "../../test/unittestfiles/res/pkg3res/3.1.0_scheme1")
    resMap, _ := repository.Src.GetbyPrefix("")
    type args struct {
        infrastructure []byte
        rawData        map[string][]byte
        envNum         string
        topicIpRange   []string
        topicPortRange []string
        tcpPortRange   []string
    }
    tests := []struct {
        name    string
        args    args
        want    map[string][]byte
        wantErr bool
    }{
        {
            name: "normal gen",
            args: args{
                infrastructure: infraData,
                rawData:        rawData,
                envNum:         "01",
                topicIpRange:   []string{"156.10.1.1", "156.10.11.2"},
                topicPortRange: []string{"10000", "32768"},
                tcpPortRange:   []string{"10000", "32768"},
            },
            want:    resMap,
            wantErr: false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := Generate(tt.args.infrastructure, tt.args.rawData, tt.args.envNum, tt.args.topicIpRange, tt.args.topicPortRange, tt.args.tcpPortRange)
            if (err != nil) != tt.wantErr {
                t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            /*for k, v := range got {
                if _, ok := resMap[k]; !ok {
                    t.Errorf("lack of file %s in res", k)
                }
                if len(v) != len(resMap[k]) {
                    t.Errorf("differ from std res:%s", k)
                    fmt.Println(v)
                    fmt.Println(resMap[k])
                }
            }*/
        })
    }
}

func Test_addThirdPartFiles(t *testing.T) {
    type args struct {
        resMap         map[string][]byte
        infrastructure []byte
        dplyStructList []ChartDeployMain
        envNum         string
    }
    tests := []struct {
        name    string
        args    args
        wantErr bool
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := addThirdPartFiles(tt.args.resMap, tt.args.infrastructure, tt.args.dplyStructList, tt.args.envNum); (err != nil) != tt.wantErr {
                t.Errorf("addThirdPartFiles() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func Test_sortRawData(t *testing.T) {
    type args struct {
        rawData map[string][]byte
    }
    tests := []struct {
        name string
        args args
        want []RawFile
    }{
        // TODO: Add test cases.
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := sortRawData(tt.args.rawData); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("sortRawData() = %v, want %v", got, tt.want)
            }
        })
    }
}
