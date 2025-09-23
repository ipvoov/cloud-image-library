package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// Review 审核图片
func (s *sPicture) Review(ctx context.Context, req *v1.PictureReviewReq) (res *v1.PictureReviewRes, err error) {
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).
		Where(pic.IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 2. 验证用户权限 - 只有管理员可以审核
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	if user.UserRole != "admin" {
		return nil, gerror.New("无权限审核图片，只有管理员可以审核")
	}

	// 4. 更新审核信息
	updateData := do.Picture{
		ReviewStatus:  req.ReviewStatus,
		ReviewMessage: req.ReviewMessage,
		ReviewerId:    user.Id,
		ReviewTime:    gtime.Now(),
		UpdateTime:    gtime.Now(),
	}

	_, err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).Data(updateData).Update()
	if err != nil {
		return nil, gerror.New("审核操作失败")
	}

	return &v1.PictureReviewRes{
		Success: true,
	}, nil
}
