// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Space is the golang structure for table space.
type Space struct {
	Id         int64       `json:"id"         orm:"id"         description:"id"`                     // id
	SpaceName  string      `json:"spaceName"  orm:"spaceName"  description:"空间名称"`                   // 空间名称
	SpaceLevel int         `json:"spaceLevel" orm:"spaceLevel" description:"空间级别：0-普通版 1-专业版 2-旗舰版"` // 空间级别：0-普通版 1-专业版 2-旗舰版
	MaxSize    int64       `json:"maxSize"    orm:"maxSize"    description:"空间图片的最大总大小"`             // 空间图片的最大总大小
	MaxCount   int64       `json:"maxCount"   orm:"maxCount"   description:"空间图片的最大数量"`              // 空间图片的最大数量
	TotalSize  int64       `json:"totalSize"  orm:"totalSize"  description:"当前空间下图片的总大小"`            // 当前空间下图片的总大小
	TotalCount int64       `json:"totalCount" orm:"totalCount" description:"当前空间下的图片数量"`             // 当前空间下的图片数量
	UserId     int64       `json:"userId"     orm:"userId"     description:"创建用户 id"`                // 创建用户 id
	CreateTime *gtime.Time `json:"createTime" orm:"createTime" description:"创建时间"`                   // 创建时间
	EditTime   *gtime.Time `json:"editTime"   orm:"editTime"   description:"编辑时间"`                   // 编辑时间
	UpdateTime *gtime.Time `json:"updateTime" orm:"updateTime" description:"更新时间"`                   // 更新时间
	IsDelete   int         `json:"isDelete"   orm:"isDelete"   description:"是否删除"`                   // 是否删除
	SpaceType  int         `json:"spaceType"  orm:"spaceType"  description:"空间类型：0-私有 1-团队"`         // 空间类型：0-私有 1-团队
}
