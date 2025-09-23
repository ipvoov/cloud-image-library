package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"
)

var SpaceUser = cSpaceUser{}

type cSpaceUser struct{}

// ListMy 获取我的团队空间
func (c *cSpaceUser) ListMy(ctx context.Context, req *v1.SpaceUserListMyReq) (res *v1.SpaceUserListMyRes, err error) {
	return service.SpaceUser().ListMy(ctx, req)
}

// Add 添加空间用户
func (c *cSpaceUser) Add(ctx context.Context, req *v1.SpaceUserAddReq) (res *v1.SpaceUserAddRes, err error) {
	return service.SpaceUser().Add(ctx, req)
}

// Edit 编辑空间用户
func (c *cSpaceUser) Edit(ctx context.Context, req *v1.SpaceUserEditReq) (res *v1.SpaceUserEditRes, err error) {
	return service.SpaceUser().Edit(ctx, req)
}

// Delete 删除空间用户
func (c *cSpaceUser) Delete(ctx context.Context, req *v1.SpaceUserDeleteReq) (res *v1.SpaceUserDeleteRes, err error) {
	return service.SpaceUser().Delete(ctx, req)
}

// Get 获取空间用户
func (c *cSpaceUser) Get(ctx context.Context, req *v1.SpaceUserGetReq) (res *v1.SpaceUserGetRes, err error) {
	return service.SpaceUser().Get(ctx, req)
}

// List 获取空间用户列表
func (c *cSpaceUser) List(ctx context.Context, req *v1.SpaceUserListReq) (res *v1.SpaceUserListRes, err error) {
	return service.SpaceUser().List(ctx, req)
}
