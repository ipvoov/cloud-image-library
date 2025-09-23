package user

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// UserAdd 添加用户
func (s *sUser) UserAdd(ctx context.Context, req *v1.UserAddReq) (*v1.UserAddRes, error) {
	// 1. 参数校验
	if req.UserAccount == "" || req.UserPassword == "" {
		return nil, gerror.New("账号和密码不能为空")
	}

	// 2. 检查账号是否已存在
	count, err := dao.User.Ctx(ctx).Where(do.User{
		UserAccount: req.UserAccount,
		IsDelete:    0,
	}).Count()
	if err != nil {
		return nil, gerror.New("查询用户失败")
	}
	if count > 0 {
		return nil, gerror.New("账号已存在")
	}

	// 3. 密码加密
	encryptPassword := s.encryptPassword(req.UserPassword)

	// 4. 设置默认值
	userName := req.UserName
	if userName == "" {
		userName = consts.DefaultName
	}
	userRole := req.UserRole
	if userRole == "" {
		userRole = consts.DefaultRole
	}

	// 5. 插入用户数据
	userId, err := dao.User.Ctx(ctx).Data(do.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		UserName:     userName,
		UserRole:     userRole,
		IsDelete:     0,
	}).InsertAndGetId()

	if err != nil {
		return nil, gerror.New("添加用户失败")
	}

	return &v1.UserAddRes{UserId: userId}, nil
}

// UserUpdate 更新用户
func (s *sUser) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (*v1.UserUpdateRes, error) {
	// 1. 检查用户是否存在
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(do.User{
		Id:       req.Id,
		IsDelete: 0,
	}).Scan(&user)

	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 2. 更新用户信息
	updateData := do.User{}
	if req.UserName != "" {
		updateData.UserName = req.UserName
	}
	if req.UserAvatar != "" {
		updateData.UserAvatar = req.UserAvatar
	}
	if req.UserRole != "" {
		updateData.UserRole = req.UserRole
	}

	_, err = dao.User.Ctx(ctx).Where(do.User{Id: req.Id}).Data(updateData).Update()
	if err != nil {
		return nil, gerror.New("更新用户失败")
	}

	return &v1.UserUpdateRes{Success: true}, nil
}

// UserDelete 删除用户
func (s *sUser) UserDelete(ctx context.Context, req *v1.DeleteReq) (*v1.DeleteRes, error) {
	// 1. 检查用户是否存在
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(do.User{
		Id:       req.Id,
		IsDelete: 0,
	}).Scan(&user)

	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	// 2. 软删除用户
	_, err = dao.User.Ctx(ctx).Where(do.User{Id: req.Id}).Data(do.User{
		IsDelete: 1,
	}).Update()

	if err != nil {
		return nil, gerror.New("删除用户失败")
	}

	return &v1.DeleteRes{Success: true}, nil
}
