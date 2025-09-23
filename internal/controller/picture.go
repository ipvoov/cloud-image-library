package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Picture = cPicture{}

type cPicture struct{}

// TagCategory 获取图片标签分类
func (c *cPicture) TagCategory(ctx context.Context, req *v1.PictureTagCategoryReq) (res *v1.PictureTagCategoryRes, err error) {
	return service.Picture().TagCategory(ctx, req)
}

// Edit 编辑图片
func (c *cPicture) Edit(ctx context.Context, req *v1.PictureEditReq) (res *v1.PictureEditRes, err error) {
	return service.Picture().Edit(ctx, req)
}

// Upload 上传图片
func (c *cPicture) Upload(ctx context.Context, req *v1.PictureUploadReq) (res *v1.PictureUploadRes, err error) {
	r := ghttp.RequestFromCtx(ctx)

	// 获取上传的文件
	file := r.GetUploadFile("file")
	if file == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "文件不能为空")
	}

	// 调用服务层处理
	return service.Picture().Upload(ctx, req, file)
}

// UploadByUrl 通过URL上传图片
func (c *cPicture) UploadByUrl(ctx context.Context, req *v1.PictureUploadByUrlReq) (res *v1.PictureUploadByUrlRes, err error) {
	return service.Picture().UploadByUrl(ctx, req)
}

// Delete 删除图片
func (c *cPicture) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	return service.Picture().Delete(ctx, req)
}

// Get 获取图片详情
func (c *cPicture) Get(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureAdminGetRes, err error) {
	return service.Picture().Get(ctx, req)
}

// GetVO 获取图片详情VO
func (c *cPicture) GetVO(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureGetRes, err error) {
	return service.Picture().GetVO(ctx, req)
}

// ListByPage 分页查询图片
func (c *cPicture) ListByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureAdminQueryRes, err error) {
	return service.Picture().ListByPage(ctx, req)
}

// ListVOByPage 分页查询图片VO
func (c *cPicture) ListVOByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureQueryRes, err error) {
	return service.Picture().ListVOByPage(ctx, req)
}

// Update 更新图片
func (c *cPicture) Update(ctx context.Context, req *v1.PictureUpdateReq) (res *v1.PictureUpdateRes, err error) {
	return service.Picture().Update(ctx, req)
}

// Review 审核图片
func (c *cPicture) Review(ctx context.Context, req *v1.PictureReviewReq) (res *v1.PictureReviewRes, err error) {
	return service.Picture().Review(ctx, req)
}

// UploadByBatch 批量上传图片
func (c *cPicture) UploadByBatch(ctx context.Context, req *v1.PictureUploadByBatchReq) (res *v1.PictureUploadByBatchRes, err error) {
	return service.Picture().UploadByBatch(ctx, req)
}

// EditByBatch 批量编辑图片
func (c *cPicture) EditByBatch(ctx context.Context, req *v1.PictureEditByBatchReq) (res *v1.PictureEditByBatchRes, err error) {
	return service.Picture().EditByBatch(ctx, req)
}

// SearchByPicture 以图搜图
func (c *cPicture) SearchByPicture(ctx context.Context, req *v1.SearchPictureByPictureReq) (res []v1.SearchPictureByPictureRes, err error) {
	return service.Picture().SearchByPicture(ctx, req)
}

// SearchByColor 按颜色搜索图片
func (c *cPicture) SearchByColor(ctx context.Context, req *v1.SearchPictureByColorReq) (res []v1.SearchPictureByColorRes, err error) {
	return service.Picture().SearchByColor(ctx, req)
}

// CreateAIEditingTask 创建AI编辑任务
func (c *cPicture) CreateAIEditingTask(ctx context.Context, req *v1.CreatePictureAIEditingTaskReq) (res *v1.CreatePictureAIEditingTaskRes, err error) {
	return service.Picture().CreateAIEditingTask(ctx, req)
}

// GetAIEditingTask 获取AI编辑任务
func (c *cPicture) GetAIEditingTask(ctx context.Context, req *v1.GetPictureAIEditingTaskReq) (res *v1.GetPictureAIEditingTaskRes, err error) {
	return service.Picture().GetAIEditingTask(ctx, req)
}

// CreateOutPainting 创建扩图
func (c *cPicture) CreateOutPainting(ctx context.Context, req *v1.CreatePictureOutPaintingReq) (res *v1.CreatePictureOutPaintingRes, err error) {
	return service.Picture().CreateOutPainting(ctx, req)
}
