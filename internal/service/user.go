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
	IUser interface {
		// UserRegister 用户注册
		UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*v1.UserRegisterRes, error)
		// UserLogin 用户登录
		UserLogin(ctx context.Context, req *v1.UserLoginReq) (*v1.UserLoginRes, error)
		// GetLoginUser 获取当前登录用户
		GetLoginUser(ctx context.Context, req *v1.GetLoginUserReq) (*v1.GetLoginUserRes, error)
		// UserLogout 用户登出
		UserLogout(ctx context.Context, req *v1.UserLogoutReq) (*v1.UserLogoutRes, error)
		// UserAdd 添加用户
		UserAdd(ctx context.Context, req *v1.UserAddReq) (*v1.UserAddRes, error)
		// UserUpdate 更新用户
		UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (*v1.UserUpdateRes, error)
		// UserDelete 删除用户
		UserDelete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error)
		// GetUserById 根据ID获取用户
		GetUserById(ctx context.Context, req *v1.GetUserByIdReq) (*v1.GetUserByIdRes, error)
		// ListUserByPage 分页查询用户
		ListUserByPage(ctx context.Context, req *v1.UserQueryReq) (*v1.UserQueryRes, error)
		// GetUserProfile 获取用户详细信息
		GetUserProfile(ctx context.Context, req *v1.GetUserProfileReq) (*v1.GetUserProfileRes, error)
		// UpdateUserProfile 更新用户资料
		UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (*v1.UpdateUserProfileRes, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
