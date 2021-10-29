package manage

import (
    "context"

    "github.com/configcenter/pkg/pb"
    "github.com/configcenter/pkg/repository"
)

//delete接口用于删除用户名下所有缓存数据
func deleteInManager(ctx context.Context, req *pb.CfgReq) (error, []string, string, []byte) {
    err := repository.Src.DeletebyPrefix(req.UserName)
    //log.Sugar().Infof("call delete func, delete UserName:%s, delete result %v", req.UserName, err)
    return err, nil, "", nil
}
