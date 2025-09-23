package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"
)

var Space = cSpace{}

type cSpace struct{}

// Add 添加空间
func (c *cSpace) Add(ctx context.Context, req *v1.SpaceAddReq) (res *v1.SpaceAddRes, err error) {
	return service.Space().Add(ctx, req)
}

// Edit 编辑空间
func (c *cSpace) Edit(ctx context.Context, req *v1.SpaceEditReq) (res *v1.SpaceEditRes, err error) {
	return service.Space().Edit(ctx, req)
}

// Update 更新空间
func (c *cSpace) Update(ctx context.Context, req *v1.SpaceUpdateReq) (res *v1.SpaceUpdateRes, err error) {
	return service.Space().Update(ctx, req)
}

// Delete 删除空间
func (c *cSpace) Delete(ctx context.Context, req *v1.SpaceDeleteReq) (res *v1.SpaceDeleteRes, err error) {
	return service.Space().Delete(ctx, req)
}

// Get 获取空间详情
func (c *cSpace) Get(ctx context.Context, req *v1.SpaceGetReq) (res *v1.SpaceGetRes, err error) {
	return service.Space().Get(ctx, req)
}

// GetVO 获取空间详情VO
func (c *cSpace) GetVO(ctx context.Context, req *v1.SpaceGetVOReq) (res *v1.SpaceGetVORes, err error) {
	return service.Space().GetVO(ctx, req)
}

// ListByPage 分页查询空间
func (c *cSpace) ListByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryRes, err error) {
	return service.Space().ListByPage(ctx, req)
}

// ListVOByPage 分页查询空间VO
func (c *cSpace) ListVOByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryVORes, err error) {
	return service.Space().ListVOByPage(ctx, req)
}

// ListLevel 获取空间级别列表
func (c *cSpace) ListLevel(ctx context.Context, req *v1.SpaceLevelListReq) (res *v1.SpaceLevelListRes, err error) {
	return service.Space().ListLevel(ctx, req)
}
