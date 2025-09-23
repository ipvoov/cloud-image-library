package user

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// GetUserById 根据ID获取用户
func (s *sUser) GetUserById(ctx context.Context, req *v1.GetUserByIdReq) (*v1.GetUserByIdRes, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Where(do.User{
		Id:       req.Id,
		IsDelete: 0,
	}).Scan(&user)

	if err != nil || user == nil {
		return nil, gerror.New("用户不存在")
	}

	loginUserVO := s.getLoginUserVO(user)
	return &v1.GetUserByIdRes{LoginUserVO: loginUserVO}, nil
}

// ListUserByPage 分页查询用户
func (s *sUser) ListUserByPage(ctx context.Context, req *v1.UserQueryReq) (*v1.UserQueryRes, error) {

	u := dao.User.Columns()
	query := dao.User.Ctx(ctx).Where(u.IsDelete, 0)

	if req.UserAccount != "" {
		query = query.WhereLike(dao.User.Columns().UserAccount, req.UserAccount+"%")
	}
	if req.UserName != "" {
		query = query.WhereLike(dao.User.Columns().UserName, req.UserName+"%")
	}
	if req.UserRole != "" {
		query = query.Where(dao.User.Columns().UserRole, req.UserRole)
	}

	// 3. 查询总数
	total, err := query.Count()
	if err != nil {
		return nil, gerror.New("查询用户总数失败")
	}

	// 排序
	var orderBy string
	if req.SortField == "" {
		orderBy = "createTime DESC"
	} else {
		var sort string
		switch req.SortOrder {
		case "ascend":
			sort = "ASC"
		case "descend":
			sort = "DESC"
		default:
			sort = "DESC"
		}
		orderBy = req.SortField + " " + sort
	}

	// 4. 分页查询
	var users []*entity.User
	err = query.Page(req.Current, req.PageSize).Order(orderBy).Scan(&users)
	if err != nil {
		return nil, gerror.New("查询用户列表失败")
	}

	// 5. 转换为VO
	var records []v1.LoginUserVO
	for _, user := range users {
		loginUserVO := s.getLoginUserVO(user)
		records = append(records, *loginUserVO)
	}

	return &v1.UserQueryRes{
		Records: records,
		Total:   total,
		Current: req.Current,
		Size:    req.PageSize,
	}, nil
}

// GetUserProfile 获取用户详细信息
func (s *sUser) GetUserProfile(ctx context.Context, req *v1.GetUserProfileReq) (*v1.GetUserProfileRes, error) {
	// 从session获取当前登录用户
	loginUser, err := s.GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 获取用户详细信息
	var user *entity.User
	err = dao.User.Ctx(ctx).Where(do.User{
		Id:       loginUser.Id,
		IsDelete: 0,
	}).Scan(&user)

	if err != nil || user == nil {
		return nil, gerror.New("用户信息不存在")
	}

	// 统计用户数据
	pictureCount, _ := dao.Picture.Ctx(ctx).Where(do.Picture{
		UserId:   user.Id,
		IsDelete: 0,
	}).Count()

	spaceCount, _ := dao.Space.Ctx(ctx).Where(do.Space{
		UserId:   user.Id,
		IsDelete: 0,
	}).Count()

	teamCount, _ := dao.SpaceUser.Ctx(ctx).Where(do.SpaceUser{
		UserId: user.Id,
	}).Count()

	// 构建返回数据
	userProfileVO := &v1.UserProfileVO{
		Id:            user.Id,
		UserAccount:   user.UserAccount,
		UserName:      user.UserName,
		UserAvatar:    user.UserAvatar,
		UserProfile:   user.UserProfile,
		UserRole:      user.UserRole,
		CreateTime:    user.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateTime:    user.UpdateTime.Format("2006-01-02 15:04:05"),
		VipExpireTime: "",
		VipCode:       user.VipCode,
		VipNumber:     user.VipNumber,
		PictureCount:  int64(pictureCount),
		SpaceCount:    int64(spaceCount),
		TeamCount:     int64(teamCount),
	}

	if user.VipExpireTime != nil {
		userProfileVO.VipExpireTime = user.VipExpireTime.Format("2006-01-02 15:04:05")
	}

	return &v1.GetUserProfileRes{UserProfileVO: userProfileVO}, nil
}

// UpdateUserProfile 更新用户资料
func (s *sUser) UpdateUserProfile(ctx context.Context, req *v1.UpdateUserProfileReq) (*v1.UpdateUserProfileRes, error) {
	// 从session获取当前登录用户
	loginUser, err := s.GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	_, err = dao.User.Ctx(ctx).Where(do.User{
		Id:       loginUser.Id,
		IsDelete: 0,
	}).Data(do.User{
		UserName:    req.UserName,
		UserAvatar:  req.UserAvatar,
		UserProfile: req.UserProfile,
	}).Update()

	if err != nil {
		return nil, gerror.New("更新用户信息失败")
	}

	return &v1.UpdateUserProfileRes{Success: true}, nil
}
