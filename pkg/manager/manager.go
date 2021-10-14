//Manager类的具体实现
package manage

import (
	"context"
	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/pb"
	"regexp"
)

//请求体Target字段取值范围
const (
	VersionString  = `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$` //初始化用于版本号校验的正则表达式
	TemplateString = `templates/`                                          //初始化用于模板文件筛选的正则表达式
)

var (
	manager *Manager
)

type regExpStruct struct {
	RegExpOfVersion  *regexp.Regexp
	RegExpOfTemplate *regexp.Regexp
}

//调度各模块实现请求
type Manager struct {
	grpcInfo GrpcInfoStruct
	regExp   regExpStruct
}

func (m *Manager) GET(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
	go log.Sugar().Infof("GET CfgReq Recieved: %+v", CfgReq)
	if CfgReq == nil {
		return &pb.CfgResp{Status: "nil req deliverd"}, nil
	}
	err, sliceData, fileName, fileData := Get(ctx, CfgReq)
	var status string
	if err != nil {
		status = err.Error()
	} else {
		status = "ok"
	}
	CfgResp := &pb.CfgResp{
		Status:    status,
		SliceData: sliceData,
		File: &pb.CompressedFile{
			FileName: fileName,
			FileData: fileData,
		},
	}
	log.Sugar().Infof("GET CfgResp Created: %+v", CfgResp)
	return CfgResp, nil
}

func (m *Manager) POST(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
	go log.Sugar().Infof("POST CfgReq Recieved: %+v", CfgReq)
	if CfgReq == nil {
		return &pb.CfgResp{Status: "nil req deliverd"}, nil
	}
	err, sliceData, fileName, fileData := post(ctx, CfgReq)
	var status string
	if err != nil {
		status = err.Error()
	} else {
		status = "ok"
	}
	CfgResp := &pb.CfgResp{
		Status:    status,
		SliceData: sliceData,
		File: &pb.CompressedFile{
			FileName: fileName,
			FileData: fileData,
		},
	}
	log.Sugar().Infof("POST CfgResp Created: %+v", CfgResp)
	return CfgResp, nil
}

func (m *Manager) DELETE(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
	go log.Sugar().Infof("DELETE CfgReq Recieved: %+v", CfgReq)
	if CfgReq == nil {
		return &pb.CfgResp{Status: "nil req deliverd"}, nil
	}
	err, sliceData, fileName, fileData := deleteInManager(ctx, CfgReq)
	var status string
	if err != nil {
		status = err.Error()
	} else {
		status = "ok"
	}
	CfgResp := &pb.CfgResp{
		Status:    status,
		SliceData: sliceData,
		File: &pb.CompressedFile{
			FileName: fileName,
			FileData: fileData,
		},
	}
	log.Sugar().Infof("DELETE CfgResp Created: %+v", CfgResp)
	return CfgResp, nil
}

func (m *Manager) PUT(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
	//go log.Sugar().Infof("PUT CfgReq Recieved: %+v", CfgReq)
	if CfgReq == nil {
		return &pb.CfgResp{Status: "nil req deliverd"}, nil
	}
	err, sliceData, fileName, fileData := put(ctx, CfgReq)
	var status string
	if err != nil {
		status = err.Error()
	} else {
		status = "ok"
	}
	CfgResp := &pb.CfgResp{
		Status:    status,
		SliceData: sliceData,
		File: &pb.CompressedFile{
			FileName: fileName,
			FileData: fileData,
		},
	}
	log.Sugar().Infof("PUT CfgResp Created: %+v", CfgResp)
	return CfgResp, nil
}
