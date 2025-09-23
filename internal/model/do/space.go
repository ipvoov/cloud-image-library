// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Space is the golang structure of table space for DAO operations like Where/Data.
type Space struct {
	g.Meta     `orm:"table:space, do:true"`
	Id         any         // id
	SpaceName  any         // 空间名称
	SpaceLevel any         // 空间级别：0-普通版 1-专业版 2-旗舰版
	MaxSize    any         // 空间图片的最大总大小
	MaxCount   any         // 空间图片的最大数量
	TotalSize  any         // 当前空间下图片的总大小
	TotalCount any         // 当前空间下的图片数量
	UserId     any         // 创建用户 id
	CreateTime *gtime.Time // 创建时间
	EditTime   *gtime.Time // 编辑时间
	UpdateTime *gtime.Time // 更新时间
	IsDelete   any         // 是否删除
	SpaceType  any         // 空间类型：0-私有 1-团队
}
