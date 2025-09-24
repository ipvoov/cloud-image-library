// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "cloud/api/user/v1"
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IPicture interface {
		// CreateAIEditingTask 创建AI编辑任务
		CreateAIEditingTask(ctx context.Context, req *v1.CreatePictureAIEditingTaskReq) (res *v1.CreatePictureAIEditingTaskRes, err error)
		// GetAIEditingTask 获取AI编辑任务
		GetAIEditingTask(ctx context.Context, req *v1.GetPictureAIEditingTaskReq) (res *v1.GetPictureAIEditingTaskRes, err error)
		// UploadByBatch 批量上传图片
		UploadByBatch(ctx context.Context, req *v1.PictureUploadByBatchReq) (res *v1.PictureUploadByBatchRes, err error)
		// EditByBatch 批量编辑图片
		EditByBatch(ctx context.Context, req *v1.PictureEditByBatchReq) (res *v1.PictureEditByBatchRes, err error)
		// Edit 编辑图片
		Edit(ctx context.Context, req *v1.PictureEditReq) (res *v1.PictureEditRes, err error)
		// Delete 删除图片
		Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
		// Update 更新图片
		Update(ctx context.Context, req *v1.PictureUpdateReq) (res *v1.PictureUpdateRes, err error)
		// CreateOutPainting 创建扩图
		CreateOutPainting(ctx context.Context, req *v1.CreatePictureOutPaintingReq) (res *v1.CreatePictureOutPaintingRes, err error)
		// TagCategory 获取图片标签分类
		TagCategory(ctx context.Context, req *v1.PictureTagCategoryReq) (res *v1.PictureTagCategoryRes, err error)
		// Get 获取图片详情
		Get(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureAdminGetRes, err error)
		// GetVO 获取图片详情VO
		GetVO(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureGetRes, err error)
		// ListByPage 分页查询图片
		ListByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureAdminQueryRes, err error)
		// ListVOByPage 分页查询图片VO
		ListVOByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureQueryRes, err error)
		// SearchByPicture 以图搜图
		SearchByPicture(ctx context.Context, req *v1.SearchPictureByPictureReq) (res []v1.SearchPictureByPictureRes, err error)
		// SearchByColor 按颜色搜索图片（基于欧氏距离的相似度搜索）
		SearchByColor(ctx context.Context, req *v1.SearchPictureByColorReq) (res []v1.SearchPictureByColorRes, err error)
		// Review 审核图片
		Review(ctx context.Context, req *v1.PictureReviewReq) (res *v1.PictureReviewRes, err error)
		// Upload 上传图片
		Upload(ctx context.Context, req *v1.PictureUploadReq, file *ghttp.UploadFile) (res *v1.PictureUploadRes, err error)
		// UploadByUrl 通过URL上传图片
		UploadByUrl(ctx context.Context, req *v1.PictureUploadByUrlReq) (res *v1.PictureUploadByUrlRes, err error)
	}
)

var (
	localPicture IPicture
)

func Picture() IPicture {
	if localPicture == nil {
		panic("implement not found for interface IPicture, forgot register?")
	}
	return localPicture
}

func RegisterPicture(i IPicture) {
	localPicture = i
}
