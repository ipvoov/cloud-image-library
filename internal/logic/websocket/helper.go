package websocket

import (
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"context"
	"net/http"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gorilla/websocket"
)

var wsUpGrader = websocket.Upgrader{
	// 开发环境下允许任何来源
	// 生产环境下需要实现适当的来源检查以确保安全
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// 升级失败时的错误处理器
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		// 在这里实现错误处理逻辑
		g.Log().Error(r.Context(), reason)
	},
}

// GetPictureById 根据ID获取图片
func GetPictureById(ctx context.Context, pictureID int64) (*entity.Picture, error) {
	var picture *entity.Picture
	err := dao.Picture.Ctx(ctx).Where(dao.Picture.Columns().Id, pictureID).
		Where(dao.Picture.Columns().IsDelete, 0).Scan(&picture)

	if err != nil {
		return nil, gerror.New("查询图片失败")
	}

	if picture == nil {
		return nil, gerror.New("图片不存在")
	}

	return picture, nil
}

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

// CheckEditPermission 检查用户是否有编辑权限
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
