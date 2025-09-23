// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SpaceUser is the golang structure of table space_user for DAO operations like Where/Data.
type SpaceUser struct {
	g.Meta     `orm:"table:space_user, do:true"`
	Id         any         // id
	SpaceId    any         // 空间 id
	UserId     any         // 用户 id
	SpaceRole  any         // 空间角色：viewer/editor/admin
	CreateTime *gtime.Time // 创建时间
	UpdateTime *gtime.Time // 更新时间
}
