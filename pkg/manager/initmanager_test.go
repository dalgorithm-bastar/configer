package manage

import (
	"reflect"
	"testing"
)

const grpcLocation = "../../config/grpc.json"

func TestGetGrpcInfo(t *testing.T) {
	NewManager(grpcLocation)
	tests := []struct {
		name string
		want *GrpcInfoStruct
	}{
		{
			name: "get grpcInfo",
			want: &manager.grpcInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGrpcInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGrpcInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManager(t *testing.T) {
	NewManager(grpcLocation)
	tests := []struct {
		name string
		want *Manager
	}{
		{
			name: "get manager test",
			want: manager,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		grpcConfigLocation string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "correct load",
			args: args{
				grpcConfigLocation: grpcLocation,
			},
			wantErr: false,
		},
		{
			name: "incorrect path load",
			args: args{
				grpcConfigLocation: "",
			},
			wantErr: true,
		},
		{
			name: "incorrect file load",
			args: args{
				grpcConfigLocation: "delreqhdlr.go",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewManager(tt.args.grpcConfigLocation); (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
