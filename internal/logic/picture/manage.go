package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// Edit 编辑图片
func (s *sPicture) Edit(ctx context.Context, req *v1.PictureEditReq) (res *v1.PictureEditRes, err error) {
	// 1. 验证图片是否存在
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).
		Where(pic.IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 2. 验证用户权限（可选，根据业务需求）
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 只有图片创建者或管理员可以编辑
	if picture.UserId != user.Id && user.UserRole != consts.Admin {
		return nil, gerror.New("无权限编辑此图片")
	}

	// 3. 准备更新数据
	updateData := do.Picture{
		EditTime:     gtime.Now(),
		UpdateTime:   gtime.Now(),
		ReviewStatus: consts.DefRwStatus,
	}

	// 只更新非空字段
	if req.Name != "" {
		updateData.Name = req.Name
	}
	if req.Introduction != "" {
		updateData.Introduction = req.Introduction
	}
	if req.Category != "" {
		updateData.Category = req.Category
	}
	if len(req.Tags) > 0 {
		tagsJson, _ := gjson.New(req.Tags).ToJson()
		updateData.Tags = string(tagsJson)
	}
	if req.SpaceId > 0 {
		updateData.SpaceId = req.SpaceId
	}

	// 4. 更新图片信息到数据库
	_, err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).Data(updateData).Update()
	if err != nil {
		return nil, gerror.New("更新图片失败")
	}

	return &v1.PictureEditRes{
		Success: true,
	}, nil
}

// Delete 删除图片
func (s *sPicture) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).
		Where(pic.IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 2. 验证用户权限
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 只有图片创建者或管理员可以删除
	if picture.UserId != user.Id && user.UserRole != consts.Admin {
		return nil, gerror.New("无权限删除此图片")
	}

	// 3. 软删除图片记录
	_, err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).Data(do.Picture{
		IsDelete:   1,
		UpdateTime: gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.New("删除图片失败")
	}

	// 4. 删除对象存储中的文件
	if _, err = service.Bucket().Delete(ctx, &v1.BucketDeleteReq{
		FileName: picture.Name,
	}); err != nil {
		return nil, gerror.New("删除图片失败")
	}

	return &v1.DeleteRes{
		Success: true,
	}, nil
}

// Update 更新图片
func (s *sPicture) Update(ctx context.Context, req *v1.PictureUpdateReq) (res *v1.PictureUpdateRes, err error) {
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).
		Where(pic.IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 2. 验证用户权限
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 只有图片创建者或管理员可以更新
	if picture.UserId != user.Id && user.UserRole != "admin" {
		return nil, gerror.New("无权限更新此图片")
	}

	// 3. 准备更新数据
	updateData := do.Picture{
		ReviewStatus: consts.DefRwStatus,
		UpdateTime:   gtime.Now(),
	}

	// 只更新非空字段
	if req.Name != "" {
		updateData.Name = req.Name
	}
	if req.Introduction != "" {
		updateData.Introduction = req.Introduction
	}
	if req.Category != "" {
		updateData.Category = req.Category
	}
	if len(req.Tags) > 0 {
		tagsJson, _ := gjson.New(req.Tags).ToJson()
		updateData.Tags = string(tagsJson)
	}

	// 4. 更新图片信息到数据库
	_, err = dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).Data(updateData).Update()
	if err != nil {
		return nil, gerror.New("更新图片失败")
	}

	return &v1.PictureUpdateRes{
		Success: true,
	}, nil
}
