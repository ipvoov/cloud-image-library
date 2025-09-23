// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PictureDao is the data access object for the table picture.
type PictureDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  PictureColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// PictureColumns defines and stores column names for the table picture.
type PictureColumns struct {
	Id            string // id
	Url           string // 图片 url
	Name          string // 图片名称
	Introduction  string // 简介
	Category      string // 分类
	Tags          string // 标签（JSON 数组）
	PicSize       string // 图片体积
	PicWidth      string // 图片宽度
	PicHeight     string // 图片高度
	PicScale      string // 图片宽高比例
	PicFormat     string // 图片格式
	UserId        string // 创建用户 id
	CreateTime    string // 创建时间
	EditTime      string // 编辑时间
	UpdateTime    string // 更新时间
	IsDelete      string // 是否删除
	ReviewStatus  string // 审核状态：0-待审核; 1-通过; 2-拒绝
	ReviewMessage string // 审核信息
	ReviewerId    string // 审核人 ID
	ReviewTime    string // 审核时间
	ThumbnailUrl  string // 缩略图 url
	SpaceId       string // 空间 id（为空表示公共空间）
	PicColor      string // 图片主色调
}

// pictureColumns holds the columns for the table picture.
var pictureColumns = PictureColumns{
	Id:            "id",
	Url:           "url",
	Name:          "name",
	Introduction:  "introduction",
	Category:      "category",
	Tags:          "tags",
	PicSize:       "picSize",
	PicWidth:      "picWidth",
	PicHeight:     "picHeight",
	PicScale:      "picScale",
	PicFormat:     "picFormat",
	UserId:        "userId",
	CreateTime:    "createTime",
	EditTime:      "editTime",
	UpdateTime:    "updateTime",
	IsDelete:      "isDelete",
	ReviewStatus:  "reviewStatus",
	ReviewMessage: "reviewMessage",
	ReviewerId:    "reviewerId",
	ReviewTime:    "reviewTime",
	ThumbnailUrl:  "thumbnailUrl",
	SpaceId:       "spaceId",
	PicColor:      "picColor",
}

// NewPictureDao creates and returns a new DAO object for table data access.
func NewPictureDao(handlers ...gdb.ModelHandler) *PictureDao {
	return &PictureDao{
		group:    "default",
		table:    "picture",
		columns:  pictureColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *PictureDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *PictureDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *PictureDao) Columns() PictureColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *PictureDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *PictureDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *PictureDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
