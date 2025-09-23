package space

import (
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// GetSpaceById 根据ID获取空间 - WebSocket专用
func GetSpaceById(ctx context.Context, spaceID int64) (*entity.Space, error) {
	var space *entity.Space
	err := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceID).
		Where(dao.Space.Columns().IsDelete, 0).Scan(&space)

	if err != nil {
		return nil, gerror.New("查询空间失败")
	}

	if space == nil {
		return nil, gerror.New("空间不存在")
	}

	return space, nil
}

// CheckEditPermission 检查用户是否有编辑权限 - WebSocket专用
func CheckEditPermission(ctx context.Context, spaceID int64, userID int64) bool {
	// 检查空间是否存在
	var space *entity.Space
	err := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceID).
		Where(dao.Space.Columns().IsDelete, 0).Scan(&space)

	if err != nil || space == nil {
		return false
	}

	// 空间所有者有编辑权限
	if space.UserId == userID {
		return true
	}

	// 检查用户是否是空间成员，以及角色权限
	var spaceUser *entity.SpaceUser
	err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, spaceID).
		Where(dao.SpaceUser.Columns().UserId, userID).Scan(&spaceUser)

	if err != nil || spaceUser == nil {
		return false // 不是空间成员
	}

	// 根据角色判断编辑权限
	switch spaceUser.SpaceRole {
	case "admin", "editor":
		return true // admin 和 editor 都有编辑权限
	case "viewer":
		return false // viewer 只有查看权限
	default:
		return false
	}
}
