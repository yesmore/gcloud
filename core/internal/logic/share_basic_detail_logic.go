package logic

import (
	"context"

	"gcloud/core/internal/svc"
	"gcloud/core/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShareBasicDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareBasicDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareBasicDetailLogic {
	return &ShareBasicDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareBasicDetailLogic) ShareBasicDetail(req *types.ShareBasicDetailRequest) (resp *types.ShareBasicDetailReply, err error) {
	resp = new(types.ShareBasicDetailReply)
	// 1 更新分享记录的点击次数
	err = l.svcCtx.Engine.
		Table("share_basic").
		Where("identity = ?", req.Identity).
		Exec("UPDATE share_basic SET click_num = click_num + 1 where identity = ?", req.Identity).Error
	if err != nil {
		resp.Msg = "error"
		return
	}

	// 2 获取资源详细信息
	err = l.svcCtx.Engine.
		Table("share_basic").
		Select("share_basic.repository_identity, user_repository.name, repository_pool.ext, repository_pool.path, repository_pool.size, share_basic.click_num").
		Joins("LEFT JOIN repository_pool ON repository_pool.identity = share_basic.repository_identity").
		Joins("LEFT JOIN user_repository ON user_repository.identity = share_basic.user_repository_identity").
		Where("share_basic.identity = ?", req.Identity).
		First(resp).Error
	if err != nil {
		resp.Msg = "error"
		return
	}

	return
}
