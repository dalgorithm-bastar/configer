package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/configcenter/internal/mock"
	"github.com/configcenter/pkg/repository"
	"github.com/golang/mock/gomock"
	"testing"
)

var (
	jsonString = `{
  "test1": "1",
  "test2": ["1","2","3"],
  "test31": {
    "test32": "4"
  },
  "test4":[{
	  "test41":"test41"
	}
  ],
  "test5":1
}`
	infrastructureJson = `{
  "normal": {
    "hostname1": {
      "oplan": {
        "ip": "10.0.1.1"
      },
      "biznet": {
        "ip": "192.168.0.1"
      }
    }
  },
  "hot_backup": {
    "hostname4": {
      "oplan": {
        "ip": "10.0.1.4"
      },
      "biznet": {
        "ip": "192.168.0.4"
      }
    }
  }
}
`
	serviceListJson = `{
  "replicator_number": "3",
  "deployment_info": [
    {
      "hostname": "hostname1",
      "MUDP_BIND_IP": "{biznet.ip}"
    }
  ],
  "MUDP_IP": "224.1.0.4",
  "MUDP_NORMAL_PORT": "9999",
  "MUDP_SUPPLEMNET_PORT": "0",
  "MUDP_NOTIFY_PORT": "0",
  "kafka_broker": [
    "179.7.89.3:29092",
    "179.7.89.7:29092",
    "179.7.89.11:29092",
    "179.7.219.4:29092",
    "179.7.219.6:29092"
  ]
}`
	deploymentString1 = `{"replicator_number":"3","deployment_info":[{"hostname":"hostname1","MUDP_BIND_IP":"192.168.0.1"}]}`
	deploymentString2 = `{"replicator_number":"3","deployment_info":[{"MUDP_BIND_IP":"192.168.0.1","hostname":"hostname1"}]}`
)

func TestConstructMap(t *testing.T) {
	var interData interface{}
	err := json.Unmarshal([]byte(jsonString), &interData)
	if err != nil {
		t.Fatal(err)
	}
	targetMap := map[string]string{
		"test1":          "1",
		"test2":          "1,2,3",
		"test31.test32":  "4",
		"test4.0.test41": "test41",
		"test5":          "1",
	}
	if err != nil {
		t.Fatal(err)
	}
	resMapTest := make(map[string]string)
	type args struct {
		resMap      map[string]string
		data        interface{}
		currentPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				resMap:      resMapTest,
				currentPath: "",
				data:        interData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConstructMap(tt.args.resMap, tt.args.data, tt.args.currentPath)
			for targetKey, targetv := range targetMap {
				if _, ok := tt.args.resMap[targetKey]; !ok {
					t.Fatal(fmt.Sprintf("cannot get key %s in resmap, resmap of %v", targetKey, tt.args.resMap))
				}
				if targetv != tt.args.resMap[targetKey] {
					t.Fatal(fmt.Sprintf("different value %s in resmap, that of %s in target", tt.args.resMap[targetKey], targetv))
				}
			}
		})
	}
}

func TestCtlFind(t *testing.T) {
	type args struct {
		tar           string
		ver           string
		env           string
		clusterObject string
		service       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal meet",
			args: args{
				tar:           Services,
				ver:           "0.0.1",
				env:           "00",
				clusterObject: "DTP.MC.set0",
				service:       "test1",
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "normal not meet",
			args: args{
				tar:           Services,
				ver:           "0.0.1",
				env:           "00",
				clusterObject: "DTP.MC.set0",
				service:       "testx",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "get err",
			args: args{
				tar:           Services,
				ver:           "0.0.1",
				env:           "00",
				clusterObject: "DTP.MC.set0",
				service:       "testx",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "get nil",
			args: args{
				tar:           Infrastructure,
				ver:           "0.0.1",
				env:           "00",
				clusterObject: "DTP.MC.set0",
				service:       "testx",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "json err",
			args: args{
				tar:           Infrastructure,
				ver:           "0.0.1",
				env:           "00",
				clusterObject: "DTP.MC.set0",
				service:       "testx",
			},
			want:    "",
			wantErr: true,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(jsonString), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(jsonString), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, errors.New("get err")),
		mockSrc.EXPECT().Get(gomock.Any()).Return(nil, nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("json err"), nil),
	)
	repository.Src = mockSrc
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CtlFind(tt.args.tar, tt.args.ver, tt.args.env, tt.args.clusterObject, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("CtlFind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CtlFind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDeploymentInfo(t *testing.T) {
	type args struct {
		serviceData        []byte
		infrastructureData []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				infrastructureData: []byte(infrastructureJson),
				serviceData:        []byte(serviceListJson),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDeploymentInfo(tt.args.serviceData, tt.args.infrastructureData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeploymentInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != deploymentString1 && got != deploymentString2 {
				t.Errorf("GetDeploymentInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSrc := mock.NewMockStorage(ctrl)
	gomock.InOrder(
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(serviceListJson), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(serviceListJson), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(serviceListJson), nil),
		mockSrc.EXPECT().Get(gomock.Any()).Return([]byte(infrastructureJson), nil),
		//mockSrc.EXPECT().Get(gomock.Any()).Return([]byte("json err"), nil),
	)
	repository.Src = mockSrc
	type args struct {
		src                repository.Storage
		infrastructureData []byte
		defaultIndex       bool
		globalId           string
		localId            string
		ver                string
		env                string
		clusterObject      string
		service            string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "normal replace",
			args: args{
				src:                repository.Src,
				infrastructureData: []byte(infrastructureJson),
				defaultIndex:       true,
				globalId:           "3254",
				localId:            "0",
				ver:                "0.0.1",
				env:                "00",
				clusterObject:      "DTP.MC.set0",
				service:            "deployment_info.MUDP_BIND_IP",
			},
			want:    "192.168.0.1",
			wantErr: false,
		},
		{
			name: "normal not replace",
			args: args{
				src:                repository.Src,
				infrastructureData: []byte(infrastructureJson),
				defaultIndex:       true,
				globalId:           "3254",
				localId:            "0",
				ver:                "0.0.1",
				env:                "00",
				clusterObject:      "DTP.MC.set0",
				service:            "deployment_info.hostname",
			},
			want:    "hostname1",
			wantErr: false,
		},
		{
			name: "normal no infrastructure file",
			args: args{
				src:                repository.Src,
				infrastructureData: nil,
				defaultIndex:       true,
				globalId:           "3254",
				localId:            "0",
				ver:                "0.0.1",
				env:                "00",
				clusterObject:      "DTP.MC.set0",
				service:            "deployment_info.MUDP_BIND_IP",
			},
			want:    "192.168.0.1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := baseGet(tt.args.src, tt.args.infrastructureData, tt.args.defaultIndex, tt.args.globalId, tt.args.localId, tt.args.ver, tt.args.env, tt.args.clusterObject, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("baseGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("baseGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}
