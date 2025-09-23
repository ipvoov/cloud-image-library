package space

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Add 添加空间
func (s *sSpace) Add(ctx context.Context, req *v1.SpaceAddReq) (res *v1.SpaceAddRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	if user.UserRole == consts.DefaultRole && req.SpaceLevel != 0 {
		return nil, gerror.New("无权限创建指定级别空间")
	}

	// 根据空间级别设置限制
	var maxSize, maxCount int64
	switch req.SpaceLevel {
	case 0: // 普通版
		maxSize = 100 * 1024 * 1024 // 100MB
		maxCount = 100
	case 1: // 专业版
		maxSize = 1024 * 1024 * 1024 // 1GB
		maxCount = 1000
	case 2: // 旗舰版
		maxSize = 10 * 1024 * 1024 * 1024 // 10GB
		maxCount = 10000
	default:
		return nil, gerror.New("无效的空间级别")
	}

	// 使用分布式锁防止并发创建
	lockKey := fmt.Sprintf("space:create:user:%d", user.Id)
	_, err = g.Redis().Do(ctx, "set", lockKey, 1, "NX", "EX", 10)
	if err != nil {
		return nil, gerror.New("服务繁忙,请稍后再试")
	}

	var id int64
	err = dao.Space.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 在锁内再次检查用户空间数量限制
		count, checkErr := dao.Space.Ctx(ctx).Where(dao.Space.Columns().UserId, user.Id).
			Where(dao.Space.Columns().IsDelete, 0).Count()
		if checkErr != nil {
			return gerror.New("检查用户空间失败")
		}

		// 检查团队空间限制：每个用户只能创建一个团队空间
		if req.SpaceType == consts.SpaceTypeTeam {
			teamSpaceCount, teamErr := dao.Space.Ctx(ctx).Where(dao.Space.Columns().UserId, user.Id).
				Where(dao.Space.Columns().SpaceType, consts.SpaceTypeTeam).
				Where(dao.Space.Columns().IsDelete, 0).Count()
			if teamErr != nil {
				return gerror.New("检查团队空间失败")
			}
			if teamSpaceCount > 0 {
				return gerror.New("每个用户只能创建一个团队空间")
			}
		}

		// 根据空间级别检查限制
		if req.SpaceLevel == 0 && req.SpaceType == consts.SpaceTypePrivate && count > 0 {
			return gerror.New("每个用户只能创建一个私有空间")
		}
		// 可以添加其他级别的限制，比如专业版最多3个空间等
		if req.SpaceType == consts.SpaceTypePrivate && count >= 3 { // 总空间数限制
			return gerror.New("空间数量已达上限")
		}

		// 插入空间记录
		result, insertErr := dao.Space.Ctx(ctx).Data(do.Space{
			SpaceName:  req.SpaceName,
			SpaceLevel: req.SpaceLevel,
			MaxSize:    maxSize,
			MaxCount:   maxCount,
			TotalSize:  0,
			TotalCount: 0,
			UserId:     user.Id,
			SpaceType:  req.SpaceType, // 使用请求中的空间类型
		}).Insert()
		if insertErr != nil {
			g.Log().Errorf(ctx, "创建空间失败: %v", insertErr)
			return gerror.New("创建空间失败")
		}

		id, insertErr = result.LastInsertId()
		if insertErr != nil {
			return gerror.New("获取空间ID失败")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	_, err = g.Redis().Do(ctx, "del", lockKey)
	if err != nil {
		return nil, gerror.New("服务繁忙, 请稍后再试")
	}

	return &v1.SpaceAddRes{Id: id}, nil
}

// Edit 编辑空间
func (s *sSpace) Edit(ctx context.Context, req *v1.SpaceEditReq) (res *v1.SpaceEditRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 验证空间权限
	if !s.hasSpacePermission(ctx, req.Id, user, "spaceUser:edit") {
		return nil, gerror.New("无权限编辑此空间")
	}

	// 更新空间信息
	_, err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.Id).
		Where(dao.Space.Columns().IsDelete, 0).
		Data(do.Space{
			SpaceName: req.SpaceName,
			EditTime:  gtime.Now(),
		}).Update()
	if err != nil {
		g.Log().Errorf(ctx, "编辑空间失败: %v", err)
		return nil, gerror.New("编辑空间失败")
	}

	return &v1.SpaceEditRes{
		Success: true,
	}, nil
}

// Update 更新空间
func (s *sSpace) Update(ctx context.Context, req *v1.SpaceUpdateReq) (res *v1.SpaceUpdateRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 验证空间权限（只有管理员可以更新空间配置）
	if user.UserRole != consts.Admin {
		return nil, gerror.New("无权限更新空间配置")
	}

	// 更新空间配置
	updateData := do.Space{
		UpdateTime: gtime.Now(),
	}
	if req.MaxCount > 0 {
		updateData.MaxCount = req.MaxCount
	}
	if req.MaxSize > 0 {
		updateData.MaxSize = req.MaxSize
	}

	_, err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.Id).
		Where(dao.Space.Columns().IsDelete, 0).
		Data(updateData).Update()
	if err != nil {
		g.Log().Errorf(ctx, "更新空间失败: %v", err)
		return nil, gerror.New("更新空间失败")
	}

	return &v1.SpaceUpdateRes{
		Success: true,
	}, nil
}

// Delete 删除空间
func (s *sSpace) Delete(ctx context.Context, req *v1.SpaceDeleteReq) (res *v1.SpaceDeleteRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 验证空间权限
	if !s.hasSpacePermission(ctx, req.Id, user, "spaceUser:manage") {
		return nil, gerror.New("无权限删除此空间")
	}

	// 软删除空间
	_, err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.Id).
		Data(do.Space{
			IsDelete:   1,
			UpdateTime: gtime.Now(),
		}).Update()
	if err != nil {
		g.Log().Errorf(ctx, "删除空间失败: %v", err)
		return nil, gerror.New("删除空间失败")
	}

	return &v1.SpaceDeleteRes{
		Success: true,
	}, nil
}

// hasSpacePermission 检查用户是否有空间权限
func (s *sSpace) hasSpacePermission(ctx context.Context, spaceId int64, user *v1.GetLoginUserRes, action string) bool {
	// 1. 检查是否是空间创建者
	// 2. 检查是否是空间管理员
	// 3. 根据action检查具体权限

	exist, err := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceId).
		Where(dao.Space.Columns().UserId, user.Id).
		Where(dao.Space.Columns().IsDelete, 0).Exist()
	if err != nil {
		g.Log().Errorf(ctx, "检查空间权限失败: %v", err)
		return false
	}

	return exist || user.UserRole == consts.Admin
}
