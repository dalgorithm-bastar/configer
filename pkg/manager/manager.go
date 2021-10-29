//使用Manager类实现service各个接口的单独文件
package manage

import (
    "context"

    "github.com/configcenter/internal/log"
    "github.com/configcenter/pkg/pb"
)

func (m *Manager) GET(ctx context.Context, CfgReq *pb.CfgReq) (*pb.CfgResp, error) {
    go log.Sugar().Infof("GET CfgReq Recieved: %+v", CfgReq)
    if CfgReq == nil {
        return &pb.CfgResp{Status: "nil req deliverd"}, nil
    }
    err, sliceData, fileName, fileData := Get(m.ctx, CfgReq)
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
    err, sliceData, fileName, fileData := post(m.ctx, CfgReq)
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
    err, sliceData, fileName, fileData := deleteInManager(m.ctx, CfgReq)
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
    err, sliceData, fileName, fileData := put(m.ctx, CfgReq)
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
