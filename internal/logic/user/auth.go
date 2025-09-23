package user

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserRegister 用户注册
func (s *sUser) UserRegister(ctx context.Context, req *v1.UserRegisterReq) (*v1.UserRegisterRes, error) {
	// 1. 参数校验
	if req.UserAccount == "" || req.UserPassword == "" || req.CheckPassword == "" {
		return nil, gerror.New("参数不能为空")
	}

	// 2. 校验两次密码是否一致
	if req.UserPassword != req.CheckPassword {
		return nil, gerror.New("两次输入的密码不一致")
	}

	// 3. 检查账号是否已存在
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

	// 4. 密码加密
	encryptPassword := s.encryptPassword(req.UserPassword)

	// 5. 插入用户数据
	userId, err := dao.User.Ctx(ctx).Data(do.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		UserName:     consts.DefaultName, // 默认用户名
		UserRole:     consts.DefaultRole, // 默认角色
		IsDelete:     0,
	}).InsertAndGetId()

	if err != nil {
		return nil, gerror.New("注册失败")
	}

	res := &v1.UserRegisterRes{UserId: userId}
	g.Log().Info(ctx, "用户注册成功，返回结果:", res)
	return res, nil
}

// UserLogin 用户登录
func (s *sUser) UserLogin(ctx context.Context, req *v1.UserLoginReq) (*v1.UserLoginRes, error) {
	// 1. 参数校验
	if req.UserAccount == "" || req.UserPassword == "" {
		return nil, gerror.New("账号或密码不能为空")
	}

	// 2. 密码加密
	encryptPassword := s.encryptPassword(req.UserPassword)

	// 3. 查询用户
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(do.User{
		UserAccount:  req.UserAccount,
		UserPassword: encryptPassword,
		IsDelete:     0,
	}).Scan(&user)

	if err != nil {
		return nil, gerror.New("查询用户失败")
	}

	if user == nil {
		return nil, gerror.New("用户不存在或密码错误")
	}

	// 4. 设置登录态
	session := g.RequestFromCtx(ctx).Session
	session.Set(consts.LoginState, user)

	// 5. 返回用户信息
	loginUserVO := s.getLoginUserVO(user)
	return &v1.UserLoginRes{LoginUserVO: loginUserVO}, nil
}

// GetLoginUser 获取当前登录用户
func (s *sUser) GetLoginUser(ctx context.Context, req *v1.GetLoginUserReq) (*v1.GetLoginUserRes, error) {
	// 从上下文获取用户信息（中间件已验证）
	userObj, _ := g.RequestFromCtx(ctx).Session.Get(consts.LoginState)
	if userObj == nil {
		return nil, gerror.New("用户未登录")
	}
	var user *entity.User
	err := gconv.Struct(userObj.Map(), &user)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 从数据库查询最新用户信息
	err = dao.User.Ctx(ctx).Where(do.User{
		Id:       user.Id,
		IsDelete: 0,
	}).Scan(&user)

	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	loginUserVO := s.getLoginUserVO(user)
	return &v1.GetLoginUserRes{LoginUserVO: loginUserVO}, nil
}

// UserLogout 用户登出
func (s *sUser) UserLogout(ctx context.Context, req *v1.UserLogoutReq) (*v1.UserLogoutRes, error) {
	session := g.RequestFromCtx(ctx).Session
	userObj, _ := session.Get(consts.LoginState)

	if userObj == nil {
		return nil, gerror.New("用户未登录")
	}

	// 清除登录态
	session.Remove(consts.LoginState)

	return &v1.UserLogoutRes{Success: true}, nil
}
