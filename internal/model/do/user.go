// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table user for DAO operations like Where/Data.
type User struct {
	g.Meta        `orm:"table:user, do:true"`
	Id            any         // id
	UserAccount   any         // 账号
	UserPassword  any         // 密码
	UserName      any         // 用户昵称
	UserAvatar    any         // 用户头像
	UserProfile   any         // 用户简介
	UserRole      any         // 用户角色：user/admin
	EditTime      *gtime.Time // 编辑时间
	CreateTime    *gtime.Time // 创建时间
	UpdateTime    *gtime.Time // 更新时间
	IsDelete      any         // 是否删除
	VipExpireTime *gtime.Time // 会员过期时间
	VipCode       any         // 会员兑换码
	VipNumber     any         // 会员编号
}
