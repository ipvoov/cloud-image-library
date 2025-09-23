package space

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

// Get 获取空间详情
func (s *sSpace) Get(ctx context.Context, req *v1.SpaceGetReq) (res *v1.SpaceGetRes, err error) {
	var space *entity.Space
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.Id).
		Where(dao.Space.Columns().IsDelete, 0).Scan(&space)
	if err != nil || space == nil {
		return nil, gerror.New("空间不存在")
	}

	return &v1.SpaceGetRes{
		Space: s.entityToSpace(ctx, space),
	}, nil
}

// GetVO 获取空间详情VO
func (s *sSpace) GetVO(ctx context.Context, req *v1.SpaceGetVOReq) (res *v1.SpaceGetVORes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	var spaceId int
	if req.Id == "[object Object]" {
		space := dao.Space.Columns()
		id, err := dao.Space.Ctx(ctx).Fields(space.Id).Where(space.UserId, user.Id).Value()
		if err != nil || id.IsEmpty() {
			return nil, gerror.New("空间不存在")
		}
		spaceId = id.Int()
	}

	if spaceId != 0 {
		req.Id = gconv.String(spaceId)
	}

	var space *entity.Space
	err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, gconv.Int64(req.Id)).
		Where(dao.Space.Columns().IsDelete, 0).Scan(&space)
	if err != nil || space == nil {
		return nil, gerror.New("空间不存在")
	}

	// 检查访问权限
	if space.UserId != user.Id && user.UserRole != consts.Admin && space.SpaceType == consts.SpaceTypePrivate {
		return nil, gerror.New("无权限访问此空间")
	}

	return &v1.SpaceGetVORes{
		SpaceVO: s.entityToSpaceVO(ctx, space, user.Id),
	}, nil
}

// ListByPage 分页查询空间
func (s *sSpace) ListByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 只有管理员可以查看所有空间
	if user.UserRole != consts.Admin {
		return nil, gerror.New("无权限查看空间列表")
	}

	// 构建查询条件
	query := dao.Space.Ctx(ctx).Where(dao.Space.Columns().IsDelete, 0)

	if req.SpaceId != "" && req.SpaceId != "[object Object]" {
		query = query.Where(dao.Space.Columns().Id, gconv.Int64(req.SpaceId))
	} else if req.SpaceId == "[object Object]" {
		var spaceId int
		space := dao.Space.Columns()
		id, err := dao.Space.Ctx(ctx).Fields(space.Id).Where(space.UserId, user.Id).Value()
		if err != nil || id.IsEmpty() {
			return nil, gerror.New("空间不存在")
		}
		spaceId = id.Int()
		query = query.Where(dao.Space.Columns().Id, gconv.Int64(spaceId))
	}
	if req.SpaceName != "" {
		query = query.WhereLike(dao.Space.Columns().SpaceName, "%"+req.SpaceName+"%")
	}
	if req.UserId > 0 {
		query = query.Where(dao.Space.Columns().UserId, req.UserId)
	}

	// 查询总数
	total, err := query.Count()
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	// 分页查询
	var spaces []entity.Space
	err = query.Page(req.Current, req.PageSize).
		Order(dao.Space.Columns().CreateTime + " DESC").
		Scan(&spaces)
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	var records []v1.Space
	for _, space := range spaces {
		records = append(records, *s.entityToSpace(ctx, &space))
	}

	// 计算总页数
	pages := (total + req.PageSize - 1) / req.PageSize

	return &v1.SpaceQueryRes{
		Records: records,
		PageInfo: &v1.PageInfo{
			Current: req.Current,
			Size:    req.PageSize,
			Total:   total,
			Pages:   pages,
		},
	}, nil
}

// ListVOByPage 分页查询空间VO
func (s *sSpace) ListVOByPage(ctx context.Context, req *v1.SpaceQueryReq) (res *v1.SpaceQueryVORes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 获取用户有权限访问的空间ID列表
	accessibleSpaceIds, err := s.getUserAccessibleSpaceIds(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	if len(accessibleSpaceIds) == 0 {
		// 用户没有任何可访问的空间
		return &v1.SpaceQueryVORes{
			Records: []v1.SpaceVO{},
			PageInfo: &v1.PageInfo{
				Current: req.Current,
				Size:    req.PageSize,
				Total:   0,
				Pages:   0,
			},
		}, nil
	}

	// 构建查询条件 - 用户只能查看自己有权限的空间
	query := dao.Space.Ctx(ctx).Where(dao.Space.Columns().IsDelete, 0).
		WhereIn(dao.Space.Columns().Id, accessibleSpaceIds)

	if req.SpaceId != "" && req.SpaceId != "[object Object]" {
		query = query.Where(dao.Space.Columns().Id, gconv.Int64(req.SpaceId))
	} else if req.SpaceId == "[object Object]" {
		var spaceId int
		space := dao.Space.Columns()
		id, err := dao.Space.Ctx(ctx).Fields(space.Id).Where(space.UserId, user.Id).Value()
		if err != nil || id.IsEmpty() || id.Int() == 0 {
			return nil, gerror.New("空间不存在")
		}
		spaceId = id.Int()
		query = query.Where(dao.Space.Columns().Id, gconv.Int64(spaceId))
	}
	if req.SpaceName != "" {
		query = query.WhereLike(dao.Space.Columns().SpaceName, "%"+req.SpaceName+"%")
	}

	// 查询总数
	total, err := query.Count()
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	// 分页查询
	var spaces []entity.Space
	err = query.Page(req.Current, req.PageSize).
		Order(dao.Space.Columns().CreateTime + " DESC").
		Scan(&spaces)
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	var records []v1.SpaceVO
	for _, space := range spaces {
		records = append(records, *s.entityToSpaceVO(ctx, &space, user.Id))
	}

	// 计算总页数
	pages := (total + req.PageSize - 1) / req.PageSize

	return &v1.SpaceQueryVORes{
		Records: records,
		PageInfo: &v1.PageInfo{
			Current: req.Current,
			Size:    req.PageSize,
			Total:   total,
			Pages:   pages,
		},
	}, nil
}

// ListLevel 获取空间级别列表
func (s *sSpace) ListLevel(ctx context.Context, req *v1.SpaceLevelListReq) (res *v1.SpaceLevelListRes, err error) {
	records := []v1.SpaceLevel{
		{
			Level:       0,
			Name:        "普通版",
			Description: "适合个人用户，提供基础功能",
			MaxSize:     100 * 1024 * 1024, // 100MB
			MaxCount:    100,
		},
		{
			Level:       1,
			Name:        "专业版",
			Description: "适合小团队，提供更多存储空间",
			MaxSize:     1024 * 1024 * 1024, // 1GB
			MaxCount:    1000,
		},
		{
			Level:       2,
			Name:        "旗舰版",
			Description: "适合大团队，提供最大存储空间",
			MaxSize:     10 * 1024 * 1024 * 1024, // 10GB
			MaxCount:    10000,
		},
	}

	return &v1.SpaceLevelListRes{
		Records: records,
	}, nil
}

// entityToSpace 将entity转换为Space
func (s *sSpace) entityToSpace(ctx context.Context, space *entity.Space) *v1.Space {
	return &v1.Space{
		Id:         space.Id,
		SpaceName:  space.SpaceName,
		SpaceLevel: space.SpaceLevel,
		MaxSize:    space.MaxSize,
		MaxCount:   space.MaxCount,
		TotalSize:  space.TotalSize,
		TotalCount: space.TotalCount,
		UserId:     space.UserId,
		CreateTime: space.CreateTime.Format(consts.Y_m_d_His),
		EditTime:   space.EditTime.Format(consts.Y_m_d_His),
		UpdateTime: space.UpdateTime.Format(consts.Y_m_d_His),
		IsDelete:   space.IsDelete,
		SpaceType:  space.SpaceType,
	}
}

// entityToSpaceVO 将entity转换为SpaceVO
func (s *sSpace) entityToSpaceVO(ctx context.Context, space *entity.Space, userId int64) *v1.SpaceVO {
	// 获取用户权限列表
	permissions := s.getUserPermissions(ctx, space, userId)

	// 确定显示的空间名称（如果是私有空间，显示为"私有空间"，团队空间显示原名）
	displayName := space.SpaceName
	if space.SpaceType == consts.SpaceTypePrivate {
		displayName = space.SpaceName + " (私有空间)"
	} else if space.SpaceType == consts.SpaceTypeTeam {
		displayName = space.SpaceName + " (团队空间)"
	}

	return &v1.SpaceVO{
		Id:             space.Id,
		SpaceName:      displayName,
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
	}
}

// getUserAccessibleSpaceIds 获取用户有权限访问的空间ID列表
func (s *sSpace) getUserAccessibleSpaceIds(ctx context.Context, userId int64) ([]int64, error) {
	var spaceIds []int64

	// 1. 查询用户创建的所有空间（私有空间）
	var ownedSpaces []entity.Space
	err := dao.Space.Ctx(ctx).Fields(dao.Space.Columns().Id).
		Where(dao.Space.Columns().UserId, userId).
		Where(dao.Space.Columns().IsDelete, 0).
		Scan(&ownedSpaces)
	if err != nil {
		return nil, gerror.New("查询用户空间失败")
	}

	// 添加用户创建的空间ID
	for _, space := range ownedSpaces {
		spaceIds = append(spaceIds, space.Id)
	}

	// 2. 查询用户被邀请加入的团队空间
	var spaceUsers []entity.SpaceUser
	err = dao.SpaceUser.Ctx(ctx).Fields(dao.SpaceUser.Columns().SpaceId).
		Where(dao.SpaceUser.Columns().UserId, userId).
		Scan(&spaceUsers)
	if err != nil {
		return nil, gerror.New("查询用户空间权限失败")
	}

	// 添加用户被邀请的空间ID（需要确保这些空间存在且未删除）
	if len(spaceUsers) > 0 {
		var invitedSpaceIds []int64
		for _, spaceUser := range spaceUsers {
			invitedSpaceIds = append(invitedSpaceIds, spaceUser.SpaceId)
		}

		var invitedSpaces []entity.Space
		err = dao.Space.Ctx(ctx).Fields(dao.Space.Columns().Id).
			WhereIn(dao.Space.Columns().Id, invitedSpaceIds).
			Where(dao.Space.Columns().IsDelete, 0).
			Scan(&invitedSpaces)
		if err != nil {
			return nil, gerror.New("查询邀请空间失败")
		}

		// 添加有效的邀请空间ID
		for _, space := range invitedSpaces {
			// 避免重复添加（用户既是创建者又是成员的情况）
			found := false
			for _, existingId := range spaceIds {
				if existingId == space.Id {
					found = true
					break
				}
			}
			if !found {
				spaceIds = append(spaceIds, space.Id)
			}
		}
	}

	return spaceIds, nil
}

// getUserPermissions 获取用户在空间中的权限
func (s *sSpace) getUserPermissions(ctx context.Context, space *entity.Space, userId int64) []string {
	// 如果是空间所有者，拥有所有权限
	if space.UserId == userId {
		return []string{
			"spaceUser:manage", // 空间用户管理权限
			"picture:view",     // 查看图片权限
			"picture:upload",   // 上传图片权限
			"picture:edit",     // 编辑图片权限
			"picture:delete",   // 删除图片权限
		}
	}

	// 查询用户在这个空间中的角色
	var spaceUser *entity.SpaceUser
	err := dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, space.Id).
		Where(dao.SpaceUser.Columns().UserId, userId).Scan(&spaceUser)
	if err != nil || spaceUser == nil {
		// 如果用户不在空间成员表中，返回基础查看权限
		return []string{"picture:view"}
	}

	// 根据用户角色分配权限
	switch spaceUser.SpaceRole {
	case "admin":
		return []string{
			"spaceUser:manage", // 空间用户管理权限
			"picture:view",     // 查看图片权限
			"picture:upload",   // 上传图片权限
			"picture:edit",     // 编辑图片权限
			"picture:delete",   // 删除图片权限
		}
	case "editor":
		return []string{
			"picture:view",   // 查看图片权限
			"picture:upload", // 上传图片权限
			"picture:edit",   // 编辑图片权限
		}
	case "viewer":
		return []string{
			"picture:view", // 查看图片权限
		}
	default:
		// 默认只有查看权限
		return []string{"picture:view"}
	}
}
