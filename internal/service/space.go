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
	ISpace interface {
		// Add 添加空间
		Add(ctx context.Context, req *v1.SpaceAddReq) (res *v1.SpaceAddRes, err error)
		// Edit 编辑空间
		Edit(ctx context.Context, req *v1.SpaceEditReq) (res *v1.SpaceEditRes, err error)
		// Update 更新空间
		Update(ctx context.Context, req *v1.SpaceUpdateReq) (res *v1.SpaceUpdateRes, err error)
		// Delete 删除空间
		Delete(ctx context.Context, req *v1.SpaceDeleteReq) (res *v1.SpaceDeleteRes, err error)
		// Get 获取空间详情
		Get(ctx context.Context, req *v1.SpaceGetReq) (res *v1.SpaceGetRes, err error)
		// GetVO 获取空间详情VO
		GetVO(ctx context.Context, req *v1.SpaceGetVOReq) (res *v1.SpaceGetVORes, err error)
		// ListByPage 分页查询空间
		ListByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryRes, err error)
		// ListVOByPage 分页查询空间VO
		ListVOByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryVORes, err error)
		// ListLevel 获取空间级别列表
		ListLevel(ctx context.Context, req *v1.SpaceLevelListReq) (res *v1.SpaceLevelListRes, err error)
	}
)

var (
	localSpace ISpace
)

func Space() ISpace {
	if localSpace == nil {
		panic("implement not found for interface ISpace, forgot register?")
	}
	return localSpace
}

func RegisterSpace(i ISpace) {
	localSpace = i
}
