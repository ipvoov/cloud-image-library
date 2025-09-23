package controller

import (
	"context"

	v1 "cloud/api/user/v1"
	"cloud/internal/service"
)

var User = cUser{}

type cUser struct{}

func (c *cUser) Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error) {
	return service.User().UserRegister(ctx, req)
}

// Login 用户登录
func (c *cUser) Login(ctx context.Context, req *v1.UserLoginReq) (res *v1.UserLoginRes, err error) {
	return service.User().UserLogin(ctx, req)
}

// GetLoginUser 获取当前登录用户
func (c *cUser) GetLoginUser(ctx context.Context, req *v1.GetLoginUserReq) (res *v1.GetLoginUserRes, err error) {
	return service.User().GetLoginUser(ctx, req)
}

// Logout 用户登出
func (c *cUser) Logout(ctx context.Context, req *v1.UserLogoutReq) (res *v1.UserLogoutRes, err error) {
	return service.User().UserLogout(ctx, req)
}

// Add 添加用户
func (c *cUser) Add(ctx context.Context, req *v1.UserAddReq) (res *v1.UserAddRes, err error) {
	return service.User().UserAdd(ctx, req)
}

// Update 更新用户
func (c *cUser) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	return service.User().UserUpdate(ctx, req)
}

// Delete 删除用户
func (c *cUser) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	return service.User().UserDelete(ctx, req)
}

// GetById 根据ID获取用户
func (c *cUser) GetById(ctx context.Context, req *v1.GetUserByIdReq) (res *v1.GetUserByIdRes, err error) {
	return service.User().GetUserById(ctx, req)
}

// ListByPage 分页查询用户
func (c *cUser) ListByPage(ctx context.Context, req *v1.UserQueryReq) (res *v1.UserQueryRes, err error) {
	return service.User().ListUserByPage(ctx, req)
}

// GetProfile 获取用户详细信息
func (c *cUser) GetProfile(ctx context.Context, req *v1.GetUserProfileReq) (res *v1.GetUserProfileRes, err error) {
	return service.User().GetUserProfile(ctx, req)
}

// UpdateProfile 更新用户资料
func (c *cUser) UpdateProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (res *v1.UpdateUserProfileRes, err error) {
	return service.User().UpdateUserProfile(ctx, req)
}
