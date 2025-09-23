// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SpaceUser is the golang structure for table space_user.
type SpaceUser struct {
	Id         int64       `json:"id"         orm:"id"         description:"id"`                       // id
	SpaceId    int64       `json:"spaceId"    orm:"spaceId"    description:"空间 id"`                    // 空间 id
	UserId     int64       `json:"userId"     orm:"userId"     description:"用户 id"`                    // 用户 id
	SpaceRole  string      `json:"spaceRole"  orm:"spaceRole"  description:"空间角色：viewer/editor/admin"` // 空间角色：viewer/editor/admin
	CreateTime *gtime.Time `json:"createTime" orm:"createTime" description:"创建时间"`                     // 创建时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"updateTime" description:"更新时间"`                     // 更新时间
}
