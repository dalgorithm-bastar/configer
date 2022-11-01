package service

import (
	"context"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/pb"
	"github.com/configcenter/pkg/repository"
	"go.uber.org/zap"
)

//DeleteInManager 接口用于删除用户名下所有缓存数据
func DeleteInManager(ctx context.Context, req *pb.CfgReq) (error, []*pb.VersionInfo, *pb.AnyFile) {
	err := repository.Src.DeletebyPrefix(req.UserName)
	log.Logger.Info("call delete func delete UserName:%s, delete result %v", zap.String("delete UserName", req.UserName), zap.Any("delete result", err))
	return err, nil, nil
}
