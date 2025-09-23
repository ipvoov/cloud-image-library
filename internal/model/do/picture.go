// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Picture is the golang structure of table picture for DAO operations like Where/Data.
type Picture struct {
	g.Meta        `orm:"table:picture, do:true"`
	Id            any         // id
	Url           any         // 图片 url
	Name          any         // 图片名称
	Introduction  any         // 简介
	Category      any         // 分类
	Tags          any         // 标签（JSON 数组）
	PicSize       any         // 图片体积
	PicWidth      any         // 图片宽度
	PicHeight     any         // 图片高度
	PicScale      any         // 图片宽高比例
	PicFormat     any         // 图片格式
	UserId        any         // 创建用户 id
	CreateTime    *gtime.Time // 创建时间
	EditTime      *gtime.Time // 编辑时间
	UpdateTime    *gtime.Time // 更新时间
	IsDelete      any         // 是否删除
	ReviewStatus  any         // 审核状态：0-待审核; 1-通过; 2-拒绝
	ReviewMessage any         // 审核信息
	ReviewerId    any         // 审核人 ID
	ReviewTime    *gtime.Time // 审核时间
	ThumbnailUrl  any         // 缩略图 url
	SpaceId       any         // 空间 id（为空表示公共空间）
	PicColor      any         // 图片主色调
}
