// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id            int64       `json:"id"            orm:"id"            description:"id"`              // id
	UserAccount   string      `json:"userAccount"   orm:"userAccount"   description:"账号"`              // 账号
	UserPassword  string      `json:"userPassword"  orm:"userPassword"  description:"密码"`              // 密码
	UserName      string      `json:"userName"      orm:"userName"      description:"用户昵称"`            // 用户昵称
	UserAvatar    string      `json:"userAvatar"    orm:"userAvatar"    description:"用户头像"`            // 用户头像
	UserProfile   string      `json:"userProfile"   orm:"userProfile"   description:"用户简介"`            // 用户简介
	UserRole      string      `json:"userRole"      orm:"userRole"      description:"用户角色：user/admin"` // 用户角色：user/admin
	EditTime      *gtime.Time `json:"editTime"      orm:"editTime"      description:"编辑时间"`            // 编辑时间
	CreateTime    *gtime.Time `json:"createTime"    orm:"createTime"    description:"创建时间"`            // 创建时间
	UpdateTime    *gtime.Time `json:"updateTime"    orm:"updateTime"    description:"更新时间"`            // 更新时间
	IsDelete      int         `json:"isDelete"      orm:"isDelete"      description:"是否删除"`            // 是否删除
	VipExpireTime *gtime.Time `json:"vipExpireTime" orm:"vipExpireTime" description:"会员过期时间"`          // 会员过期时间
	VipCode       string      `json:"vipCode"       orm:"vipCode"       description:"会员兑换码"`           // 会员兑换码
	VipNumber     int64       `json:"vipNumber"     orm:"vipNumber"     description:"会员编号"`            // 会员编号
}
