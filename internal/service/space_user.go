// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "cloud/api/user/v1"
	"context"
)

type (
	ISpaceUser interface {
		// ListMy 获取我的团队空间
		ListMy(ctx context.Context, req *v1.SpaceUserListMyReq) (res *v1.SpaceUserListMyRes, err error)
		// Add 添加空间用户
		Add(ctx context.Context, req *v1.SpaceUserAddReq) (res *v1.SpaceUserAddRes, err error)
		// Edit 编辑空间用户
		Edit(ctx context.Context, req *v1.SpaceUserEditReq) (res *v1.SpaceUserEditRes, err error)
		// Delete 删除空间用户
		Delete(ctx context.Context, req *v1.SpaceUserDeleteReq) (res *v1.SpaceUserDeleteRes, err error)
		// Get 获取空间用户
		Get(ctx context.Context, req *v1.SpaceUserGetReq) (res *v1.SpaceUserGetRes, err error)
		// List 获取空间用户列表
		List(ctx context.Context, req *v1.SpaceUserListReq) (res *v1.SpaceUserListRes, err error)
	}
)

var (
	localSpaceUser ISpaceUser
)

func SpaceUser() ISpaceUser {
	if localSpaceUser == nil {
		panic("implement not found for interface ISpaceUser, forgot register?")
	}
	return localSpaceUser
}

func RegisterSpaceUser(i ISpaceUser) {
	localSpaceUser = i
}
