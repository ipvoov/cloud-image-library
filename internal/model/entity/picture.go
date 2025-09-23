// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Picture is the golang structure for table picture.
type Picture struct {
	Id            int64       `json:"id"            orm:"id"            description:"id"`                     // id
	Url           string      `json:"url"           orm:"url"           description:"图片 url"`                 // 图片 url
	Name          string      `json:"name"          orm:"name"          description:"图片名称"`                   // 图片名称
	Introduction  string      `json:"introduction"  orm:"introduction"  description:"简介"`                     // 简介
	Category      string      `json:"category"      orm:"category"      description:"分类"`                     // 分类
	Tags          string      `json:"tags"          orm:"tags"          description:"标签（JSON 数组）"`            // 标签（JSON 数组）
	PicSize       int64       `json:"picSize"       orm:"picSize"       description:"图片体积"`                   // 图片体积
	PicWidth      int         `json:"picWidth"      orm:"picWidth"      description:"图片宽度"`                   // 图片宽度
	PicHeight     int         `json:"picHeight"     orm:"picHeight"     description:"图片高度"`                   // 图片高度
	PicScale      float64     `json:"picScale"      orm:"picScale"      description:"图片宽高比例"`                 // 图片宽高比例
	PicFormat     string      `json:"picFormat"     orm:"picFormat"     description:"图片格式"`                   // 图片格式
	UserId        int64       `json:"userId"        orm:"userId"        description:"创建用户 id"`                // 创建用户 id
	CreateTime    *gtime.Time `json:"createTime"    orm:"createTime"    description:"创建时间"`                   // 创建时间
	EditTime      *gtime.Time `json:"editTime"      orm:"editTime"      description:"编辑时间"`                   // 编辑时间
	UpdateTime    *gtime.Time `json:"updateTime"    orm:"updateTime"    description:"更新时间"`                   // 更新时间
	IsDelete      int         `json:"isDelete"      orm:"isDelete"      description:"是否删除"`                   // 是否删除
	ReviewStatus  int         `json:"reviewStatus"  orm:"reviewStatus"  description:"审核状态：0-待审核; 1-通过; 2-拒绝"` // 审核状态：0-待审核; 1-通过; 2-拒绝
	ReviewMessage string      `json:"reviewMessage" orm:"reviewMessage" description:"审核信息"`                   // 审核信息
	ReviewerId    int64       `json:"reviewerId"    orm:"reviewerId"    description:"审核人 ID"`                 // 审核人 ID
	ReviewTime    *gtime.Time `json:"reviewTime"    orm:"reviewTime"    description:"审核时间"`                   // 审核时间
	ThumbnailUrl  string      `json:"thumbnailUrl"  orm:"thumbnailUrl"  description:"缩略图 url"`                // 缩略图 url
	SpaceId       int64       `json:"spaceId"       orm:"spaceId"       description:"空间 id（为空表示公共空间）"`        // 空间 id（为空表示公共空间）
	PicColor      string      `json:"picColor"      orm:"picColor"      description:"图片主色调"`                  // 图片主色调
}
