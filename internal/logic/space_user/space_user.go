package space_user

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func init() {
	service.RegisterSpaceUser(New())
}

type sSpaceUser struct{}

func New() *sSpaceUser {
	return &sSpaceUser{}
}

// ListMy 获取我的团队空间
func (s *sSpaceUser) ListMy(ctx context.Context, req *v1.SpaceUserListMyReq) (res *v1.SpaceUserListMyRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 查询用户参与的所有团队空间
	var spaceUsers []*entity.SpaceUser
	err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().UserId, user.Id).
		OrderDesc(dao.SpaceUser.Columns().CreateTime).Scan(&spaceUsers)
	if err != nil {
		return nil, gerror.New("查询我的团队空间失败")
	}

	// 转换为VO，包含空间信息
	records := make([]v1.SpaceUserVO, 0, len(spaceUsers))
	for _, spaceUser := range spaceUsers {
		// 获取空间信息，只包含团队空间
		space := &entity.Space{}
		spaceErr := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceUser.SpaceId).
			Where(dao.Space.Columns().SpaceType, consts.SpaceTypeTeam).
			Where(dao.Space.Columns().IsDelete, 0).Scan(space)
		if spaceErr != nil {
			continue // 跳过不存在或非团队空间
		}

		vo, convertErr := s.entityToSpaceUserVOWithSpace(ctx, spaceUser, space)
		if convertErr != nil {
			g.Log().Errorf(ctx, "转换空间用户VO失败: %v", convertErr)
			continue
		}
		records = append(records, *vo)
	}

	return (*v1.SpaceUserListMyRes)(&records), nil
}

// Add 添加空间用户
func (s *sSpaceUser) Add(ctx context.Context, req *v1.SpaceUserAddReq) (res *v1.SpaceUserAddRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 检查空间是否存在且为团队空间
	space := &entity.Space{}
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).Scan(space)
	if err != nil {
		return nil, gerror.New("空间不存在")
	}

	// 只有团队空间才能添加用户
	if space.SpaceType != consts.SpaceTypeTeam {
		return nil, gerror.New("只有团队空间才能添加用户")
	}

	// 检查权限：只有空间创建者或管理员才能添加用户
	if space.UserId != user.Id && user.UserRole != consts.Admin {
		return nil, gerror.New("无权限添加用户到此空间")
	}

	// 检查目标用户是否存在
	userExists, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, req.UserId).
		Where(dao.User.Columns().IsDelete, 0).Exist()
	if err != nil {
		return nil, gerror.New("检查用户失败")
	}
	if !userExists {
		return nil, gerror.New("目标用户不存在")
	}

	// 检查用户是否已经在空间中
	exists, err := dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, req.SpaceId).
		Where(dao.SpaceUser.Columns().UserId, req.UserId).Exist()
	if err != nil {
		return nil, gerror.New("检查用户关联失败")
	}
	if exists {
		return nil, gerror.New("用户已在空间中")
	}

	// 插入空间用户记录
	result, err := dao.SpaceUser.Ctx(ctx).Data(do.SpaceUser{
		SpaceId:   req.SpaceId,
		UserId:    req.UserId,
		SpaceRole: consts.SpaceRoleViewer,
	}).Insert()
	if err != nil {
		g.Log().Errorf(ctx, "添加空间用户失败: %v", err)
		return nil, gerror.New("添加空间用户失败")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.New("获取空间用户ID失败")
	}

	return &v1.SpaceUserAddRes{
		Id: id,
	}, nil
}

// Edit 编辑空间用户
func (s *sSpaceUser) Edit(ctx context.Context, req *v1.SpaceUserEditReq) (res *v1.SpaceUserEditRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 获取空间用户记录
	spaceUser := &entity.SpaceUser{}
	err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().Id, req.Id).Scan(spaceUser)
	if err != nil {
		return nil, gerror.New("空间用户记录不存在")
	}

	// 获取空间信息
	space := &entity.Space{}
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceUser.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).Scan(space)
	if err != nil {
		return nil, gerror.New("空间不存在")
	}

	// 检查权限：只有空间创建者或管理员才能编辑用户角色
	if space.UserId != user.Id && user.UserRole != consts.Admin {
		return nil, gerror.New("无权限编辑此空间用户")
	}

	// 不能修改空间创建者的角色
	if spaceUser.UserId == space.UserId {
		return nil, gerror.New("不能修改空间创建者的角色")
	}

	// 更新空间用户角色
	_, err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().Id, req.Id).
		Data(do.SpaceUser{
			SpaceRole:  req.SpaceRole,
			UpdateTime: gtime.Now(),
		}).Update()
	if err != nil {
		g.Log().Errorf(ctx, "编辑空间用户失败: %v", err)
		return nil, gerror.New("编辑空间用户失败")
	}

	return &v1.SpaceUserEditRes{
		Success: true,
	}, nil
}

// Delete 删除空间用户
func (s *sSpaceUser) Delete(ctx context.Context, req *v1.SpaceUserDeleteReq) (res *v1.SpaceUserDeleteRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 获取空间用户记录
	spaceUser := &entity.SpaceUser{}
	err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().Id, req.Id).Scan(spaceUser)
	if err != nil {
		return nil, gerror.New("空间用户记录不存在")
	}

	// 获取空间信息
	space := &entity.Space{}
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceUser.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).Scan(space)
	if err != nil {
		return nil, gerror.New("空间不存在")
	}

	// 检查权限：只有空间创建者或管理员才能删除用户，或者用户可以删除自己
	if space.UserId != user.Id && user.UserRole != consts.Admin && spaceUser.UserId != user.Id {
		return nil, gerror.New("无权限删除此空间用户")
	}

	// 不能删除空间创建者
	if spaceUser.UserId == space.UserId {
		return nil, gerror.New("不能删除空间创建者")
	}

	// 删除空间用户记录
	_, err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().Id, req.Id).Delete()
	if err != nil {
		g.Log().Errorf(ctx, "删除空间用户失败: %v", err)
		return nil, gerror.New("删除空间用户失败")
	}

	return &v1.SpaceUserDeleteRes{
		Success: true,
	}, nil
}

// Get 获取空间用户
func (s *sSpaceUser) Get(ctx context.Context, req *v1.SpaceUserGetReq) (res *v1.SpaceUserGetRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	var spaceUser *entity.SpaceUser

	// 根据不同参数查询空间用户
	query := dao.SpaceUser.Ctx(ctx)
	if req.Id > 0 {
		query = query.Where(dao.SpaceUser.Columns().Id, req.Id)
	} else if req.SpaceId > 0 && req.UserId > 0 {
		query = query.Where(dao.SpaceUser.Columns().SpaceId, req.SpaceId).
			Where(dao.SpaceUser.Columns().UserId, req.UserId)
	} else {
		return nil, gerror.New("参数不足")
	}

	spaceUser = &entity.SpaceUser{}
	err = query.Scan(spaceUser)
	if err != nil {
		return nil, gerror.New("空间用户记录不存在")
	}

	// 检查权限：只有空间创建者、管理员或相关用户才能查看
	space := &entity.Space{}
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceUser.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).Scan(space)
	if err != nil {
		return nil, gerror.New("空间不存在")
	}

	if space.UserId != user.Id && user.UserRole != consts.Admin && spaceUser.UserId != user.Id {
		return nil, gerror.New("无权限查看此空间用户")
	}

	return &v1.SpaceUserGetRes{
		SpaceUser: s.entityToSpaceUser(spaceUser),
	}, nil
}

// List 获取空间用户列表
func (s *sSpaceUser) List(ctx context.Context, req *v1.SpaceUserListReq) (res *v1.SpaceUserListRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 检查空间是否存在
	space := &entity.Space{}
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).Scan(space)
	if err != nil {
		return nil, gerror.New("空间不存在")
	}

	// 检查权限：只有空间成员、创建者或管理员才能查看成员列表
	hasPermission := false
	if space.UserId == user.Id || user.UserRole == consts.Admin {
		hasPermission = true
	} else {
		// 检查是否是空间成员
		exists, checkErr := dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, req.SpaceId).
			Where(dao.SpaceUser.Columns().UserId, user.Id).Exist()
		if checkErr == nil && exists {
			hasPermission = true
		}
	}

	if !hasPermission {
		return nil, gerror.New("无权限查看此空间成员列表")
	}

	// 查询空间用户列表
	var spaceUsers []*entity.SpaceUser
	err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, req.SpaceId).
		OrderAsc(dao.SpaceUser.Columns().CreateTime).Scan(&spaceUsers)
	if err != nil {
		return nil, gerror.New("查询空间用户列表失败")
	}

	// 转换为VO
	records := make([]v1.SpaceUserVO, 0, len(spaceUsers)+1) // +1 为空间所有者预留位置

	// 首先添加空间所有者
	ownerVO, ownerErr := s.createSpaceOwnerVO(ctx, space)
	if ownerErr != nil {
		g.Log().Errorf(ctx, "创建空间所有者VO失败: %v", ownerErr)
	} else {
		records = append(records, *ownerVO)
	}

	// 然后添加其他成员
	for _, spaceUser := range spaceUsers {
		// 跳过空间所有者，避免重复
		if spaceUser.UserId == space.UserId {
			continue
		}

		vo, convertErr := s.entityToSpaceUserVO(ctx, spaceUser)
		if convertErr != nil {
			g.Log().Errorf(ctx, "转换空间用户VO失败: %v", convertErr)
			continue
		}
		records = append(records, *vo)
	}

	return (*v1.SpaceUserListRes)(&records), nil
}

// entityToSpaceUser 将entity转换为SpaceUser
func (s *sSpaceUser) entityToSpaceUser(spaceUser *entity.SpaceUser) *v1.SpaceUser {
	return &v1.SpaceUser{
		Id:         spaceUser.Id,
		UserId:     spaceUser.UserId,
		SpaceId:    spaceUser.SpaceId,
		SpaceRole:  spaceUser.SpaceRole,
		CreateTime: spaceUser.CreateTime.Format(consts.Y_m_d_His),
		UpdateTime: spaceUser.UpdateTime.Format(consts.Y_m_d_His),
	}
}

// entityToSpaceUserVO 将entity转换为SpaceUserVO
func (s *sSpaceUser) entityToSpaceUserVO(ctx context.Context, spaceUser *entity.SpaceUser) (*v1.SpaceUserVO, error) {
	// 获取用户信息
	user := &entity.User{}
	err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, spaceUser.UserId).
		Where(dao.User.Columns().IsDelete, 0).Scan(user)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	return &v1.SpaceUserVO{
		Id:         spaceUser.Id,
		UserId:     spaceUser.UserId,
		SpaceId:    spaceUser.SpaceId,
		SpaceRole:  spaceUser.SpaceRole,
		CreateTime: spaceUser.CreateTime.Format(consts.Y_m_d_His),
		UpdateTime: spaceUser.UpdateTime.Format(consts.Y_m_d_His),
		User: &v1.UserVO{
			Id:         user.Id,
			UserName:   user.UserName,
			UserAvatar: user.UserAvatar,
		},
	}, nil
}

// entityToSpaceUserVOWithSpace 将entity转换为SpaceUserVO（包含空间信息）
func (s *sSpaceUser) entityToSpaceUserVOWithSpace(ctx context.Context, spaceUser *entity.SpaceUser, space *entity.Space) (*v1.SpaceUserVO, error) {
	// 获取用户信息
	user := &entity.User{}
	err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, spaceUser.UserId).
		Where(dao.User.Columns().IsDelete, 0).Scan(user)
	if err != nil {
		return nil, gerror.New("获取用户信息失败")
	}

	// 获取用户在空间中的权限列表
	permissions := s.getUserPermissions(ctx, space.UserId, spaceUser.UserId, spaceUser.SpaceRole)

	return &v1.SpaceUserVO{
		Id:         spaceUser.Id,
		UserId:     spaceUser.UserId,
		SpaceId:    spaceUser.SpaceId,
		SpaceRole:  spaceUser.SpaceRole,
		CreateTime: spaceUser.CreateTime.Format(consts.Y_m_d_His),
		UpdateTime: spaceUser.UpdateTime.Format(consts.Y_m_d_His),
		User: &v1.UserVO{
			Id:         user.Id,
			UserName:   user.UserName,
			UserAvatar: user.UserAvatar,
		},
		Space: &v1.SpaceVO{
			Id:             space.Id,
			SpaceName:      space.SpaceName,
			SpaceLevel:     space.SpaceLevel,
			MaxSize:        space.MaxSize,
			MaxCount:       space.MaxCount,
			TotalSize:      space.TotalSize,
			TotalCount:     space.TotalCount,
			UserId:         space.UserId,
			CreateTime:     space.CreateTime.Format(consts.Y_m_d_His),
			EditTime:       space.EditTime.Format(consts.Y_m_d_His),
			UpdateTime:     space.UpdateTime.Format(consts.Y_m_d_His),
			SpaceType:      space.SpaceType,
			PermissionList: permissions,
		},
	}, nil
}

// getUserPermissions 获取用户在空间中的权限列表
func (s *sSpaceUser) getUserPermissions(ctx context.Context, spaceOwnerId, userId int64, spaceRole string) []string {
	permissions := make([]string, 0)

	// 空间创建者拥有所有权限
	if spaceOwnerId == userId {
		return []string{
			"spaceUser:manage",
			"picture:view",
			"picture:upload",
			"picture:edit",
			"picture:delete",
		}
	}

	// 根据角色分配权限
	switch spaceRole {
	case "admin":
		permissions = []string{
			"spaceUser:manage",
			"picture:view",
			"picture:upload",
			"picture:edit",
			"picture:delete",
		}
	case "editor":
		permissions = []string{
			"picture:view",
			"picture:upload",
			"picture:edit",
		}
	case "viewer":
		permissions = []string{
			"picture:view",
		}
	}

	return permissions
}

// createSpaceOwnerVO 创建空间所有者的VO
func (s *sSpaceUser) createSpaceOwnerVO(ctx context.Context, space *entity.Space) (*v1.SpaceUserVO, error) {
	// 获取空间所有者的用户信息
	user := &entity.User{}
	err := dao.User.Ctx(ctx).Where(dao.User.Columns().Id, space.UserId).
		Where(dao.User.Columns().IsDelete, 0).Scan(user)
	if err != nil {
		return nil, gerror.New("获取空间所有者信息失败")
	}

	return &v1.SpaceUserVO{
		Id:         0, // 空间所有者没有在space_user表中的记录，所以ID为0
		UserId:     space.UserId,
		SpaceId:    space.Id,
		SpaceRole:  "owner", // 空间所有者角色
		CreateTime: space.CreateTime.Format(consts.Y_m_d_His),
		UpdateTime: space.UpdateTime.Format(consts.Y_m_d_His),
		User: &v1.UserVO{
			Id:         user.Id,
			UserName:   user.UserName,
			UserAvatar: user.UserAvatar,
		},
	}, nil
}
